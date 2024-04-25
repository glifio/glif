package fsjournal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	clk "github.com/raulk/clock"
	"golang.org/x/xerrors"

	"github.com/glifio/glif/v2/journal"
)

var clock = clk.New()

const RFC3339nocolon = "2006-01-02T150405Z0700"

// fsJournal is a basic journal backed by files on a filesystem.
type fsJournal struct {
	journal.EventTypeRegistry

	dir       string
	sizeLimit int64

	fi    *os.File
	fSize int64

	incoming chan *journal.Event

	closing chan struct{}
	closed  chan struct{}
}

// OpenFSJournal constructs a rolling filesystem journal, with a default
// per-file size limit of 1GiB.
func OpenFSJournal(journalPath string, disabled journal.DisabledEvents) (journal.Journal, error) {
	dir := filepath.Join(journalPath, "journal")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to mk directory %s for file journal: %w", dir, err)
	}

	f := &fsJournal{
		EventTypeRegistry: journal.NewEventTypeRegistry(disabled),
		dir:               dir,
		sizeLimit:         1 << 30,
		incoming:          make(chan *journal.Event, 32),
		closing:           make(chan struct{}),
		closed:            make(chan struct{}),
	}

	var nfi *os.File
	var nfSize int64
	current := filepath.Join(f.dir, "glif-journal.ndjson")
	if fi, err := os.Stat(current); err == nil && !fi.IsDir() {
		nfi, err = os.OpenFile(current, os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return nil, xerrors.Errorf("failed to open journal file: %w", err)
		}
		nfSize = fi.Size()
	} else {
		nfi, err = os.Create(current)
		if err != nil {
			return nil, xerrors.Errorf("failed to create journal file: %w", err)
		}
	}
	f.fi = nfi
	f.fSize = nfSize

	go f.runLoop()

	return f, nil
}

func (f *fsJournal) RecordEvent(evtType journal.EventType, supplier func() interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic while recording journal event; type=%s, err=%v", evtType, r)
		}
	}()

	if !evtType.Enabled() {
		return
	}

	je := &journal.Event{
		EventType: evtType,
		Timestamp: clock.Now(),
		Data:      supplier(),
	}
	select {
	case f.incoming <- je:
	case <-f.closing:
		log.Printf("journal closed but tried to log event: %s\n", je)
	}
}

func (f *fsJournal) ReadEvents() ([]journal.Event, error) {
	evts := []journal.Event{}
	b := make([]byte, f.fSize)

	_, err := f.fi.Read(b)
	if err != nil {
		return nil, err
	}

	bstrs := strings.Split(string(b), "\n")
	for _, bs := range bstrs {
		evt := &journal.Event{}
		if bs == "" {
			continue
		}
		err := json.Unmarshal([]byte(bs), evt)
		if err != nil {
			return nil, err
		}
		evts = append(evts, *evt)
	}
	return evts, nil
}

func (f *fsJournal) Close() error {
	close(f.closing)
	<-f.closed
	return nil
}

func (f *fsJournal) putEvent(evt *journal.Event) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	n, err := f.fi.Write(append(b, '\n'))
	if err != nil {
		return err
	}

	f.fSize += int64(n)

	if f.fSize >= f.sizeLimit {
		_ = f.rollJournalFile()
	}

	return nil
}

func (f *fsJournal) rollJournalFile() error {
	if f.fi != nil {
		_ = f.fi.Close()
	}
	current := filepath.Join(f.dir, "glif-journal.ndjson")
	rolled := filepath.Join(f.dir, fmt.Sprintf(
		"glif-journal-%s.ndjson",
		clock.Now().Format(RFC3339nocolon),
	))

	// check if journal file exists
	if fi, err := os.Stat(current); err == nil && !fi.IsDir() {
		err := os.Rename(current, rolled)
		if err != nil {
			return xerrors.Errorf("failed to roll journal file: %w", err)
		}
	}

	nfi, err := os.Create(current)
	if err != nil {
		return xerrors.Errorf("failed to create journal file: %w", err)
	}

	f.fi = nfi
	f.fSize = 0

	return nil
}

func (f *fsJournal) runLoop() {
	defer close(f.closed)

	for {
		select {
		case je := <-f.incoming:
			if err := f.putEvent(je); err != nil {
				log.Print("failed to write out journal event", "event", je, "err", err)
			}
		case <-f.closing:
			_ = f.fi.Close()
			return
		}
	}
}

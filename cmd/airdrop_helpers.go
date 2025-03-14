package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/go-pools/constants"
	"github.com/spf13/cobra"
)

const (
	MAINNET_GLIF_LTD_DELEGATEE = "0xfF0000000000000000000000000000000035909b"
	TESTNET_GLIF_LTD_DELEGATEE = "0xFf000000000000000000000000000000000254d7"
)

type language string

// question 1
const (
	VOTE_QUESTION_LANGUAGE                = "You are about to claim your airdrop, we have a few extra questions for you. What language would you like to continue in?"
	VOTE_OPTION_LANGUAGE_ENGLISH language = "English"
	VOTE_OPTION_LANGUAGE_CHINESE language = "中文"
)

func getUserLanguagePreference() (language, error) {
	var selectedLanguage string
	languagePrompt := &survey.Select{
		Message: VOTE_QUESTION_LANGUAGE,
		Options: []string{string(VOTE_OPTION_LANGUAGE_ENGLISH), string(VOTE_OPTION_LANGUAGE_CHINESE)},
	}
	err := survey.AskOne(languagePrompt, &selectedLanguage)
	if err != nil {
		return "", fmt.Errorf("failed to get user input: %w", err)
	}

	return language(selectedLanguage), nil
}

// question 2
const (
	VOTE_QUESTION_DELEGATEE               = "The GLF Token has governance capabilities. Who do you want to delegate your vote to?"
	VOTE_QUESTION_DELEGATEE_CHINESE       = "GLF 代币具有治理功能。你想将你的投票委托给谁？"
	VOTE_OPTION_DELEGATEE_GLIF_LTD        = "Glif Ltd. - the development company"
	VOTE_OPTION_DELEGATEE_MYSELF          = "Myself - I will vote on proposals"
	VOTE_OPTION_DELEGATEE_SOMEONE_ELSE    = "Someone else - I know their address"
	VOTE_OPTION_DELEGATEE_GLIF_LTD_CH     = "Glif Ltd. - 开发公司"
	VOTE_OPTION_DELEGATEE_MYSELF_CH       = "我自己 - 我将对提案投票"
	VOTE_OPTION_DELEGATEE_SOMEONE_ELSE_CH = "其他人 - 我知道他们的地址"
)

// question 2.5
const (
	VOTE_QUESTION_DELEGATEE_ADDRESS         = "Please enter the address of the person you want to delegate your vote to:"
	VOTE_QUESTION_DELEGATEE_ADDRESS_CHINESE = "请输入你想委托投票的地址:"
)

func getDelegateeAddressByLanguage(ctx context.Context, self common.Address, lang language) (common.Address, error) {
	var delegatee common.Address
	var selectedOption string
	var delegatePrompt *survey.Select

	if lang == VOTE_OPTION_LANGUAGE_CHINESE {
		delegatePrompt = &survey.Select{
			Message: VOTE_QUESTION_DELEGATEE_CHINESE,
			Options: []string{VOTE_OPTION_DELEGATEE_GLIF_LTD_CH, VOTE_OPTION_DELEGATEE_MYSELF_CH, VOTE_OPTION_DELEGATEE_SOMEONE_ELSE_CH},
		}
	} else {
		delegatePrompt = &survey.Select{
			Message: VOTE_QUESTION_DELEGATEE,
			Options: []string{VOTE_OPTION_DELEGATEE_GLIF_LTD, VOTE_OPTION_DELEGATEE_MYSELF, VOTE_OPTION_DELEGATEE_SOMEONE_ELSE},
		}
	}

	err := survey.AskOne(delegatePrompt, &selectedOption)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get user input: %w", err)
	}

	switch selectedOption {
	case VOTE_OPTION_DELEGATEE_GLIF_LTD, VOTE_OPTION_DELEGATEE_GLIF_LTD_CH:
		testDrop := PoolsSDK.Query().ChainID().Int64() != constants.MainnetChainID
		if testDrop {
			delegatee = common.HexToAddress(TESTNET_GLIF_LTD_DELEGATEE)
		} else {
			delegatee = common.HexToAddress(MAINNET_GLIF_LTD_DELEGATEE)
		}
	case VOTE_OPTION_DELEGATEE_MYSELF, VOTE_OPTION_DELEGATEE_MYSELF_CH:
		// Assuming the account address is available in the context
		delegatee = self
	case VOTE_OPTION_DELEGATEE_SOMEONE_ELSE, VOTE_OPTION_DELEGATEE_SOMEONE_ELSE_CH:
		message := VOTE_QUESTION_DELEGATEE_ADDRESS
		if lang == VOTE_OPTION_LANGUAGE_CHINESE {
			message = VOTE_QUESTION_DELEGATEE_ADDRESS_CHINESE
		}
		prompt := &survey.Input{
			Message: message,
		}
		var delegateAddress string
		err := survey.AskOne(prompt, &delegateAddress)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to get user input: %w", err)
		}
		delegatee, err = AddressOrAccountNameToEVM(ctx, delegateAddress)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to parse address: %w", err)
		}
	}

	return delegatee, nil
}

const (
	VOTE_QUESTION_ROADMAP         = "<English roadmap goes here> "
	VOTE_QUESTION_ROADMAP_CHINESE = "请查看路线图"
	VOTE_OPTION_ROADMAP_CONT_EN   = "Continue"
	VOTE_OPTION_ROADMAP_CONT_CH   = "继续"
)

// question 3
const (
	VOTE_QUESTION_ACCEPT_TERMS         = "Please accept the terms and conditions at https://glif.io/terms. Do you accept?"
	VOTE_QUESTION_ACCEPT_TERMS_CHINESE = "请接受 https://glif.io/terms 的条款和条件。你接受吗？"
)

func getDelegateeAddress(ctx context.Context, userAddr common.Address) (common.Address, error) {
	lang, err := getUserLanguagePreference()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get user language preference: %w", err)
	}
	return getDelegateeAddressByLanguage(ctx, userAddr, lang)
}

var getDelegateeCmd = &cobra.Command{
	Use:   "get-delegatee [address]",
	Short: "Get the delegatee address for a given address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strAddr := args[0]
		addr, err := AddressOrAccountNameToEVM(cmd.Context(), strAddr)
		if err != nil {
			logFatalf("Failed to parse address %s", err)
		}

		delegatee, err := getDelegateeAddress(cmd.Context(), addr)
		if err != nil {
			logFatalf("Failed to get delegatee address %s", err)
		}

		fmt.Printf("%s\n", delegatee.Hex())
	},
}

func init() {
	airdropCmd.AddCommand(getDelegateeCmd)
}

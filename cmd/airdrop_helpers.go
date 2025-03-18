package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/go-pools/constants"
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
	VOTE_OPTION_DELEGATEE_MYSELF_CH       = "我自己 - 我将参与提案投票"
	VOTE_OPTION_DELEGATEE_SOMEONE_ELSE_CH = "其他人 - 我知道他们的地址"
)

// question 2.5
const (
	VOTE_QUESTION_DELEGATEE_ADDRESS         = "Please enter the address of the person you want to delegate your vote to:"
	VOTE_QUESTION_DELEGATEE_ADDRESS_CHINESE = "请输入你想委托投票的人的地址："
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
	VOTE_QUESTION_ROADMAP         = "GLIF is entering the next phase of its journey. Here are the most important things you need to know:\n\n1. GLIF is going multichain - the $GLF Token is the first step of our multichain expansion.\n\n2. The token will be useful in the near future for obtaining benefits from GLIF - like cheaper borrowing rates.\n\n3. Join our discord https://discord.gg/5qsJjsP3Re for the most important announcements.\n\nRead more about it: https://medium.com/@glifio/chapter-iii-expansion-7953b7cc6444"
	VOTE_QUESTION_ROADMAP_CHINESE = "GLIF 正在进入其旅程的下一个阶段。以下是最重要的事情，您需要了解：\n\nGLIF 正在向多链扩展 — $GLF 代币是我们多链扩展的第一步。\n\n2. 代币将在不久的将来用于获取 GLIF 的好处，例如更便宜的借贷利率。\n\n3. 加入我们的 Discord https://discord.gg/5qsJjsP3Re 获取最重要的公告。\n\n阅读更多内容：https://medium.com/@glifio/chapter-iii-expansion-7953b7cc6444"
	VOTE_OPTION_ROADMAP_CONT_EN   = "Finish"
	VOTE_OPTION_ROADMAP_CONT_CH   = "完成"
)

func getRoadmapByLanguage(lang language) error {
	var selectedOption string
	var prompt *survey.Select
	if lang == VOTE_OPTION_LANGUAGE_CHINESE {
		prompt = &survey.Select{
			Message: VOTE_QUESTION_ROADMAP_CHINESE,
			Options: []string{VOTE_OPTION_ROADMAP_CONT_CH},
			Default: VOTE_OPTION_ROADMAP_CONT_CH,
		}
	} else {
		prompt = &survey.Select{
			Message: VOTE_QUESTION_ROADMAP,
			Options: []string{VOTE_OPTION_ROADMAP_CONT_EN},
			Default: VOTE_OPTION_ROADMAP_CONT_EN,
		}
	}

	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}

	return nil
}

// question 3
const (
	VOTE_QUESTION_ACCEPT_TERMS         = "By claiming your airdrop, you are agreeing to the terms and conditions at https://glif.io/terms. Do you accept?"
	VOTE_QUESTION_ACCEPT_TERMS_CHINESE = "通过领取空投，你同意 https://glif.io/terms 的条款和条件。你接受吗？"
	VOTE_OPTION_ACCEPT_TERMS_YES       = "Yes"
	VOTE_OPTION_ACCEPT_TERMS_YES_CH    = "是"
	VOTE_OPTION_ACCEPT_TERMS_NO        = "No"
	VOTE_OPTION_ACCEPT_TERMS_NO_CH     = "否"
)

func getAcceptTermsByLanguage(lang language) error {
	var selectedOption string
	var prompt *survey.Select
	if lang == VOTE_OPTION_LANGUAGE_CHINESE {
		prompt = &survey.Select{
			Message: VOTE_QUESTION_ACCEPT_TERMS_CHINESE,
			Options: []string{VOTE_OPTION_ACCEPT_TERMS_YES_CH, VOTE_OPTION_ACCEPT_TERMS_NO_CH},
			Default: VOTE_OPTION_ACCEPT_TERMS_YES_CH,
		}
	} else {
		prompt = &survey.Select{
			Message: VOTE_QUESTION_ACCEPT_TERMS,
			Options: []string{VOTE_OPTION_ACCEPT_TERMS_YES, VOTE_OPTION_ACCEPT_TERMS_NO},
			Default: VOTE_OPTION_ACCEPT_TERMS_YES,
		}
	}

	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}

	if selectedOption == VOTE_OPTION_ACCEPT_TERMS_YES || selectedOption == VOTE_OPTION_ACCEPT_TERMS_YES_CH {
		return nil
	} else if selectedOption == VOTE_OPTION_ACCEPT_TERMS_NO || selectedOption == VOTE_OPTION_ACCEPT_TERMS_NO_CH {
		return fmt.Errorf("user did not accept terms")
	}

	return fmt.Errorf("invalid option: %s", selectedOption)
}

func interactiveClaimExp(ctx context.Context, addr common.Address) (common.Address, error) {
	fromAddr, err := AddressOrAccountNameToEVM(ctx, addr.Hex())
	if err != nil {
		logFatalf("Failed to parse address %s", err)
	}

	lang, err := getUserLanguagePreference()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get user language preference: %w", err)
	}

	if err := getAcceptTermsByLanguage(lang); err != nil {
		return common.Address{}, err
	}

	delegatee, err := getDelegateeAddressByLanguage(ctx, fromAddr, lang)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get delegatee address: %w", err)
	}

	roadmapErr := getRoadmapByLanguage(lang)
	if roadmapErr != nil {
		return common.Address{}, fmt.Errorf("failed to get roadmap: %w", roadmapErr)
	}

	return delegatee, nil
}

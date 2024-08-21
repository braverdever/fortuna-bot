package main

import (
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gochain/gochain/v4/crypto"
	"github.com/l3njo/rochambeau/database"
	"github.com/l3njo/rochambeau/models"
	"github.com/yanzay/tbot/v2"
)

type ChainStatus struct {
	Name         string
	AnotherField int
}

func createOtherWallet(privateKeyString string) (*Wallet, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key to ECDSA: %w", err)
	}

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)

	// Return a Wallet1 struct instead of just the address
	return &Wallet{
		Address:    address,
		PrivateKey: privateKeyString,
		// Set other fields as necessary
	}, nil
}

func (a *application) startHandler(m *tbot.Message) {
	welcomeMsg := `
	The best trading bot-Fortuna
	Your favourite trading bot

	[📖 Docs](https://example.com/docs)
	💬 Official Chat
	🌍 Website

	Paste a contract address or pick an option to get started.
	    `

	buttons := makeGeneralButtons() // Removed the trailing comma

	a.client.SendMessage(m.Chat.ID, welcomeMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

// In handlers.go, modify the walletHandler function

func (a *application) handleBridge(m *tbot.Message) {
	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}

	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	var walletDetails strings.Builder
	for i, wallet := range wallets {
		//walletUrl := fmt.Sprintf("https://etherscan.io/address/%s", wallet.Address)
		linkText := "Balance"
		// Create a clickable "Balance" text that opens the walletUrl when clicked
		//clickableBalance := fmt.Sprintf("[%s](%s)", linkText, walletUrl)
		// balanceButtons := tbot.InlineKeyboardButton{Text: "Balance", URL: walletUrl}
		walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n %s\n", i+1, linkText, wallet.AccountWorth, wallet.Address))
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s)

Wallet Worth: $0

Your currently added wallets:
%s`, currentChainStatus, walletDetails.String())

	buttons := makeBridgeButtons()

	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) walletHandler(m *tbot.Message) {
	chainCurrentStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}
	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	var walletDetails strings.Builder
	for i, wallet := range wallets {
		//walletUrl := fmt.Sprintf("https://etherscan.io/address/%s", wallet.Address)
		linkText := "Balance"
		// Create a clickable "Balance" text that opens the walletUrl when clicked
		//clickableBalance := fmt.Sprintf("[%s](%s)", linkText, walletUrl)
		// balanceButtons := tbot.InlineKeyboardButton{Text: "Balance", URL: walletUrl}
		walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n %s\n", i+1, linkText, wallet.AccountWorth, wallet.Address))
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s)

Wallet Worth: $0

Your currently added wallets:
%s`, chainCurrentStatus, walletDetails.String())

	buttons := makeWalletButtons()
	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) walletSettingsHandler(m *tbot.Message) {
	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}
	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	var walletDetails strings.Builder
	for i, wallet := range wallets {
		//walletUrl := fmt.Sprintf("https://etherscan.io/address/%s", wallet.Address)
		linkText := "Balance"
		// Create a clickable "Balance" text that opens the walletUrl when clicked
		//clickableBalance := fmt.Sprintf("[%s](%s)", linkText, walletUrl)
		// balanceButtons := tbot.InlineKeyboardButton{Text: "Balance", URL: walletUrl}
		walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n %s\n", i+1, linkText, wallet.AccountWorth, wallet.Address))
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s)

Wallet Worth: $0

Your currently added wallets:
%s`, currentChainStatus, walletDetails.String())

	buttons := makeWalletSettingsButtons()
	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) handleTransferCrypto(m *tbot.Message) {

	transferMsg := `
Settings > Transfers

Use these options to transfer balances and tokens between wallets on the same chain.
`

	buttons := makeTransferButtons()
	a.client.SendMessage(m.Chat.ID, transferMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) privateKeyHandler(m *tbot.Message) {

	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}

	privateKeyArray, err := database.GetAllWallets(a.db)

	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	var walletDetails strings.Builder
	for i, wallet := range privateKeyArray {
		walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", i+1, wallet.ChainScanLabel, wallet.AccountWorth, wallet.PrivateKey))
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s)

Wallet Worth: $0

Your currently added wallets:
%s`, currentChainStatus, walletDetails.String())

	buttons := makeWalletButtons()
	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) createWalletHandler(m *tbot.Message) {

	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	/*privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyString := hex.EncodeToString(privateKeyBytes)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		/*client, err := ethclient.Dial("https://privatekeyfinder.io/private-keys/ethereum")
		if err != nil {
			log.Fatal(err)
		}

		account := common.HexToAddress(address)
		_balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Fatal(err)
		}

		_accountWorth := int(_balance.Int64())


		// Increment the click count

		newWallet := &models.Wallet{
			ChatId:         2,
			ChainScanLabel: "Balance",
			AccountWorth:   0,
			PrivateKey:     privateKeyString,
			Address:        address,
		}
		database.CreateWallet(a.db, newWallet)
		getWallet, err := database.GetAllWallets(a.db)
		if err != nil {
			// Handle error appropriately
			log.Printf("Error fetching wallet: %v", err)
			return
		}

		// Define the messages
		var walletDetails strings.Builder
		for i, wallet := range getWallet {
			walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", i+1, wallet.ChainScanLabel, wallet.AccountWorth, wallet.Address))
		}

		walletMsg := fmt.Sprintf(`
	Settings > Wallets (🔗ETH)

	Wallet Worth: $0

	Your currently added wallets:
	%s
	`, walletDetails.String())

		walletMsg += "\n"

		buttons := makeWalletButtons()
		a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))*/

	url := "https://app.dynamicauth.com/api/v0/environments/ccda0d02-8563-4108-a6d1-7b0167e0ca81/embeddedWallets"

	payload := strings.NewReader(`{
			"type": "email",
			"chains": ["EVM", "SOL"],
			"smsCountryCode": {
				"isoCountryCode": "FI",
				"phoneCountryCode": "358"
			},
			"identifier": "artur.dranhoi1025@gmail.com"
		}`)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Authorization", "Bearer dyn_y0VDm6LdL259JZkrmcNfDJtSsBcCk3BcnALaSByiqnnfpBLJRo4fE2tc")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	var userResponse UserResponse
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return
	}

	// Check if there are any verified credentials and handle accordingly
	if len(userResponse.User.VerifiedCredentials) > 0 {
		// Assuming you want to use the first verified credential's address
		firstCredential := userResponse.User.VerifiedCredentials[0]
		address := firstCredential.Address

		// Now, create a new wallet entry using this address
		newWallet := &models.Wallet{
			ChatId:         2,
			ChainScanLabel: "Balance",
			AccountWorth:   0,
			PrivateKey:     "", // Assuming you don't have a private key to set here; adjust as needed
			Address:        address,
		}

		database.CreateWallet(a.db, newWallet)
		getWallet, err := database.GetAllWallets(a.db)
		if err != nil {
			log.Printf("Error fetching wallet: %v", err)
			return
		}

		// Define the messages
		var walletDetail strings.Builder
		for i, wallet := range getWallet {
			walletDetail.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", i+1, wallet.ChainScanLabel, wallet.AccountWorth, wallet.Address))
		}

		walletMsg := fmt.Sprintf(`
		Settings > Wallets (🔗%s)

		Wallet Worth: $0

		Your currently added wallets:
		%s
		`, currentChainStatus, walletDetail.String())

		walletMsg += "\n"

		buttons := makeWalletButtons()
		a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
	} else {
		log.Println("No verified credentials found.")
		return
	}
}

func (a *application) defaultWalletHandler(m *tbot.Message) {

	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}

	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	var walletDetails strings.Builder
	for i, wallet := range wallets {
		//walletUrl := fmt.Sprintf("https://etherscan.io/address/%s", wallet.Address)
		linkText := "Balance"
		// Create a clickable "Balance" text that opens the walletUrl when clicked
		//clickableBalance := fmt.Sprintf("[%s](%s)", linkText, walletUrl)
		// balanceButtons := tbot.InlineKeyboardButton{Text: "Balance", URL: walletUrl}
		walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n %s\n", i+1, linkText, wallet.AccountWorth, wallet.Address))
	}

	walletMsg := fmt.Sprintf(`
	Settings > Wallets (🔗%s) > Default Wallets

Select the wallets that you want to be preselected when you create a new buy or snipe monitor.
	%s`, currentChainStatus, walletDetails.String())

	buttons := makeDefaultButtons()

	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) snipeWalletsSelectHandler(m *tbot.Message) {
	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}

	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s) > Default Wallets

Select the wallets to be preselected and click Done to confirm.
`, currentChainStatus)

	// Correct approach to create inline keyboard buttons with tbot
	buttons := make([][]tbot.InlineKeyboardButton, len(wallets))
	for i := range wallets {
		buttonText := fmt.Sprintf("Wallet %d", i+1)
		button := tbot.InlineKeyboardButton{
			Text:         buttonText,
			CallbackData: fmt.Sprintf("wallet_select_%d", i+1),
		}
		buttons[i] = []tbot.InlineKeyboardButton{button}
	}

	// Assuming tbot provides a method to create inline keyboard markup directly
	inlineKeyboardMarkup := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	a.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(inlineKeyboardMarkup))
}

func (a *application) selectFromWalletHandler(m *tbot.Message) {
	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}

	wallets, err := database.GetAllWallets(app.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		app.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	walletMsg := fmt.Sprintf(`
Settings > Wallets (🔗%s) > Default Wallets

Select the wallet you want to transfer from.
`, currentChainStatus)

	buttons := make([][]tbot.InlineKeyboardButton, len(wallets))

	for i, wallet := range wallets {
		buttonText := fmt.Sprintf("%s (%d)", wallet.Address, i+1)
		button := tbot.InlineKeyboardButton{
			Text:         buttonText,
			CallbackData: fmt.Sprintf("wallet_from_select_%d", i+1),
		}
		buttons[i] = []tbot.InlineKeyboardButton{button}
	}

	inlineKeyboardMarkup := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(inlineKeyboardMarkup))

}

func (a *application) regenerateWalletButtons(selectedWallets map[int]bool) [][]tbot.InlineKeyboardButton {
	var buttons [][]tbot.InlineKeyboardButton

	// Assuming you have a list of all wallets
	for i := range a.allWallets {
		isSelected := selectedWallets[i]
		buttonText := fmt.Sprintf("Wallet %d", i+1)
		if isSelected {
			buttonText = fmt.Sprintf("[Selected] Wallet %d", i+1)
		}

		button := tbot.InlineKeyboardButton{
			Text:         buttonText,
			CallbackData: fmt.Sprintf("wallet_select_%d", i),
		}
		buttons = append(buttons, []tbot.InlineKeyboardButton{button})
	}

	return buttons
}

func (a *application) UpdateWalletSelection(db *sql.DB, walletIndex int, m *tbot.Message) {
	// Assuming db is your database connection and you have a way to update the wallet's selection state
	// This is a placeholder for whatever logic you need to execute to update the wallet's selection state in your database or application state
	// For example, you might insert or remove an entry in a table or update a field in a document in your database
	log.Println("UpdateWalletSelection entered")
	// Toggle the selection state of the wallet

	if _, exists := a.selectedWalletIndices[walletIndex]; exists {
		delete(a.selectedWalletIndices, walletIndex)
	} else {
		if a.selectedWalletIndices == nil {
			a.selectedWalletIndices = make(map[int]bool)
		}
		a.selectedWalletIndices[walletIndex] = true
	}
	log.Printf("State toggled: %+v", a.selectedWalletIndices)
	// After toggling, regenerate the wallet buttons to reflect the new selection state
	buttons := a.regenerateWalletButtons(a.selectedWalletIndices)

	log.Printf("Regenerated buttons: %+v", buttons)
	markup := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	if _, err := a.client.SendMessage(m.Chat.ID, "Updated Wallet Selection", tbot.OptInlineKeyboardMarkup(markup)); err != nil {
		log.Printf("Error sending message with buttons: %v", err)
	}

	formattedMessage := fmt.Sprintf("Your Wallet %d is set as default.", walletIndex)
	a.client.SendMessage(m.Chat.ID, formattedMessage)
	// Here, you would send the updated buttons back to the user, similar to the callbackHandler example provided earlier
}

func (a *application) handleWalletSelectionForTransfer(walletIndex int, m *tbot.Message) {
	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	if walletIndex >= len(wallets) || walletIndex < 0 {
		a.client.SendMessage(m.Chat.ID, "Invalid wallet selection.")
		return
	}

	selectedWallet := wallets[walletIndex]

	walletMsg := fmt.Sprintf(`Settings > Transfers (🔗ETH)
From: %s
Available balance: 0Ξ

Select the wallet you want to transfer to:
`, selectedWallet.Address)

	buttons := make([][]tbot.InlineKeyboardButton, len(wallets))

	for i, wallet := range wallets {
		buttonText := fmt.Sprintf("%s (%d)", wallet.Address, i+1)
		button := tbot.InlineKeyboardButton{
			Text:         buttonText,
			CallbackData: fmt.Sprintf("wallet_to_select_%d", i+1),
		}
		buttons[i] = []tbot.InlineKeyboardButton{button}
	}

	inlineKeyboardMarkup := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(inlineKeyboardMarkup))

	// Here, you would typically start the transfer process or another action involving the selected wallet
	// For example, prompt the user for the recipient address or confirm the action
}

func (a *application) defaultSettingsHandler(m *tbot.Message) {
	status, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}
	defaultSettingMsg := fmt.Sprintf(
		`Settings > Defaults (🔗%s)
Set default parameters for your bot:
___________________________________`, status)

	buttons := makeDefaultSettingsButtons()

	a.client.SendMessage(m.Chat.ID, defaultSettingMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) handleTokenTransfer(m *tbot.Message) {
	// Example parameters, replace with actual values from the user input or database
	fromChain := "POL"
	toChain := "BSC"
	fromToken := "USDC"
	toToken := "USDC"
	amount := "1"
	fromAddress := "0x60d64E6F87bdf0FB45e424bf6bBC84b2DA0cfCD7"
	toAddress := "0xA81AB00dDDD0aD908FF5422bbfcBf6daDacA6F25"

	// Create the request
	transferRequest := WormholeRequest{
		FromChain:   fromChain,
		ToChain:     toChain,
		FromToken:   fromToken,
		ToToken:     toToken,
		Amount:      amount,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
	}

	// Call the transfer function
	transactionID, err := transferTokensViaWormhole(transferRequest)
	if err != nil {
		a.client.SendMessage(m.Chat.ID, "Failed to initiate token transfer: "+err.Error(), nil)
		return
	}

	// Send a confirmation message with the transaction ID
	a.client.SendMessage(m.Chat.ID, "Token transfer initiated. Transaction ID: "+transactionID, nil)
}

func convertChainStatusNameToType(name string) database.ChainStatusOne {
	switch name {
	case "Ethereum":
		return database.Ethereum
	case "BSC":
		return database.BSC
	case "Blast":
		return database.Blast
	case "Base":
		return database.Base
	case "Avax":
		return database.Avax
	case "Solana":
		return database.Solana
	default:
		// Handle unknown cases appropriately, maybe log an error or return a special value indicating an unknown status
		return database.Ethereum // Assuming there's an Unknown value or similar
	}
}

func (a *application) warRoomHandler(m *tbot.Message) {
	t := time.Now()
	remainTime := t.Second()
	warRoomMsg := fmt.Sprintf(`
	You are a Recruit Trader (Level 1)

The journey begins with Recruit Traders, the brave enlistees of
the DeFi army. As they embark on their trading campaign, they
accumulate valuable experience points with each successful 
trade. These points serve as the initial rank insignia on their path 
to becoming formidable DeFi Commanders.

Daily Missions:
⏳ Time left: %d seconds
10 Trades in 24 hours: +50 points bonus
5 Profitable trades: +100 points bonus
5x profit on one trade +100 points bonus

Stats:
⚖️ Trades: 0 
📊 Volume: 0.00
📅 Streak: 1
🤑 Daily Profitable trades: 0

🔥 Points earned: 10
💣 Points to next promotion: 290

 Keep on dominating the blockchain!⚔️`, remainTime)

	button := tbot.InlineKeyboardButton{Text: "Dismiss Message", CallbackData: "dismiss_warroom_message"}
	inlineKeyboard := tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{{button}},
	}

	a.client.SendMessage(m.Chat.ID, warRoomMsg, tbot.OptInlineKeyboardMarkup(&inlineKeyboard))
}

func (a *application) settingsHandler(m *tbot.Message, status ChainStatus) {

	settingsMsg := fmt.Sprintf(`
Settings (🔗%s)
Select an option below to configure:
___________________________________
`, status.Name)
	buttons := makeSettingsButtons()
	a.client.SendMessage(m.Chat.ID, settingsMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) gasPresetHandler(m *tbot.Message) {

	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		// Handle the error appropriately. For example, log it or return it from the function.
		log.Printf("Error getting chain status: %v", err)
		return // Or handle the error in another way suitable for your application's flow.
	}
	gasPresetMsg := fmt.Sprintf(`
	Settings > Presets > Gas (🔗%s)

Click on the buttons below to specify new gas values
	`, currentChainStatus)
	inlineKeyboard := makeGasPresetButtons()

	a.client.SendMessage(m.Chat.ID, gasPresetMsg, tbot.OptInlineKeyboardMarkup(inlineKeyboard))
}
func (a *application) sendCodeSnippet(m *tbot.Message) {
	code := "my fortuna bot"
	formattedCode := fmt.Sprintf("```\n%s\n```", code) // Markdown formatting for code block
	_, err := a.client.SendMessage(m.Chat.ID, formattedCode, tbot.OptParseModeMarkdown)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (a *application) presetSettingsHander(m *tbot.Message) {
	status, err := database.GetChainStatus(a.db)
	if err != nil {
		log.Printf("Error getting chain status: %v", err)
		// Handle the error appropriately, maybe send a message to the user or return early
		return
	}
	presetSettingsMsg := fmt.Sprintf(`
	Settings > Presets (🔗%s)

This menu allows you to specify defaults for buy and sell settings`, status)

	buttons := makePresetButtons()

	a.client.SendMessage(m.Chat.ID, presetSettingsMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) defaultPresetAutoBuyHandler(m *tbot.Message) {
	presetSettingsMsg := `
	Settings > Auto Buy

Select a chain below:
	`
	buttons := makeAutoBuyChainButtons()

	a.client.SendMessage(m.Chat.ID, presetSettingsMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) chainSettingsHandler(m *tbot.Message, status ChainCurrentStatus) {
	chainSettingsMsg := `
	Settings > Chains

Select the chain you'd like to use. You can only have one chain 
selected at the same time. Your defaults and presets will be 
different for each chain.
`
	buttons := makeChainSettingsButtons(status)
	a.client.SendMessage(m.Chat.ID, chainSettingsMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) copyTradingHandler(m *tbot.Message) {
	copyTradingMsg := `Copy Trading

Select the chain you want to copy:`

	buttons := makeCopyTradingChainButtons()

	a.client.SendMessage(m.Chat.ID, copyTradingMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) handleWalletAmountToTransfer(walletIndex int, m *tbot.Message) {
	wallets, err := database.GetAllWallets(a.db)
	if err != nil {
		log.Printf("Error retrieving wallets: %v", err)
		a.client.SendMessage(m.Chat.ID, "Failed to fetch wallets.")
		return
	}

	if walletIndex >= len(wallets) || walletIndex < 0 {
		a.client.SendMessage(m.Chat.ID, "Invalid wallet selection.")
		return
	}

	selectedWallet := wallets[walletIndex]
	selectedToWallet := wallets[walletIndex+1]

	walletMsg := fmt.Sprintf(`Settings > Transfers (🔗ETH)
From: %s
To: %s
Available balance: 0Ξ
Do not forget to leave some amount for gas!

Enter the amount you would like to transfer:
`, selectedWallet.Address, selectedToWallet.Address)

	app.client.SendMessage(m.Chat.ID, walletMsg)
}

func convertToChainStatus(chainStatus database.ChainStatusOne) ChainStatus {
	switch chainStatus {
	case database.Ethereum:
		return ChainStatus{Name: "Ethereum"}
	case database.BSC:
		return ChainStatus{Name: "BSC"}
	case database.Blast:
		return ChainStatus{Name: "Blast"}
	case database.Base:
		return ChainStatus{Name: "Base"}
	case database.Avax:
		return ChainStatus{Name: "Avax"}
	case database.Solana:
		return ChainStatus{Name: "Solana"}
	default:
		return ChainStatus{Name: "Unknown"} // Or however you wish to handle unknown cases
	}
}

func (a *application) defaultPresetBuyHandler(m *tbot.Message) {
	currentChainStatus, err := database.GetChainStatus(a.db)
	if err != nil {
		// Handle the error appropriately. For example, log it or return it from the function.
		log.Printf("Error getting chain status: %v", err)
		return // Or handle the error in another way suitable for your application's flow.
	}
	buyPresetMsg := fmt.Sprintf(`
	Settings > Presets > Buy Value (🔗%s)

Click on the buttons below to specify new buy button values
	`, currentChainStatus)
	inlineKeyboard := makePresetBuyButtons()

	a.client.SendMessage(m.Chat.ID, buyPresetMsg, tbot.OptInlineKeyboardMarkup(inlineKeyboard))
}

func (a *application) languageHandler(m *tbot.Message) {
	languageMsg := `
	Settings > Language

Select a preferred language below:
	`

	languageButtons := makeLanugageButtons()

	a.client.SendMessage(m.Chat.ID, languageMsg, tbot.OptInlineKeyboardMarkup(languageButtons))
}

func (a *application) autoBuyEthereumHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗ETH)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoEthereumButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) autoBuyBSCHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗BSC)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoBSCButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) autoBuyBLAHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗Blast)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoBLAButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) autoBuyBASHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗Blast)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoBASButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) autoBuyAVAHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗Blast)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoAVAButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) autoBuySOLHandler(m *tbot.Message) {
	autoBuyMsg := `Settings > Auto Buy (🔗Solana)

Auto Buy will automatically purchase any contract address that
sent to the bot when this setting is on. Configure the settings 
below to control how to buy the token.`

	buttons := makeAutoSOLButtons()

	a.client.SendMessage(m.Chat.ID, autoBuyMsg, tbot.OptInlineKeyboardMarkup(buttons))

}

func (a *application) tradeConfirmHandler(m *tbot.Message) {
	tradeConfirmMsg := `Settings > Trade Confirmation

When on, this option will show an additional confirmation dialog 
before tokens are bought or sold. This helps to avoid accidental 
transactions.`

	buttons := makeTradeConfirmButtons()

	a.client.SendMessage(m.Chat.ID, tradeConfirmMsg, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	switch cq.Data {
	case "auto_sniper":
		a.sendCodeSnippet(cq.Message)
		// Existing cases...
	case "manual_buyer":
		// Existing cases...
	case "positions_management":
		// Existing cases...
	case "copy_trading":
		a.copyTradingHandler(cq.Message)
		// Existing cases...
	case "pending_orders":
		// Existing cases...
	case "settings_":
		currentChainStatus, err := database.GetChainStatus(a.db)
		if err != nil {
			log.Printf("Error getting chain status: %v", err)
			// Optionally, send a message to the user indicating an error occurred.
			a.client.SendMessage(cq.Message.Chat.ID, "Failed to get chain status.")
			break // Or continue, depending on your desired behavior
		}
		convertedChainStatus := convertToChainStatus(currentChainStatus)
		a.settingsHandler(cq.Message, convertedChainStatus)
		// Existing cases...
	case "war_room":
		a.warRoomHandler(cq.Message)
		// Existing cases...
	case "backup_bots":
		// Existing cases...

	case "bridge_":
		a.handleBridge(cq.Message)
	case "languages_":
		a.languageHandler(cq.Message)
		// Existing cases...
	case "wallets_":
		a.walletHandler(cq.Message)
	case "close_":
		// Existing cases...
	case "create_new_wallet":
		// Handle the creation of a new wallet
		a.createWalletHandler(cq.Message)
	case "import_other_wallets":
		a.waitingForKey = true
		a.walletHandler(cq.Message)
		a.client.SendMessage(cq.Message.Chat.ID, "Please enter your private key:")

	case "gas_preset_buttons":
		a.gasPresetHandler(cq.Message)

	case "remove_wallets":
		a.waitingForRemove = true
		a.walletHandler(cq.Message)
		a.client.SendMessage(cq.Message.Chat.ID, "Please enter your address:")

	case "back_to_mainboard":
		a.startHandler(cq.Message)
	case "private_keys_management":
		a.privateKeyHandler(cq.Message)

	case "rearrange_wallets":
		a.waitingForRearrange = true
		a.walletHandler(cq.Message)
		a.client.SendMessage(cq.Message.Chat.ID, "Please enter your address:")

	case "transfer_crypto":
		a.walletHandler(cq.Message)
		a.handleTransferCrypto(cq.Message)

	case "transfer_back":
		a.walletHandler(cq.Message)

	case "select_bridge":
		a.walletHandler(cq.Message)
		a.handleBridge(cq.Message)

	case "default_wallet":
		a.defaultWalletHandler(cq.Message)

	case "balance_transfer":
		a.selectFromWalletHandler(cq.Message)

	case "token_transfer":
		a.selectFromWalletHandler(cq.Message)

	case "snipe_wallets":
		a.snipeWalletsSelectHandler(cq.Message)

	case "manual_buy_wallets":
		a.snipeWalletsSelectHandler(cq.Message)

	case "ethereum_bridge":
		a.handleTokenTransfer(cq.Message)

	case "cancel_bridge":
		a.walletHandler(cq.Message)

	case "wallet_settings":
		a.walletSettingsHandler(cq.Message)

	case "chains_settings":
		status, err := database.GetChainStatus(a.db)
		if err != nil {
			log.Printf("Error getting chain status: %v", err)
			// Optionally, send a message to the user or handle the error in another appropriate way.
			a.client.SendMessage(cq.Message.Chat.ID, "Failed to get chain status.")
			return
		}

		a.chainSettingsHandler(cq.Message, ChainCurrentStatus(status))

	case "change_referal_wallet":
		a.waitingChangeReferalWallet = true
		walletMsg := `
		Paste the public address of the new wallet for your referral earnings.

Heads up! You will only be able to change the wallet again after 48 hours!
		`
		a.client.SendMessage(cq.Message.Chat.ID, walletMsg)

	case "refer_and_earn":
		a.waitingForReferAndEarn = true
		walletMsg := `Refer & Earn
Refer other users to earn commissions on their trades!

Enter your wallet address to earn:`
		a.client.SendMessage(cq.Message.Chat.ID, walletMsg)
		//a.referAndEarnHandler(cq.Message)

	case "bsc_chain_settings":
		bscStatus := ChainStatus{Name: "BSC"}
		a.settingsHandler(cq.Message, bscStatus)

	case "ehteruem_chain_settings":
		ethereumStatus := ChainStatus{Name: "Ethereum"}
		chainStatus := convertChainStatusNameToType(ethereumStatus.Name)
		err := database.SaveChainStatus(a.db, chainStatus)
		if err != nil {
			log.Printf("Error saving chain status: %v", err)
			// Handle error, possibly send a message to the user or log it
		}
		a.settingsHandler(cq.Message, ethereumStatus)

	case "blast_chain_settings":
		blastStatus := ChainStatus{Name: "Blast"}
		chainStatus := convertChainStatusNameToType(blastStatus.Name)
		err := database.SaveChainStatus(a.db, chainStatus)
		if err != nil {
			log.Printf("Error saving chain status: %v", err)
			// Handle error, possibly send a message to the user or log it
		}
		a.settingsHandler(cq.Message, blastStatus)

	case "base_chain_settings":
		baseStatus := ChainStatus{Name: "Base"}
		chainStatus := convertChainStatusNameToType(baseStatus.Name)
		err := database.SaveChainStatus(a.db, chainStatus)
		if err != nil {
			log.Printf("Error saving chain status: %v", err)
			// Handle error, possibly send a message to the user or log it
		}
		a.settingsHandler(cq.Message, baseStatus)

	case "avax_chain_settings":
		avaxStatus := ChainStatus{Name: "Avax"}
		chainStatus := convertChainStatusNameToType(avaxStatus.Name)
		err := database.SaveChainStatus(a.db, chainStatus)
		if err != nil {
			log.Printf("Error saving chain status: %v", err)
			// Handle error, possibly send a message to the user or log it
		}
		a.settingsHandler(cq.Message, avaxStatus)

	case "solana_chain_settings":
		solanaStatus := ChainStatus{Name: "Solana"}
		chainStatus := convertChainStatusNameToType(solanaStatus.Name)
		err := database.SaveChainStatus(a.db, chainStatus)
		if err != nil {
			log.Printf("Error saving chain status: %v", err)
			// Handle error, possibly send a message to the user or log it
		}
		a.settingsHandler(cq.Message, solanaStatus)

	case "presets_settings":
		a.presetSettingsHander(cq.Message)

	case "default_settings":
		a.defaultSettingsHandler(cq.Message)

	case "buy_preset_buttons":
		a.defaultPresetBuyHandler(cq.Message)

	case "auto_buy_buttons":
		a.defaultPresetAutoBuyHandler(cq.Message)

	case "auto_buy_ehteruem_chain_settings":
		a.autoBuyEthereumHandler(cq.Message)

	case "auto_buy_bsc_chain_settings":
		a.autoBuyBSCHandler(cq.Message)

	case "auto_buy_blast_chain_settings":
		a.autoBuyBLAHandler(cq.Message)

	case "auto_buy_base_chain_settings":
		a.autoBuyBASHandler(cq.Message)

	case "auto_buy_avax_chain_settings":
		a.autoBuyAVAHandler(cq.Message)

	case "auto_buy_solana_chain_settings":
		a.autoBuySOLHandler(cq.Message)

	case "trade_confirm_buttons":
		a.tradeConfirmHandler(cq.Message)

	default:
		if strings.HasPrefix(cq.Data, "wallet_select_") {
			walletIndexStr := strings.TrimPrefix(cq.Data, "wallet_select_")
			walletIndex, err := strconv.Atoi(walletIndexStr) // Handle both returned values
			if err != nil {
				log.Printf("Parse error: %v", err)
				return
			}
			log.Printf("Index: %d", walletIndex) // Confirm index parsing
			a.UpdateWalletSelection(a.db, walletIndex, cq.Message)
		} else if strings.HasPrefix(cq.Data, "wallet_from_select_") {
			walletIndexStr := strings.TrimPrefix(cq.Data, "wallet_from_select_")
			walletIndex, err := strconv.Atoi(walletIndexStr) // Convert the suffix to an integer index
			if err != nil {
				log.Printf("Parse error: %v", err)
				return
			}
			log.Printf("Selected wallet index: %d", walletIndex) // Confirm index parsing
			a.handleWalletSelectionForTransfer(walletIndex, cq.Message)
		} else if strings.HasPrefix(cq.Data, "wallet_to_select_") {
			walletIndexStr := strings.TrimPrefix(cq.Data, "wallet_to_select_")
			walletIndex, err := strconv.Atoi(walletIndexStr) // Convert the suffix to an integer index
			if err != nil {
				log.Printf("Parse error: %v", err)
				return
			}
			log.Printf("Selected wallet index: %d", walletIndex) // Confirm index parsing
			a.handleWalletAmountToTransfer(walletIndex, cq.Message)
		} else {
			log.Println("Unhandled callback data:", cq.Data)
		}

		// Delete the original message to clean up
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	}
}

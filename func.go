package main

import (

	// Add this line

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/l3njo/rochambeau/database"
	"github.com/yanzay/tbot/v2"
)

type WormholeRequest struct {
	FromChain   string `json:"fromChain"`
	ToChain     string `json:"toChain"`
	FromToken   string `json:"fromToken"`
	ToToken     string `json:"toToken"`
	Amount      string `json:"amount"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
}

type WormholeResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

func makeGeneralButtons() *tbot.InlineKeyboardMarkup {
	btnGroup1 := []tbot.InlineKeyboardButton{
		{Text: "🚀 Auto Sniper", CallbackData: "auto_sniper"},
		{Text: "🦽Manual Buyer", CallbackData: "manual_buyer"},
	}

	btnGroup2 := []tbot.InlineKeyboardButton{
		{Text: "🏦Positions", CallbackData: "positions_management"},
		{Text: "🏛️Copy Trading", CallbackData: "copy_trading"},
	}

	btnGroup3 := []tbot.InlineKeyboardButton{
		{Text: "🔀Pending orders", CallbackData: "pending_orders"},
		{Text: "⚙️Settings", CallbackData: "settings_"},
	}

	btnGroup4 := []tbot.InlineKeyboardButton{
		{Text: "🦸‍♂️Refer and Earn", CallbackData: "refer_and_earn"},
		{Text: "⚔️War room", CallbackData: "war_room"},
	}

	btnGroup5 := []tbot.InlineKeyboardButton{
		{Text: "♻️Backup bots", CallbackData: "backup_bots"},
		{Text: "🧏‍♂️Languages", CallbackData: "languages_"},
	}

	btnGroup6 := []tbot.InlineKeyboardButton{
		{Text: "👛Wallets", CallbackData: "wallets_"},
		{Text: "🌉Bridge", CallbackData: "bridge_"},
	}

	btnClose := tbot.InlineKeyboardButton{
		Text:         "❌Close",
		CallbackData: "close_",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup1,
		btnGroup2,
		btnGroup3,
		btnGroup4,
		btnGroup5,
		btnGroup6,
		{btnClose},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeWalletButtons() *tbot.InlineKeyboardMarkup {
	btnGroup1 := []tbot.InlineKeyboardButton{
		{Text: "🦸Create", CallbackData: "create_new_wallet"},
		{Text: "🤟Import", CallbackData: "import_other_wallets"},
	}

	btnGroup2 := []tbot.InlineKeyboardButton{
		{Text: "🏢Rearrange", CallbackData: "rearrange_wallets"},
		{Text: "🚮Remove", CallbackData: "remove_wallets"},
	}

	btnGroup3 := []tbot.InlineKeyboardButton{
		{Text: "🌝Private Keys", CallbackData: "private_keys_management"},
		{Text: "🛀Default wallet", CallbackData: "default_wallet"},
	}

	btnGroup4 := []tbot.InlineKeyboardButton{
		{Text: "🚆Transfers", CallbackData: "transfer_crypto"},
		{Text: "🌉Bridge", CallbackData: "select_bridge"},
	}

	btnBack := tbot.InlineKeyboardButton{
		Text:         "Back",
		CallbackData: "back_to_mainboard",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup1,
		btnGroup2,
		btnGroup3,
		btnGroup4,
		{btnBack},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeWalletSettingsButtons() *tbot.InlineKeyboardMarkup {
	btnGroup1 := []tbot.InlineKeyboardButton{
		{Text: "🦸Create", CallbackData: "create_new_wallet"},
		{Text: "🤟Import", CallbackData: "import_other_wallets"},
	}

	btnGroup2 := []tbot.InlineKeyboardButton{
		{Text: "🏢Rearrange", CallbackData: "rearrange_wallets"},
		{Text: "🚮Remove", CallbackData: "remove_wallets"},
	}

	btnGroup3 := []tbot.InlineKeyboardButton{
		{Text: "🌝Private Keys", CallbackData: "private_keys_management"},
		{Text: "🛀Default wallet", CallbackData: "default_wallet"},
	}

	btnGroup4 := []tbot.InlineKeyboardButton{
		{Text: "🚆Transfers", CallbackData: "transfer_crypto"},
		{Text: "🌉Bridge", CallbackData: "select_bridge"},
	}

	btnBack := tbot.InlineKeyboardButton{
		Text:         "Back",
		CallbackData: "settings_",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup1,
		btnGroup2,
		btnGroup3,
		btnGroup4,
		{btnBack},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeTransferButtons() *tbot.InlineKeyboardMarkup {
	balanceTranferBtn := tbot.InlineKeyboardButton{Text: "🗨️Balance Transfer", CallbackData: "balance_transfer"}
	tokenTranferBtn := tbot.InlineKeyboardButton{Text: "🏧Token Transfer", CallbackData: "token_transfer"}
	cancelTransferBtn := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "transfer_back"}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{{balanceTranferBtn}, {tokenTranferBtn}, {cancelTransferBtn}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeBridgeButtons() *tbot.InlineKeyboardMarkup {
	ethereumBridge := tbot.InlineKeyboardButton{Text: "Ethereum", CallbackData: "ethereum_bridge"}
	avaxBridge := tbot.InlineKeyboardButton{Text: "Avax", CallbackData: "avax_bridge"}
	bscBridge := tbot.InlineKeyboardButton{Text: "Bsc", CallbackData: "bsc_bridge"}
	baseBridge := tbot.InlineKeyboardButton{Text: "Base", CallbackData: "base_bridge"}
	cancelBridge := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "cancel_bridge"}

	buttons := [][]tbot.InlineKeyboardButton{{ethereumBridge}, {avaxBridge}, {baseBridge}, {bscBridge}, {cancelBridge}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeDefaultButtons() *tbot.InlineKeyboardMarkup {
	snipeWallets := tbot.InlineKeyboardButton{Text: "Snipe Wallets", CallbackData: "snipe_wallets"}
	manualBuyWallets := tbot.InlineKeyboardButton{Text: "Manual Buy Wallets", CallbackData: "manual_buy_wallets"}
	cancelBridge := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "cancel_default"}

	buttons := [][]tbot.InlineKeyboardButton{{snipeWallets}, {manualBuyWallets}, {cancelBridge}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func transferTokensViaWormhole(request WormholeRequest) (string, error) {
	// Convert the request struct to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Replace with the actual Wormhole API endpoint
	apiUrl := "https://api.wormholebridge.com/v1/transfer"

	// Create a new request using http
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response
	var wormholeResponse WormholeResponse
	err = json.Unmarshal(body, &wormholeResponse)
	if err != nil {
		return "", err
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to transfer tokens: %s", resp.Status)
	}

	return wormholeResponse.TransactionID, nil
}

func makeReferButtons() *tbot.InlineKeyboardMarkup {
	btnGroup := []tbot.InlineKeyboardButton{
		{Text: "Change Referal Wallet", CallbackData: "change_referal_wallet"},
		{Text: "Dismiss Message", CallbackData: "dismiss_message"},
	}

	buttons := [][]tbot.InlineKeyboardButton{btnGroup}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeLanugageButtons() *tbot.InlineKeyboardMarkup {
	languageEnglish := tbot.InlineKeyboardButton{Text: "English", CallbackData: "language_english"}
	languageFrance := tbot.InlineKeyboardButton{Text: "France", CallbackData: "language_france"}
	languageChina := tbot.InlineKeyboardButton{Text: "China", CallbackData: "language_china"}
	languageSpain := tbot.InlineKeyboardButton{Text: "Spain", CallbackData: "language_spain"}
	cancelLanguage := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "cancel_language"}

	buttons := [][]tbot.InlineKeyboardButton{{languageEnglish}, {languageFrance}, {languageChina}, {languageSpain}, {cancelLanguage}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeSettingsButtons() *tbot.InlineKeyboardMarkup {
	walletsButton := tbot.InlineKeyboardButton{Text: "👛Wallets", CallbackData: "wallet_settings"}
	presetsettingsn := tbot.InlineKeyboardButton{Text: "⏸️Presets", CallbackData: "presets_settings"}
	defaulsettingsn := tbot.InlineKeyboardButton{Text: "🛌Default", CallbackData: "default_settings"}
	chainssettings := tbot.InlineKeyboardButton{Text: "🔗Chains", CallbackData: "chains_settings"}
	cancelButton := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "setting_cancel"}

	buttons := [][]tbot.InlineKeyboardButton{{walletsButton}, {presetsettingsn}, {defaulsettingsn}, {chainssettings}, {cancelButton}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeChainSettingsButtons(cs ChainCurrentStatus) *tbot.InlineKeyboardMarkup {
	ethereumButton := tbot.InlineKeyboardButton{Text: "Ethereum", CallbackData: "ehteruem_chain_settings"}
	bscButton := tbot.InlineKeyboardButton{Text: "BSC", CallbackData: "bsc_chain_settings"}
	blastButton := tbot.InlineKeyboardButton{Text: "Blast", CallbackData: "blast_chain_settings"}
	baseButton := tbot.InlineKeyboardButton{Text: "Base", CallbackData: "base_chain_settings"}
	avaxButton := tbot.InlineKeyboardButton{Text: "Avax", CallbackData: "avax_chain_settings"}
	solanaButton := tbot.InlineKeyboardButton{Text: "Solana", CallbackData: "solana_chain_settings"}
	cancelButton := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "chain_setting_cancel"}
	switch cs {
	case Ethereum:
		ethereumButton.Text = "✅" + ethereumButton.Text
	case BSC:
		bscButton.Text = "✅" + bscButton.Text
	case Blast:
		blastButton.Text = "✅" + blastButton.Text
	case Avax:
		avaxButton.Text = "✅" + avaxButton.Text
	case Solana:
		solanaButton.Text = "✅" + solanaButton.Text

	default:
		ethereumButton.Text = "✅" + ethereumButton.Text

	}

	buttons := [][]tbot.InlineKeyboardButton{
		{ethereumButton},
		{bscButton},
		{blastButton},
		{baseButton},
		{avaxButton},
		{solanaButton},
		{cancelButton},
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makePresetButtons() *tbot.InlineKeyboardMarkup {
	gasButtons := tbot.InlineKeyboardButton{Text: "⛽Gas Buttons", CallbackData: "gas_preset_buttons"}
	buyButtons := tbot.InlineKeyboardButton{Text: "💲Buy Buttons", CallbackData: "buy_preset_buttons"}
	autoBuyButtons := tbot.InlineKeyboardButton{Text: "🤖Auto Buy", CallbackData: "auto_buy_buttons"}
	tradeConfirmButtons := tbot.InlineKeyboardButton{Text: "™️Trade Confirmation", CallbackData: "trade_confirm_buttons"}
	backPresetButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "back_preset_buttons"}

	buttons := [][]tbot.InlineKeyboardButton{{gasButtons}, {buyButtons}, {autoBuyButtons}, {tradeConfirmButtons}, {backPresetButton}}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeGasPresetButtons() *tbot.InlineKeyboardMarkup {
	btnGroup := []tbot.InlineKeyboardButton{
		{Text: "10.00", CallbackData: "new_gas_values1"},
		{Text: "45.00", CallbackData: "new_gas_values2"},
		{Text: "50.00", CallbackData: "new_gas_values3"},
	}

	btnDone := tbot.InlineKeyboardButton{
		Text:         "✅Done",
		CallbackData: "confirm_gas_fee",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup,
		{btnDone},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

/*
	newSettings := models.DefaultSettings{
		Slippage:        10,
		SellGweiExtra:   7.00,
		ApproveGwei:     7.00,
		BuyTax:          100.00,
		SellTax:         100.00,
		MinLiquidity:    150,
		AlphaMode:       false,
		MultitxOrRevert: false,
		AntiRug:         false,
		Createdate:      time.Now(),
		UpdatedAt:       time.Now(),
	}

result := database.SetDefaultSettings(app.db, &newSettings)

	if result != nil {
		log.Printf("Failed to set default settings: %v", result)
	}
*/
func makeDefaultSettingsButtons() *tbot.InlineKeyboardMarkup {
	newSettings, err := database.GetDefaultSettings(app.db)
	if err != nil {
		log.Printf("Failed to get default settings: %v", err)
		// Return nil or an error indicating failure to retrieve settings
		return nil
	}
	if len(newSettings) == 0 {
		// Handle the case where no settings are found
		// For example, you could return a default keyboard or an error
		log.Println("No default settings found.")
		// Return a default keyboard or nil with an error
		return &tbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]tbot.InlineKeyboardButton{},
		}
	}

	slippageBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Slippage: %d", newSettings[0].Slippage),
		CallbackData: "slip_page_settings",
	}

	sellGweiBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Sell Gwei Extra: %f", newSettings[0].SellGweiExtra),
		CallbackData: "sell_gwei_settings",
	}

	approveGweiBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Approve Gwei: %f", newSettings[0].ApproveGwei),
		CallbackData: "approve_gwei_settings",
	}

	buyTaxBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Buy Tax: %f", newSettings[0].BuyTax),
		CallbackData: "buy_tax_settings",
	}

	sellTaxBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Sell Tax: %f", newSettings[0].SellTax),
		CallbackData: "sell_tax_settings",
	}

	minLiquidityBtn := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Min Liquidity: %d", newSettings[0].MinLiquidity),
		CallbackData: "min_liquidity_settings",
	}

	alphaModeBtn := tbot.InlineKeyboardButton{
		Text:         "Alpha Mode:🔴",
		CallbackData: "alpha_mode_settings",
	}

	maxTxOrRevertBtn := tbot.InlineKeyboardButton{
		Text:         "MaxTx or Revert:🔴",
		CallbackData: "max_tx_or_revert_settings",
	}

	antiRugBtn := tbot.InlineKeyboardButton{
		Text:         "AntiRug:🔴",
		CallbackData: "anti_rug_settings",
	}

	backDefaultButton := tbot.InlineKeyboardButton{
		Text:         "Back",
		CallbackData: "back_default_settings",
	}

	buttons := [][]tbot.InlineKeyboardButton{
		{slippageBtn},
		{sellGweiBtn},
		{approveGweiBtn},
		{buyTaxBtn},
		{sellTaxBtn},
		{minLiquidityBtn},
		{alphaModeBtn},
		{maxTxOrRevertBtn},
		{antiRugBtn},
		{backDefaultButton},
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makePresetBuyButtons() *tbot.InlineKeyboardMarkup {
	btnGroup := []tbot.InlineKeyboardButton{
		{Text: "0.1", CallbackData: "new_buy_values1"},
		{Text: "0.2", CallbackData: "new_buy_values2"},
		{Text: "0.8", CallbackData: "new_buy_values3"},
		{Text: "1", CallbackData: "new_buy_value4"},
	}

	btnDone := tbot.InlineKeyboardButton{
		Text:         "✅Done",
		CallbackData: "confirm_buy_fee",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup,
		{btnDone},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoBuyChainButtons() *tbot.InlineKeyboardMarkup {
	ethereumButton := tbot.InlineKeyboardButton{Text: "Ethereum", CallbackData: "auto_buy_ehteruem_chain_settings"}
	bscButton := tbot.InlineKeyboardButton{Text: "BSC", CallbackData: "auto_buy_bsc_chain_settings"}
	blastButton := tbot.InlineKeyboardButton{Text: "Blast", CallbackData: "auto_buy_blast_chain_settings"}
	baseButton := tbot.InlineKeyboardButton{Text: "Base", CallbackData: "auto_buy_base_chain_settings"}
	avaxButton := tbot.InlineKeyboardButton{Text: "Avax", CallbackData: "auto_buy_avax_chain_settings"}
	solanaButton := tbot.InlineKeyboardButton{Text: "Solana", CallbackData: "auto_buy_solana_chain_settings"}
	cancelButton := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "auto_buy_chain_setting_cancel"}
	buttons := [][]tbot.InlineKeyboardButton{
		{ethereumButton},
		{bscButton},
		{blastButton},
		{baseButton},
		{avaxButton},
		{solanaButton},
		{cancelButton},
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoEthereumButtons() *tbot.InlineKeyboardMarkup {
	autoBuyETHSettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_buy_settings"}
	autoBuyETHBackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuyETHSettingsButton},
		{autoBuyETHBackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoBSCButtons() *tbot.InlineKeyboardMarkup {
	autoBuyBSCSettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_BSC_buy_settings"}
	autoBuyBSCBackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_BSC_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuyBSCSettingsButton},
		{autoBuyBSCBackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoBLAButtons() *tbot.InlineKeyboardMarkup {
	autoBuyBLASettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_BLA_buy_settings"}
	autoBuyBLABackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_BLA_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuyBLASettingsButton},
		{autoBuyBLABackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoBASButtons() *tbot.InlineKeyboardMarkup {
	autoBuyBASSettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_BAS_buy_settings"}
	autoBuyBASBackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_BAS_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuyBASSettingsButton},
		{autoBuyBASBackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoAVAButtons() *tbot.InlineKeyboardMarkup {
	autoBuyAVASettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_AVA_buy_settings"}
	autoBuyAVABackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_AVA_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuyAVASettingsButton},
		{autoBuyAVABackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeAutoSOLButtons() *tbot.InlineKeyboardMarkup {
	autoBuySOLSettingsButton := tbot.InlineKeyboardButton{Text: "Auto Buy:Off", CallbackData: "auto_SOL_buy_settings"}
	autoBuySOLBackButton := tbot.InlineKeyboardButton{Text: "Back", CallbackData: "auto_SOL_buy_back"}

	buttons := [][]tbot.InlineKeyboardButton{
		{autoBuySOLSettingsButton},
		{autoBuySOLBackButton},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeTradeConfirmButtons() *tbot.InlineKeyboardMarkup {
	btnGroup := []tbot.InlineKeyboardButton{
		{Text: "On", CallbackData: "new_trade_on"},
		{Text: "Off🟢", CallbackData: "new_trade_off"},
	}

	btnCancel := tbot.InlineKeyboardButton{
		Text:         "Cancel",
		CallbackData: "cancel_trade_confirm",
	}

	// Correctly appending all button groups and close button
	buttons := [][]tbot.InlineKeyboardButton{
		btnGroup,
		{btnCancel},
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func makeCopyTradingChainButtons() *tbot.InlineKeyboardMarkup {
	ethereumButton := tbot.InlineKeyboardButton{Text: "ETH", CallbackData: "copy_trading_ehteruem_chain_settings"}
	bscButton := tbot.InlineKeyboardButton{Text: "BSC", CallbackData: "copy_trading_bsc_chain_settings"}
	blastButton := tbot.InlineKeyboardButton{Text: "BLAST", CallbackData: "copy_trading_blast_chain_settings"}
	baseButton := tbot.InlineKeyboardButton{Text: "BASE", CallbackData: "copy_trading_base_chain_settings"}
	avaxButton := tbot.InlineKeyboardButton{Text: "AVAX", CallbackData: "copy_trading_avax_chain_settings"}
	solanaButton := tbot.InlineKeyboardButton{Text: "SOLANA", CallbackData: "copy_trading_solana_chain_settings"}
	cancelButton := tbot.InlineKeyboardButton{Text: "Cancel", CallbackData: "copy_trading_chain_setting_cancel"}
	buttons := [][]tbot.InlineKeyboardButton{
		{ethereumButton},
		{bscButton},
		{blastButton},
		{baseButton},
		{avaxButton},
		{solanaButton},
		{cancelButton},
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

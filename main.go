package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/l3njo/rochambeau/database"
	"github.com/l3njo/rochambeau/models"
	"github.com/yanzay/tbot/v2"
)

type Wallet struct {
	ChainScanLabel string
	AccountWorth   int
	PrivateKey     string
	Address        string
}

type UserResponse struct {
	User struct {
		VerifiedCredentials []struct {
			Address          string                 `json:"address"`
			Chain            string                 `json:"chain"`
			ID               string                 `json:"id"`
			NameService      map[string]interface{} `json:"name_service"` // Adjusted to map[string]interface{}
			PublicIdentifier string                 `json:"public_identifier"`
			WalletName       string                 `json:"wallet_name"`
			WalletProvider   string                 `json:"wallet_provider"`
			Properties       map[string]interface{} `json:"wallet_properties"` // Assuming you might want to handle properties similarly
			Format           string                 `json:"format"`
			LastSelectedAt   string                 `json:"lastSelectedAt"`
		} `json:"verifiedCredentials"`
		Wallets []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Chain     string `json:"chain"`
			PublicKey string `json:"publicKey"`
			Provider  string `json:"provider"`
		} `json:"wallets"`
	}
}

type State int

const (
	WaitingForSenderAddress State = iota
	WaitingForRecipientAddress
)

type ChainCurrentStatus string

const (
	Ethereum ChainCurrentStatus = "Ethereum"
	BSC      ChainCurrentStatus = "BSC"
	Blast    ChainCurrentStatus = "Blast"
	Base     ChainCurrentStatus = "Base"
	Avax     ChainCurrentStatus = "Avax"
	Solana   ChainCurrentStatus = "Solana"
)

type application struct {
	waitingForRemove           bool
	waitingForKey              bool
	waitingForRearrange        bool
	messageChannel             chan *tbot.Message
	client                     *tbot.Client
	selectedWalletIndices      map[int]bool
	allWallets                 []*Wallet
	waitingForReferAndEarn     bool
	waitingChangeReferalWallet bool
	userLanguage               map[int]string
	//balanceMsg     []models.Wallet
	db *sql.DB
}

var (
	app   application
	bot   *tbot.Server
	token string
)

func init() {
	app.db = database.ConnectDatabase()
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}
	token = os.Getenv("TELEGRAM_TOKEN")
	bot = tbot.New(token)
	app.userLanguage = make(map[int]string)
	app.client = bot.Client()
	if app.client == nil {
		log.Fatal("Failed to initialize Telegram client")
	}
	app.messageChannel = make(chan *tbot.Message)
}

func main() {
	bot.HandleMessage("/start", app.startHandler)
	bot.HandleCallback(app.callbackHandler)

	bot.HandleMessage("", func(m *tbot.Message) {

		if app.waitingForKey {
			app.waitingForKey = false
			privateKey := m.Text
			// Attempt to load the wallet using the provided private key.
			wallet, err := createOtherWallet(privateKey)
			if err != nil {
				app.client.SendMessage(m.Chat.ID, "Failed to load wallet. Please check the private key.")
				log.Printf("Error loading wallet: %v", err)
			} else {
				currentChainStatus, err := database.GetChainStatus(app.db)
				if err != nil {
					log.Printf("Error getting chain status: %v", err)
					// Handle the error appropriately, maybe send a message to the user or return early
					return
				}
				newWallet := &models.Wallet{ // Assuming Wallet1 is the correct type to use here
					ChainScanLabel: "Balance",
					AccountWorth:   0,
					PrivateKey:     wallet.PrivateKey,
					Address:        wallet.Address,
				}
				// Correctly append newWallet to app.balanceMsg

				database.CreateWallet(app.db, newWallet)
				getWallet, err := database.GetAllWallets(app.db)
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
			Settings > Wallets (🔗%s)
			
			Wallet Worth: $0
			
			Your currently added wallets:
			%s
			`, currentChainStatus, walletDetails.String())

				walletMsg += "\n"

				buttons := makeWalletButtons()
				app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))
			}
		}

		if app.waitingForRemove {
			currentChainStatus, err := database.GetChainStatus(app.db)
			if err != nil {
				log.Printf("Error getting chain status: %v", err)
				// Handle the error appropriately, maybe send a message to the user or return early
				return
			}
			app.waitingForRemove = false
			wallet_address := m.Text
			database.DeleteWallet(app.db, wallet_address)
			getRefreshWallet, err := database.GetAllWallets(app.db)
			if err != nil {
				// Handle error appropriately
				log.Printf("Error fetching wallet: %v", err)
				return
			}
			var walletDetails strings.Builder
			for i, wallet := range getRefreshWallet {
				walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", i+1, wallet.ChainScanLabel, wallet.AccountWorth, wallet.Address))
			}

			walletMsg := fmt.Sprintf(`
			Settings > Wallets (🔗%s)
			
			Wallet Worth: $0
			
			Your currently added wallets:
			%s
			`, currentChainStatus, walletDetails.String())

			walletMsg += "\n"

			buttons := makeWalletButtons()
			app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))

		}

		if app.waitingForRearrange {
			currentChainStatus, err := database.GetChainStatus(app.db)
			if err != nil {
				log.Printf("Error getting chain status: %v", err)
				// Handle the error appropriately, maybe send a message to the user or return early
				return
			}
			app.waitingForRearrange = false
			wallet_address := m.Text
			// database.UpdateWalletUpdatedAt(app.db, wallet_address)
			// var updatedAt time.Time
			// updatedAt, err := database.GetUpdatedTime(app.db, wallet_address)
			// if err != nil {
			// 	// Handle error appropriately
			// 	log.Printf("Error fetching wallet: %v", err)
			// 	return
			// }
			// database.RearrangeWalletUpdatedAt(app.db, updatedAt)
			getRefreshWallet, err := database.GetAllWallets(app.db)
			if err != nil {
				// Handle error appropriately
				log.Printf("Error fetching wallet: %v", err)
				return
			}
			var walletDetails strings.Builder
			var firstWallet *models.Wallet
			index := 1
			for _, wallet := range getRefreshWallet {
				if wallet.Address == wallet_address {
					firstWallet = wallet
				}
			}
			walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", index, firstWallet.ChainScanLabel, firstWallet.AccountWorth, firstWallet.Address))
			for _, wallet := range getRefreshWallet {
				if wallet.Address != wallet_address {
					walletDetails.WriteString(fmt.Sprintf("%d: %s: $%d\n%s\n", index+1, wallet.ChainScanLabel, wallet.AccountWorth, wallet.Address))
					index++
				}
			}

			walletMsg := fmt.Sprintf(`
			Settings > Wallets (🔗%s)
			
			Wallet Worth: $0
			
			Your currently added wallets:
			%s
			`, currentChainStatus, walletDetails.String())

			walletMsg += "\n"

			buttons := makeWalletButtons()
			app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))

		}

		if app.waitingForReferAndEarn {
			app.waitingForReferAndEarn = false
			walletAddress := m.Text
			multiple, _ := database.FindMultipleWalletsByAddress(app.db, walletAddress)

			if !multiple {
				ErrMsg := "This address does not exist in your wallet. Please correct address or import."
				app.client.SendMessage(m.Chat.ID, ErrMsg)
			} else if multiple {
				log.Printf("Heelo world!")
				walletMsg := fmt.Sprintf(`
Start earning today! 🚀

Share the link below with fellow traders and start earning 20%% commission on each trade and bonus points.

Commissions Wallet 👇
%s

Your personal link 👇
https://t.me/Sigma_buyot?start=ref=7090525195

Stats:
🤝 Referred: 0 
📊 Volume: 0.000Ξ 
💰 Revenue: 0.000Ξ 
⚖️ Trades: 0`, walletAddress)

				buttons := makeReferButtons()

				app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(buttons))

			}

		}

		if app.waitingChangeReferalWallet {
			app.waitingChangeReferalWallet = false
			walletAddress := m.Text
			multiple, _ := database.FindMultipleWalletsByAddress(app.db, walletAddress)

			if !multiple {
				ErrMsg := "This address does not exist in your wallet. Please correct address or import."
				app.client.SendMessage(m.Chat.ID, ErrMsg)
			} else if multiple {
				log.Printf("Heelo world!")
				walletMsg := fmt.Sprintf(`
Start earning today! 🚀

Share the link below with fellow traders and start earning 20%% commission on each trade and bonus points.

Commissions Wallet 👇
%s

Your personal link 👇
https://t.me/Sigma_buyot?start=ref=7090525195

Stats:
🤝 Referred: 0 
📊 Volume: 0.000Ξ 
💰 Revenue: 0.000Ξ 
⚖️ Trades: 0`, walletAddress)

				button := tbot.InlineKeyboardButton{
					Text: "Dismiss", CallbackData: "button_clicked",
				}

				inlineKeyboardMarkup := tbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]tbot.InlineKeyboardButton{
						{button},
					},
				}
				app.client.SendMessage(m.Chat.ID, walletMsg, tbot.OptInlineKeyboardMarkup(&inlineKeyboardMarkup))

			}

		}

	})

	log.Fatal(bot.Start())
}

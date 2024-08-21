package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/l3njo/rochambeau/models"
)

var ErrRecordNotFound = errors.New("record not found")

type TransferRecord struct {
	ID            uuid.UUID `json:"id"`
	FromChain     string    `json:"fromChain"`
	ToChain       string    `json:"toChain"`
	FromToken     string    `json:"fromToken"`
	ToToken       string    `json:"toToken"`
	Amount        string    `json:"amount"`
	FromAddress   string    `json:"fromAddress"`
	ToAddress     string    `json:"toAddress"`
	TransactionID string    `json:"transactionId"`
	Status        string    `json:"status"`
}

// CreateWallet inserts a new wallet record into the database.
func CreateWallet(db *sql.DB, wallet *models.Wallet) error {

	query := `INSERT INTO wallet (chat_id, chain_scan_label, account_worth, private_key, wallet_address) VALUES ($1, $2, $3, $4, $5);`
	_, err := db.Exec(query, wallet.ChatId, wallet.ChainScanLabel, wallet.AccountWorth, wallet.PrivateKey, wallet.Address)
	if err != nil {
		log.Printf("Failed to insert wallet: %v", err)
		return err
	}
	return nil
}

// GetWallet retrieves a wallet record by its ID.
func GetWallet(db *sql.DB, id uuid.UUID) (*models.Wallet, error) {
	rows := db.QueryRow(`SELECT id, chat_id, chain_scan_label, account_worth, private_key, wallet_address, create_date, updated_at FROM wallet`)
	var wallet models.Wallet
	err := rows.Scan(&wallet.ID, &wallet.ChatId, &wallet.ChainScanLabel, &wallet.AccountWorth, &wallet.PrivateKey, &wallet.Address, &wallet.Createdate, &wallet.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		log.Printf("Failed to query wallet: %v", err)
		return nil, err
	}
	return &wallet, nil
}

// UpdateWallet updates an existing wallet record in the database.
func UpdateWallet(db *sql.DB, wallet *models.Wallet) error {
	query := `UPDATE wallet SET chat_id = $1, chain_scan_label = $2, account_worth = $3, private_key = $4, wallet_address = $5 create_date = $6 WHERE id = $7`
	result, err := db.Exec(query,
		wallet.ChatId, wallet.ChainScanLabel, wallet.AccountWorth, wallet.PrivateKey, wallet.Address, wallet.Createdate, wallet.ID,
	)
	if err != nil {
		log.Printf("Failed to update wallet: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

// DeleteWallet removes a wallet record from the database by its ID.
func DeleteWallet(db *sql.DB, wallet_address string) error {
	result, err := db.Exec(`DELETE FROM wallet WHERE wallet_address = $1`, wallet_address)
	if err != nil {
		log.Printf("Failed to delete wallet: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were deleted")
	}
	return nil
}

// In database.go

func GetAllWallets(db *sql.DB) ([]*models.Wallet, error) {
	rows, err := db.Query("SELECT * FROM wallet")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*models.Wallet
	for rows.Next() {
		wallet := &models.Wallet{}
		err := rows.Scan(&wallet.ID, &wallet.ChatId, &wallet.ChainScanLabel, &wallet.AccountWorth, &wallet.PrivateKey, &wallet.Address, &wallet.Createdate, &wallet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func UpdateWalletUpdatedAt(db *sql.DB, wallet_address string) error {
	query := `UPDATE wallet SET updated_at = $1 WHERE wallet_address = $2`
	_, err := db.Exec(query, time.Now(), wallet_address)
	if err != nil {
		return err
	}
	return nil
}

func RearrangeWalletUpdatedAt(db *sql.DB, updated_at time.Time) error {
	query := `SELECT * FROM wallet WHERE ROWNUM <= 1 ORDER BY updated_at DESC`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func GetUpdatedTime(db *sql.DB, wallet_address string) (time.Time, error) {
	var updatedAt time.Time
	err := db.QueryRow("SELECT updated_at FROM wallet WHERE wallet_address = $1", wallet_address).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows were found
			return time.Time{}, ErrRecordNotFound // Assuming ErrRecordNotFound is defined elsewhere
		}
		return time.Time{}, err
	}
	return updatedAt, nil
}

func CreateTransferRecord(db *sql.DB, transfer *TransferRecord) error {
	query := `INSERT INTO transfers (from_chain, to_chain, from_token, to_token, amount, from_address, to_address, transaction_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := db.Exec(query, transfer.FromChain, transfer.ToChain, transfer.FromToken, transfer.ToToken, transfer.Amount, transfer.FromAddress, transfer.ToAddress, transfer.TransactionID, transfer.Status)
	return err
}

func FindMultipleWalletsByAddress(db *sql.DB, address string) (bool, error) {
	// Use QueryRow to get a single row result
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM wallet WHERE wallet_address = $1", address).Scan(&count)
	if err != nil {
		log.Printf("Error querying wallet count by address: %v", err)
		return false, err
	}

	// Now, count holds the number of rows that match the condition
	if count >= 1 {
		return true, nil
	} else if count == 0 {
		return false, nil
	} else {
		// This should not happen unless there's an issue with the database or query
		return false, fmt.Errorf("unexpected count value: %d", count)
	}
}

func SetDefaultSettings(db *sql.DB, settings *models.DefaultSettings) error {
	// Prepare an SQL statement
	stmt, err := db.Prepare("UPDATE settings SET slippage=$1, sellgweextra=$2, approvecwei=$3, buytax=$4, selltax=$5, minliquidity=$6, alphamode=$7, multitxorrevert=$8, antirug=$9 WHERE id=1")
	if err != nil {
		return err
	}

	// Execute the statement with the settings values
	_, err = stmt.Exec(settings.Slippage, settings.SellGweiExtra, settings.ApproveGwei, settings.BuyTax, settings.SellTax, settings.MinLiquidity, settings.AlphaMode, settings.MultitxOrRevert, settings.AntiRug)
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultSettings(db *sql.DB) ([]*models.DefaultSettings, error) {
	rows, err := db.Query("SELECT * FROM default_settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defaultSettings []*models.DefaultSettings
	for rows.Next() {
		defaultSetting := &models.DefaultSettings{}
		err := rows.Scan(&defaultSetting.ID, &defaultSetting.Slippage, &defaultSetting.SellGweiExtra, &defaultSetting.ApproveGwei, &defaultSetting.BuyTax, &defaultSetting.SellTax, &defaultSetting.MinLiquidity, &defaultSetting.AlphaMode, &defaultSetting.MultitxOrRevert, &defaultSetting.AntiRug, &defaultSetting.Createdate, &defaultSetting.UpdatedAt)
		if err != nil {
			return nil, err
		}
		defaultSettings = append(defaultSettings, defaultSetting)
	}

	return defaultSettings, nil
}

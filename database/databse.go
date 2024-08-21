package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ChainStatusOne string

const (
	Ethereum ChainStatusOne = "Ethereum"
	BSC      ChainStatusOne = "BSC"
	Blast    ChainStatusOne = "Blast"
	Base     ChainStatusOne = "Base"
	Avax     ChainStatusOne = "Avax"
	Solana   ChainStatusOne = "Solana"
)

func ConnectDatabase() *sql.DB {

	err := godotenv.Load() //by default, it is .env so we don't have to write
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	//we read our .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("USER")
	dbname := os.Getenv("DATABASE_NAME")
	pass := os.Getenv("PASSWORD")

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	} else {
		fmt.Println("Successfully connected to database!")
		return db
	}
}

func SaveChainStatus(db *sql.DB, status ChainStatusOne) error {
	var statusStr string
	switch status {
	case Ethereum:
		statusStr = "Ethereum"
	case BSC:
		statusStr = "BSC"
	case Blast:
		statusStr = "Blast"
	case Base:
		statusStr = "Base"
	case Avax:
		statusStr = "Avax"
	case Solana:
		statusStr = "Solana"
	default:
		return fmt.Errorf("unknown chain status")
	}

	_, err := db.Exec("INSERT INTO chain_status (status) VALUES ($1) ON CONFLICT (id) DO UPDATE SET status = $1 WHERE chain_status.id = EXCLUDED.id", statusStr)
	return err
}

// GetChainStatus retrieves the current chain status from the database
func GetChainStatus(db *sql.DB) (ChainStatusOne, error) {
	var statusStr string
	err := db.QueryRow("SELECT status FROM chain_status ORDER BY id desc LIMIT 1").Scan(&statusStr)
	if statusStr == "" {
		return Ethereum, nil
	}
	if err != nil {
		return Ethereum, err // Default to Ethereum if no status found
	}

	switch statusStr {
	case "Ethereum":
		return Ethereum, nil
	case "BSC":
		return BSC, nil
	case "Blast":
		return Blast, nil
	case "Base":
		return Base, nil
	case "Avax":
		return Avax, nil
	case "Solana":
		return Solana, nil
	default:
		return Ethereum, fmt.Errorf("unknown chain status")
	}
}

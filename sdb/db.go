package sdb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrInvalidCommand = errors.New("INVALID COMMAND")
var ErrNoTransaction = errors.New("NO TRANSACTION")
var ErrNoValue = errors.New("NULL")

type Transaction struct {
	storage *Storage
}

type Database struct {
	storage      *Storage
	transactions []Transaction
}

// NewDatabase returns new Database instance.
func NewDatabase() *Database {
	db := Database{}
	db.storage = NewStorage()
	return &db
}

// hasTransaction returns true if transaction is in progress.
func (db Database) hasTransaction() bool {
	return len(db.transactions) != 0
}

// getTransaction returns the most recent transation.
func (db Database) getTransaction() *Transaction {
	if !db.hasTransaction() {
		return nil
	}

	return &db.transactions[len(db.transactions)-1]
}

// RollbackTransaction undos all of the changes issued in the most recent transaction.
func (db *Database) RollbackTransaction() error {
	if !db.hasTransaction() {
		return ErrNoTransaction
	}

	db.transactions = db.transactions[:len(db.transactions)-1]
	return nil
}

// CommitTransactions closes all open transactions and applies the changes made in them.
func (db *Database) CommitTransactions() error {
	if !db.hasTransaction() {
		return ErrNoTransaction
	}

	t := db.getTransaction()
	db.storage = t.storage
	db.transactions = make([]Transaction, 0, 1)

	return nil
}

// getStorage returns using Storage.
func (db *Database) getStorage() *Storage {
	if db.hasTransaction() {
		return db.getTransaction().storage
	}

	return db.storage
}

// BeginTransaction start new transaction.
func (db *Database) BeginTransaction() *Transaction {
	t := Transaction{}
	t.storage = db.getStorage().Copy()
	db.transactions = append(db.transactions, t)
	return &t
}

// Set handles "SET name value" command.
func (db *Database) Set(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) != 3 {
		return ErrInvalidCommand
	}

	if parts[0] != "SET" {
		return ErrInvalidCommand
	}

	s := db.getStorage()
	s.Set(parts[1], parts[2])
	return nil
}

// Get handles "GET name" command.
func (db Database) Get(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return "", ErrInvalidCommand
	}

	if parts[0] != "GET" {
		return "", ErrInvalidCommand
	}

	s := db.getStorage()

	name := parts[1]
	if !s.HasVariable(name) {
		return "", ErrNoValue
	}

	return s.Get(name), nil
}

// Unset handles "UNSET name" command.
func (db *Database) Unset(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return ErrInvalidCommand
	}

	if parts[0] != "UNSET" {
		return ErrInvalidCommand
	}

	s := db.getStorage()
	s.Unset(parts[1])
	return nil
}

// NumEqualTo handles "NUMEQUALTO value" command.
func (db *Database) NumEqualTo(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return "", ErrInvalidCommand
	}

	if parts[0] != "NUMEQUALTO" {
		return "", ErrInvalidCommand
	}

	s := db.getStorage()
	return fmt.Sprintf("%d", s.NumEqualTo(parts[1])), nil
}

// Serve handles commands via stdin.
func (db *Database) Serve() {
	reader := bufio.NewReader(os.Stdin)

	for {
		l, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			return
		}

		cmd := strings.TrimSpace(l)

		switch {
		case cmd == "END":
			return

		case cmd == "BEGIN":
			db.BeginTransaction()

		case cmd == "ROLLBACK":
			if err := db.RollbackTransaction(); err != nil {
				fmt.Println(err)
			}

		case cmd == "COMMIT":
			if err := db.CommitTransactions(); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(cmd, "SET"):
			if err := db.Set(cmd); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(cmd, "GET"):
			output, err := db.Get(cmd)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(output)
			}

		case strings.HasPrefix(cmd, "UNSET"):
			if err := db.Unset(cmd); err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(cmd, "NUMEQUALTO"):
			output, err := db.NumEqualTo(cmd)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(output)
			}

		default:
			fmt.Println(ErrInvalidCommand)
		}
	}
}

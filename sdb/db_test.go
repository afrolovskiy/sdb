package sdb

import (
	"reflect"
	"testing"
)

func TestDatabaseSet(t *testing.T) {
	db := NewDatabase()

	err := db.Set("SET a")
	if err != ErrInvalidCommand {
		t.Errorf("Set must return error if command is invalid")
	}
	if db.storage.Get("a") != "" {
		t.Errorf("Invalid SET command must not change storage")
	}

	err = db.Set("SET  a  10")
	if err != nil {
		t.Errorf("Set must not return error if command is valid")
	}
	if db.storage.Get("a") != "10" {
		t.Errorf("Valid SET command must change storage")
	}
}

func TestDatabaseGet(t *testing.T) {
	db := NewDatabase()

	output, err := db.Get("GET a 10")
	if err != ErrInvalidCommand {
		t.Errorf("GET must return error if command is invalid")
	}

	output, err = db.Get("GET  a")
	if err != ErrNoValue {
		t.Errorf("GET must return error if variable does not exist")
	}

	db.storage.Set("a", "10")
	output, err = db.Get("GET  a")
	if err != nil || output != "10" {
		t.Errorf("GET must return variable value if variable exists")
	}
}

func TestDatabaseUnset(t *testing.T) {
	db := NewDatabase()

	err := db.Unset("UNSET a 10")
	if err != ErrInvalidCommand {
		t.Errorf("UNSET must return error if command is invalid")
	}

	err = db.Unset("UNSET a")
	if err != nil {
		t.Errorf("UNSET must not return anything if variable does not exist")
	}
	if db.storage.Get("a") != "" {
		t.Errorf("UNSET must not set value to the variable")
	}

	db.storage.Set("a", "10")
	err = db.Unset("UNSET a")
	if err != nil {
		t.Errorf("UNSET must not return anything if variable does not exist")
	}
	if db.storage.Get("a") != "" {
		t.Errorf("UNSET must remove the variable from storage")
	}
}

func TestDatabaseBeginTransaction(t *testing.T) {
	db := NewDatabase()

	if len(db.transactions) != 0 {
		t.Errorf("Before BEGIN command database must not have transactions")
	}

	db.BeginTransaction()
	if len(db.transactions) != 1 {
		t.Errorf("After BEGIN command new transaction must be created")
	}
	if reflect.DeepEqual(db.storage, db.transactions[0].storage) != true {
		t.Errorf("After BEGIN command database and transaction storages must be equal")
	}

	db.Set("SET a 10")
	db.BeginTransaction()
	if len(db.transactions) != 2 {
		t.Errorf("After BEGIN command new transaction must be created")
	}
	if reflect.DeepEqual(db.transactions[0].storage, db.transactions[1].storage) != true {
		t.Errorf("After BEGIN command previous transaction and current transaction storages must be equal")
	}
}

func TestDatabaseRollbackTransaction(t *testing.T) {
	db := NewDatabase()

	err := db.RollbackTransaction()
	if err != ErrNoTransaction {
		t.Errorf("RollbackTransaction must return error if transaction is not in progress")
	}

	t1 := db.BeginTransaction()
	db.BeginTransaction()

	err = db.RollbackTransaction()
	if err != nil {
		t.Errorf("RollbackTransaction must not return error if transaction is in progress")
	}
	if len(db.transactions) != 1 || reflect.DeepEqual(db.transactions[0], *t1) != true {
		t.Errorf("RollbackTransaction must rollback only last transaction")
	}
}

func TestDatabaseCommitTransactions(t *testing.T) {
	db := NewDatabase()

	err := db.CommitTransactions()
	if err != ErrNoTransaction {
		t.Errorf("CommitTransactions must return error if transaction is not in progress")
	}

	db.BeginTransaction()
	db.BeginTransaction()
	db.Set("SET a 10")
	err = db.CommitTransactions()
	if err != nil {
		t.Errorf("CommitTransactions must not return error if transaction is in progress")
	}
	if db.storage.Get("a") != "10" || db.hasTransaction() {
		t.Errorf("CommitTransactions must close all active transactions and apply the changes made in them")
	}
}

func TestDatabaseNumEqualTo(t *testing.T) {
	db := NewDatabase()

	output, err := db.NumEqualTo("NUMEQUALTO a 10")
	if err != ErrInvalidCommand {
		t.Errorf("NUMEQUALTO must return error if command is invalid")
	}

	output, err = db.NumEqualTo("NUMEQUALTO 10")
	if output != "0" || err != nil {
		t.Errorf("NUMEQUALTO must return the number of variables that are currently set to <valuex></valuex>")
	}
}

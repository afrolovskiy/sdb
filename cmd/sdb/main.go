package main

import "github.com/afrolovskiy/sdb/sdb"

func main() {
	db := sdb.NewDatabase()
	db.Serve()
}

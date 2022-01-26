package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	StorageDB *sql.DB
)

func Init() {
	dataSource := "admin:admin@tcp(localhost:3306)/storage"
	// Open inicia un pool de conexiones. Sólo abrir una vez
	var err error
	StorageDB, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	if err = StorageDB.Ping(); err != nil {
		panic(err)
	}
	log.Println("database Configured")
}

func Close() {
	StorageDB.Close()
}

func main() {
	Init()
	defer Close()
	//repo := products.NewRepo(StorageDB)
	//service := products.NewService(repo)
}

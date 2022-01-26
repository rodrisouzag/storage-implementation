package products

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rodrisouzag/storage-implementation/internal/models"
	"github.com/stretchr/testify/assert"
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
func TestStore(t *testing.T) {
	Init()
	defer Close()
	input := models.Product{
		ID:       0,
		Name:     "Arroz",
		Category: "Alimentos",
		Count:    20,
		Price:    50.0,
	}
	repo := NewRepo(StorageDB)

	stored, err := repo.Store(input)
	expected := models.Product{
		ID:       1,
		Name:     "Arroz",
		Category: "Alimentos",
		Count:    20,
		Price:    50.0,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, stored)
}

func TestGetByName(t *testing.T) {
	Init()
	defer Close()
	expected := models.Product{
		ID:       1,
		Name:     "Arroz",
		Category: "Alimentos",
		Count:    20,
		Price:    50.0,
	}

	repo := NewRepo(StorageDB)
	arroz := repo.GetByName("Arroz")

	assert.Equal(t, expected, arroz)
}

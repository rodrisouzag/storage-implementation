package products

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rodrisouzag/storage-implementation/internal/models"
	"github.com/rodrisouzag/storage-implementation/util"
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

func TestGetAll(t *testing.T) {
	Init()
	defer Close()
	p1 := models.Product{
		ID:       1,
		Name:     "Arroz",
		Category: "Alimentos",
		Count:    20,
		Price:    50.0,
	}
	var expected []models.Product
	expected = append(expected, p1)
	repo := NewRepo(StorageDB)
	result, err := repo.GetAll()

	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestUpdateWithContext(t *testing.T) {
	Init()
	defer Close()
	p1 := models.Product{
		ID:       1,
		Name:     "Pan",
		Category: "Alimentos",
		Count:    100,
		Price:    120.5,
	}

	repo := NewRepo(StorageDB)
	result, err := repo.UpdateWithContext(context.Background(), p1)

	assert.Equal(t, p1, result)
	assert.NoError(t, err)
}

// Tests TXDB

func Test_sqlRepository_Store(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)
	repository := NewRepo(db)
	ctx := context.TODO()
	productId := 1
	product := models.Product{
		ID:       productId,
		Name:     "Arroz",
		Category: "Alimentos",
		Count:    20,
		Price:    50.0,
	}
	_, err = repository.Store(product)
	assert.NoError(t, err)
	getResult, err := repository.GetOneWithContext(ctx, 45)
	assert.NoError(t, err)
	assert.Equal(t, models.Product{}, getResult)
	getResult, err = repository.GetOneWithContext(ctx, productId)
	assert.NoError(t, err)
	assert.NotNil(t, getResult)
	assert.Equal(t, product.ID, getResult.ID)
}

func Test_sqlRepository_Update(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)
	repository := NewRepo(db)
	ctx := context.TODO()
	productId := 1
	product := models.Product{
		ID:       productId,
		Name:     "Leche",
		Category: "Alimentos",
		Count:    100,
		Price:    52.50,
	}
	_, err = repository.UpdateWithContext(ctx, product)
	assert.NoError(t, err)
	getResult, err := repository.GetOneWithContext(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, product, getResult)
}

func Test_sqlRepository_Delete(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)
	repository := NewRepo(db)
	ctx := context.TODO()

	err = repository.Delete(1)
	assert.NoError(t, err)
	getResult, err := repository.GetOneWithContext(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, models.Product{}, getResult)
}

/*func Test_sqlRepository_Store_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare("INSERT INTO products")
	mock.ExpectExec("INSERT INTO products").WillReturnResult(sqlmock.NewResult(1, 1))
	columns := []string{"id", "name", "category", "count", "price"}
	rows := sqlmock.NewRows(columns)
	productId := 1
	rows.AddRow(productId, "", "", "", "")
	mock.ExpectQuery("SELECT * FROM products").WithArgs(productId).WillReturnRows(rows)
	repository := NewRepo(db)
	ctx := context.TODO()
	product := models.Product{
		ID:       productId,
		Name:     "Leche",
		Category: "Alimentos",
		Count:    100,
		Price:    52.50,
	}
	getResult, err := repository.GetOneWithContext(ctx, productId)
	assert.NoError(t, err)
	assert.Nil(t, getResult)
	_, err = repository.Store(product)
	assert.NoError(t, err)
	getResult, err = repository.GetOneWithContext(ctx, productId)
	assert.NoError(t, err)
	assert.NotNil(t, getResult)
	assert.Equal(t, product.ID, getResult.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}*/

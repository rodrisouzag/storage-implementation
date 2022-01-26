package products

import (
	"context"
	"database/sql"
	"log"

	"github.com/rodrisouzag/storage-implementation/internal/models"
)

const (
	GetAllProducts = "SELECT id, name, category, count, price FROM products"
	GetByName      = "SELECT id, name, category, count, price FROM products WHERE name = ?"
	StoreProduct   = "INSERT INTO products(name, category, count, price) VALUES( ?, ?, ?, ? )"
	UpdateProduct  = "UPDATE products SET name = ?, category = ?, count = ?, price = ? WHERE id = ?"
)

type Repository interface {
	GetByName(name string) models.Product
	Store(product models.Product) (models.Product, error)
	GetAll() ([]models.Product, error)
	UpdateWithContext(ctx context.Context, product models.Product) (models.Product, error)
}

type repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByName(name string) models.Product {
	var product models.Product
	db := r.db
	rows, err := db.Query(GetByName, name)
	if err != nil {
		log.Println(err)
		return product
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) Store(product models.Product) (models.Product, error) {
	db := r.db                            // se inicializa la base
	stmt, err := db.Prepare(StoreProduct) // se prepara el SQL
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	var result sql.Result
	result, err = stmt.Exec(product.Name, product.Category, product.Count, product.Price) // retorna un sql.Result y un error
	if err != nil {
		return models.Product{}, err
	}
	insertedId, _ := result.LastInsertId() // del sql.Resul devuelto en la ejecución obtenemos el Id insertado
	product.ID = int(insertedId)
	return product, nil
}

func (r *repository) GetAll() ([]models.Product, error) {
	var products []models.Product
	db := r.db
	rows, err := db.Query(GetAllProducts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// se recorren todas las filas
	for rows.Next() {
		// por cada fila se obtiene un objeto del tipo Product
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price); err != nil {
			log.Fatal(err)
			return nil, err
		}
		//se añade el objeto obtenido al slice products
		products = append(products, product)
	}
	return products, nil
}

func (r *repository) UpdateWithContext(ctx context.Context, product models.Product) (models.Product, error) {
	db := r.db                             // se inicializa la base
	stmt, err := db.Prepare(UpdateProduct) // se prepara la sentencia SQL a ejecutar
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.ExecContext(ctx, product.Name, product.Category, product.Count, product.Price, product.ID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

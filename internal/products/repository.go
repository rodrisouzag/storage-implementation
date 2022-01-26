package products

import (
	"database/sql"
	"log"

	"github.com/rodrisouzag/storage-implementation/internal/models"
)

type Repository interface {
	GetByName(name string) models.Product
	Store(product models.Product) (models.Product, error)
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
	rows, err := db.Query("select * from products where name = ?", name)
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
	db := r.db                                                                                         // se inicializa la base
	stmt, err := db.Prepare("INSERT INTO products(name, category, count, price) VALUES( ?, ?, ?, ? )") // se prepara el SQL
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

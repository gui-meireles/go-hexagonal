package db

import (
	"database/sql"
	"github.com/gui-meireles/go-hexagonal/application"
	_ "github.com/mattn/go-sqlite3"
)

type ProductDb struct {
	db *sql.DB
}

func NewProductDb(db *sql.DB) *ProductDb {
	return &ProductDb{db: db}
}

func (p *ProductDb) Get(id string) (application.ProductInterface, error) {
	var product application.Product
	stmt, err := p.db.Prepare("select id, name, price, status from products where id=?")
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price, &product.Status)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductDb) Save(product application.ProductInterface) (application.ProductInterface, error) {
	var rows int

	// Busca no banco pelo id do produto, caso exista algum produto ele atribuirá para a var 'rows'
	p.db.QueryRow("SELECT id FROM products WHERE id=?", product.GetID()).Scan(&rows)

	if rows == 0 {
		_, err := p.create(product)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := p.update(product)
		if err != nil {
			return nil, err
		}
	}
	return product, nil
}

// Métodos com a 1° letra minuscula são métodos privates
func (p *ProductDb) create(product application.ProductInterface) (application.ProductInterface, error) {
	stmt, err := p.db.Prepare(`INSERT INTO products(id,name,price,status) VALUES(?,?,?,?);`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		product.GetID(),
		product.GetName(),
		product.GetPrice(),
		product.GetStatus(),
	)
	if err != nil {
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Métodos com a 1° letra minuscula são métodos privates
func (p *ProductDb) update(product application.ProductInterface) (application.ProductInterface, error) {
	_, err := p.db.Exec("UPDATE products SET name = ?, price=?, status=? WHERE id=?",
		product.GetName(), product.GetPrice(), product.GetStatus(), product.GetID())

	if err != nil {
		return nil, err
	}
	return product, nil
}

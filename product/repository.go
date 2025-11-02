package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

const (
	createProduct = `INSERT INTO product 
	(id, name, description, price, stock, createdAt, updatedAt, isActive, userID)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	getProduct = `SELECT * from product WHERE id = $1`

	findAll = `SELECT id, name, description, price, stock, createdAt, updatedAt, isActive, userID 
	FROM product 
	ORDER BY createdAt $1`
)

type repository struct {
	Db *sql.DB
}

type Repository interface {
	CreateProduct(ctx context.Context, id uuid.UUID, name string, description string, price float64, stock int, userID uuid.UUID) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context, page int, limit int, sort string) ([]Product, error)
}

func NewRepository(Db *sql.DB) Repository {
	return &repository{
		Db: Db,
	}
}

func (r *repository) CreateProduct(ctx context.Context, id uuid.UUID, name string, description string, price float64, stock int, userID uuid.UUID) error {
	date := time.Now()

	_, err := r.Db.ExecContext(ctx, createProduct, id, name, description, price, stock, date, date, true, userID)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Product, error) {
	var product Product

	sqlErr := r.Db.QueryRowContext(ctx, getProduct, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.IsActive,
		&product.UserID,
	)

	if sqlErr == sql.ErrNoRows {
		return nil, ErrProductNotFound
	} else if sqlErr != nil {
		return nil, sqlErr
	}

	return &product, nil
}

func (r *repository) FindAll(ctx context.Context, page int, limit int, sort string) ([]Product, error) {

	if sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		offset := (page - 1) * limit
		query := findAll + ` LIMIT $2 OFFSET $3`
		
		rows, err := r.Db.QueryContext(ctx, query, sort, limit, offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		return scanProducts(rows)
	}

	rows, err := r.Db.QueryContext(ctx, findAll, sort)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows)
}

func scanProducts(rows *sql.Rows) ([]Product, error) {
	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.IsActive,
			&product.UserID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	createProduct = `INSERT INTO product 
	(id, name, description, price, stock, createdAt, updatedAt, isActive, userID)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
)

type repository struct {
	Db *sql.DB
}

type Repository interface {
	CreateProduct(ctx context.Context, id uuid.UUID, name string, description string, price float64, stock int, userID uuid.UUID) error
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

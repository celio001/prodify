package product_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/celio001/prodify/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	product_uuid := uuid.New()
	userid := uuid.New()

	repo_product := product.NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product (id, name, description, price, stock, createdAt, updatedAt, isActive, userID) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")).
		WithArgs(product_uuid, "product1", "novo produto cadastrado", 200.00, 5, time.Now(), time.Now(), true, userid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo_product.CreateProduct(context.Background(), product_uuid, "product1", "novo produto cadastrado", 200.00, 5, userid)

	assert.NoError(t, err)
}

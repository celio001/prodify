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
		WithArgs(product_uuid, "product1", "novo produto cadastrado", 200.00, 5, sqlmock.AnyArg(), sqlmock.AnyArg(), true, userid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo_product.CreateProduct(context.Background(), product_uuid, "product1", "novo produto cadastrado", 200.00, 5, userid)

	assert.NoError(t, err)
}

func TestGetProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	product_uuid := uuid.New()
	userid := uuid.New()

	repo_product := product.NewRepository(db)

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "createdAt", "updatedAt", "isActive", "userID"}).
		AddRow(product_uuid, "product1", "novo produto cadastrado", 200.00, 5, now, now, true, userid)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * from product WHERE id = $1")).
		WithArgs(product_uuid).
		WillReturnRows(rows)

	product, err := repo_product.FindByID(context.Background(), product_uuid.String())

	assert.NoError(t, err)
	assert.Equal(t, "product1", product.Name)
}

func TestFindAll_WithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo_product := product.NewRepository(db)

	product1_uuid := uuid.New()
	product2_uuid := uuid.New()
	userid := uuid.New()
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "createdAt", "updatedAt", "isActive", "userID"}).
		AddRow(product1_uuid, "product1", "description 1", 200.00, 5, now, now, true, userid).
		AddRow(product2_uuid, "product2", "description 2", 300.00, 10, now, now, true, userid)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, price, stock, createdAt, updatedAt, isActive, userID 
	FROM product 
	ORDER BY createdAt $1 LIMIT $2 OFFSET $3`)).
		WithArgs("asc", 2, 0).
		WillReturnRows(rows)

	products, err := repo_product.FindAll(context.Background(), 1, 2, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "product1", products[0].Name)
	assert.Equal(t, "product2", products[1].Name)
}

func TestFindAll_WithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo_product := product.NewRepository(db)

	product1_uuid := uuid.New()
	product2_uuid := uuid.New()
	product3_uuid := uuid.New()
	userid := uuid.New()
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "createdAt", "updatedAt", "isActive", "userID"}).
		AddRow(product1_uuid, "product1", "description 1", 200.00, 5, now, now, true, userid).
		AddRow(product2_uuid, "product2", "description 2", 300.00, 10, now, now, true, userid).
		AddRow(product3_uuid, "product3", "description 3", 400.00, 15, now, now, true, userid)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, price, stock, createdAt, updatedAt, isActive, userID 
	FROM product 
	ORDER BY createdAt $1`)).
		WithArgs("desc").
		WillReturnRows(rows)

	products, err := repo_product.FindAll(context.Background(), 0, 0, "desc")

	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "product1", products[0].Name)
	assert.Equal(t, "product2", products[1].Name)
	assert.Equal(t, "product3", products[2].Name)
}

func TestDeleteProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo_product := product.NewRepository(db)

	product1_uuid := uuid.New()

	userid := uuid.New()
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "createdAt", "updatedAt", "isActive", "userID"}).
		AddRow(product1_uuid, "product1", "description 1", 200.00, 5, now, now, true, userid)
	
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * from product WHERE id = $1")).
		WithArgs(product1_uuid).
		WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM product
	WHERE id = $1`)).
	WithArgs(product1_uuid).
	WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo_product.DeleteProduct(context.Background(), product1_uuid.String())

	assert.NoError(t, err)
}

func TestUpdateProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo_product := product.NewRepository(db)

	product_uuid := uuid.New()
	userid := uuid.New()
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	findRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "createdAt", "updatedAt", "isActive", "userID"}).
		AddRow(product_uuid, "product1", "description 1", 200.00, 5, now, now, true, userid)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * from product WHERE id = $1")).
		WithArgs(product_uuid).
		WillReturnRows(findRows)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE product 
	SET name = $1, 
		description = $2, 
		price = $3, 
		stock = $4, 
		updatedAt = $5, 
		isActive = $6 
	WHERE id = $7`)).
		WithArgs("Updated Product", "Updated description", 299.99, 10, sqlmock.AnyArg(), true, product_uuid).
		WillReturnResult(sqlmock.NewResult(0, 1))

	productToUpdate := &product.Product{
		ID:          product_uuid,
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       299.99,
		Stock:       10,
		IsActive:    true,
		UserID:      userid,
	}

	updatedProduct, err := repo_product.UpdateProduct(context.Background(), productToUpdate)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProduct.Name)
	assert.Equal(t, "Updated description", updatedProduct.Description)
	assert.Equal(t, 299.99, updatedProduct.Price)
	assert.Equal(t, 10, updatedProduct.Stock)
}
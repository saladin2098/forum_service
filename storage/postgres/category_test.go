package postgres

import (
	// "database/sql"
	"testing"

	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCategoryStorage(db)

	category := &pb.Category{
		CategoryId: "1",
		Name:       "Test Category",
	}

	mock.ExpectExec("insert into categories").
		WithArgs(category.CategoryId, category.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.CreateCategory(category)
	assert.NoError(t, err)
	assert.Equal(t, category, result)
}

func TestGetCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCategoryStorage(db)

	name := &pb.ByName{Name: "Test Category"}

	rows := sqlmock.NewRows([]string{"category_id", "name"}).
		AddRow("1", "Test Category")

	mock.ExpectQuery("select category_id, name from categories where name = \\$1 and deleted_at = 0").
		WithArgs(name.Name).
		WillReturnRows(rows)

	result, err := storage.GetCategory(name)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Category{CategoryId: "1", Name: "Test Category"}, result)
}

func TestListCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCategoryStorage(db)

	rows := sqlmock.NewRows([]string{"category_id", "name"}).
		AddRow("1", "Test Category").
		AddRow("2", "Another Category")

	mock.ExpectQuery("select category_id, name from categories where deleted_at = 0").
		WillReturnRows(rows)

	result, err := storage.ListCategories(&pb.Void{})
	assert.NoError(t, err)
	assert.Equal(t, &pb.Categories{
		Categories: []*pb.Category{
			{CategoryId: "1", Name: "Test Category"},
			{CategoryId: "2", Name: "Another Category"},
		},
	}, result)
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCategoryStorage(db)

	name := &pb.ByName{Name: "Test Category"}

	mock.ExpectExec("UPDATE categories SET deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) WHERE name = \\$1 and deleted_at = 0").
		WithArgs(name.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeleteCategory(name)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCategoryStorage(db)

	category := &pb.Category{
		CategoryId: "1",
		Name:       "Updated Category",
	}

	mock.ExpectExec("UPDATE categories SET name = \\$1 WHERE category_id = \\$2 and deleted_at = 0").
		WithArgs(category.Name, category.CategoryId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.UpdateCategory(category)
	assert.NoError(t, err)
	assert.Equal(t, category, result)
}

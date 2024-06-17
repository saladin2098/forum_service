package postgres

import (
	"testing"

	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	tag := &pb.Tag{
		TagId: "1",
		Name:  "Test Tag",
	}

	mock.ExpectExec("insert into tags").
		WithArgs(tag.TagId, tag.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.CreateTag(tag)
	assert.NoError(t, err)
	assert.Equal(t, tag, result)
}

func TestGetTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	name := &pb.ByName{Name: "Test Tag"}

	rows := sqlmock.NewRows([]string{"tag_id", "name"}).
		AddRow("1", "Test Tag")

	mock.ExpectQuery("select tag_id, name from tags where name = \\$1 and deleted_at = 0").
		WithArgs(name.Name).
		WillReturnRows(rows)

	result, err := storage.GetTag(name)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Tag{TagId: "1", Name: "Test Tag"}, result)
}

func TestDeleteTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	id := &pb.ById{Id: "1"}

	mock.ExpectExec("update tags set deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) where tag_id = \\$1 and deleted_at = 0").
		WithArgs(id.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeleteTag(id)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}

func TestUpdateTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	tag := &pb.Tag{
		TagId: "1",
		Name:  "Updated Tag",
	}

	mock.ExpectExec("UPDATE tags SET name = \\$1 WHERE tag_id = \\$2 and deleted_at = 0").
		WithArgs(tag.Name, tag.TagId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.UpdateTag(tag)
	assert.NoError(t, err)
	assert.Equal(t, tag, result)
}

func TestGetTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	rows := sqlmock.NewRows([]string{"tag_id", "name"}).
		AddRow("1", "Tag 1").
		AddRow("2", "Tag 2")

	mock.ExpectQuery("select tag_id, name from tags where deleted_at = 0").
		WillReturnRows(rows)

	result, err := storage.GetTags(&pb.Void{})
	assert.NoError(t, err)
	assert.Equal(t, &pb.TagList{
		Tags: []*pb.Tag{
			{TagId: "1", Name: "Tag 1"},
			{TagId: "2", Name: "Tag 2"},
		},
	}, result)
}

func TestGetPopularTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewTagStorage(db)

	rows := sqlmock.NewRows([]string{"tag_id", "name", "usage_count"}).
		AddRow("1", "Tag 1", 10).
		AddRow("2", "Tag 2", 8)

	mock.ExpectQuery("SELECT t.tag_id, t.name, COUNT\\(pt.tag_id\\) AS usage_count FROM tags t JOIN posts_tags pt ON t.tag_id = pt.tag_id WHERE t.deleted_at = 0 AND pt.deleted_at = 0 GROUP BY t.tag_id, t.name ORDER BY usage_count DESC LIMIT 10").
		WillReturnRows(rows)

	result, err := storage.GetPopularTags(&pb.Void{})
	assert.NoError(t, err)
	assert.Equal(t, &pb.TagList{
		Tags: []*pb.Tag{
			{TagId: "1", Name: "Tag 1"},
			{TagId: "2", Name: "Tag 2"},
		},
	}, result)
}

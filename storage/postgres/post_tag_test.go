package postgres

import (
	"testing"

	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostTagStorage(db)

	postTag := &pb.PostTag{
		PostTagId: "1",
		PostId:    "123",
		TagId:     "456",
	}

	mock.ExpectExec("insert into posts_tags").
		WithArgs(postTag.PostTagId, postTag.PostId, postTag.TagId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.CreatePostTag(postTag)
	assert.NoError(t, err)
	assert.Equal(t, postTag, result)
}

func TestDeletePostTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostTagStorage(db)

	id := &pb.ById{Id: "1"}

	mock.ExpectExec("update posts_tags set deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) where post_tag_id = \\$1 and deleted_at = 0").
		WithArgs(id.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeletePostTag(id)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}

func TestUpdatePostTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostTagStorage(db)

	postTag := &pb.PostTag{
		PostTagId: "1",
		PostId:    "123",
		TagId:     "456",
	}

	mock.ExpectExec("UPDATE posts_tags SET post_id = \\$1, tag_id = \\$2 WHERE deleted_at = 0 and post_tag_id = \\$3").
		WithArgs(postTag.PostId, postTag.TagId, postTag.PostTagId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.UpdatePostTag(postTag)
	assert.NoError(t, err)
	assert.Equal(t, postTag, result)
}

func TestGetPostTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostTagStorage(db)

	id := &pb.ById{Id: "1"}

	rows := sqlmock.NewRows([]string{"post_tag_id", "post_id", "tag_id"}).
		AddRow("1", "123", "456")

	mock.ExpectQuery("select post_tag_id, post_id, tag_id from posts_tags where post_tag_id = \\$1 and deleted_at = 0").
		WithArgs(id.Id).
		WillReturnRows(rows)

	result, err := storage.GetPostTag(id)
	assert.NoError(t, err)
	assert.Equal(t, &pb.PostTag{PostTagId: "1", PostId: "123", TagId: "456"}, result)
}

func TestGetPostTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostTagStorage(db)

	byPost := &pb.ByPost{PostId: "123"}

	rows := sqlmock.NewRows([]string{"post_tag_id", "post_id", "tag_id"}).
		AddRow("1", "123", "456").
		AddRow("2", "123", "789")

		mock.ExpectQuery("SELECT post_tag_id, post_id, tag_id FROM posts_tags WHERE deleted_at = 0 AND post_id = \\$1").
		WithArgs(byPost.PostId).
		WillReturnRows(rows)

	result, err := storage.GetPostTags(byPost)
	assert.NoError(t, err)
	assert.Equal(t, &pb.PostTags{
		PostTags: []*pb.PostTag{
			{PostTagId: "1", PostId: "123", TagId: "456"},
			{PostTagId: "2", PostId: "123", TagId: "789"},
		},
	}, result)
}


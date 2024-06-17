package postgres

import (
	// "database/sql"
	"testing"

	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCommentStorage(db)

	comment := &pb.Comment{
		CommentId: "1",
		PostId:    "123",
		UserId:    "456",
		Body:      "This is a comment",
	}

	mock.ExpectExec("insert into comments").
		WithArgs(comment.CommentId, comment.PostId, comment.UserId, comment.Body).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.CreateComment(comment)
	assert.NoError(t, err)
	assert.Equal(t, comment, result)
}

func TestGetComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCommentStorage(db)

	commentId := &pb.ById{Id: "1"}

	rows := sqlmock.NewRows([]string{"comment_id", "post_id", "user_id", "body"}).
		AddRow("1", "123", "456", "This is a comment")

	mock.ExpectQuery("select comment_id, post_id, user_id, body from comments where comment_id = \\$1 and deleted_at = 0").
		WithArgs(commentId.Id).
		WillReturnRows(rows)

	result, err := storage.GetComment(commentId)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Comment{CommentId: "1", PostId: "123", UserId: "456", Body: "This is a comment"}, result)
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCommentStorage(db)

	comment := &pb.Comment{
		CommentId: "1",
		PostId:    "123",
		UserId:    "456",
		Body:      "Updated comment",
	}

	mock.ExpectExec("UPDATE comments SET post_id = \\$1, user_id = \\$2, body = \\$3 WHERE comment_id = \\$4").
		WithArgs(comment.PostId, comment.UserId, comment.Body, comment.CommentId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.UpdateComment(comment)
	assert.NoError(t, err)
	assert.Equal(t, comment, result)
}

func TestDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCommentStorage(db)

	commentId := &pb.ById{Id: "1"}

	mock.ExpectExec("update comments set deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) where comment_id = \\$1 and deleted_at = 0").
		WithArgs(commentId.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeleteComment(commentId)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}

func TestGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewCommentStorage(db)

	filter := &pb.CommentFilter{
		PostId: "123",
		UserId: "456",
		Body:   "This is a comment",
	}

	rows := sqlmock.NewRows([]string{"comment_id", "post_id", "user_id", "body"}).
		AddRow("1", "123", "456", "This is a comment").
		AddRow("2", "123", "789", "Another comment")

	mock.ExpectQuery("SELECT comment_id, post_id, user_id, body FROM comments WHERE deleted_at = 0 AND post_id = \\$1 AND user_id = \\$2 AND body = \\$3").
		WithArgs(filter.PostId, filter.UserId, filter.Body).
		WillReturnRows(rows)

	result, err := storage.GetComments(filter)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Comments{
		Comments: []*pb.Comment{
			{CommentId: "1", PostId: "123", UserId: "456", Body: "This is a comment"},
			{CommentId: "2", PostId: "123", UserId: "789", Body: "Another comment"},
		},
	}, result)
}

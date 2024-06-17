package postgres

import (
	"testing"

	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	post := &pb.Post{
		PostId:     "1",
		UserId:     "123",
		Title:      "Test Post",
		Body:       "This is a test post",
		CategoryId: "456",
	}

	mock.ExpectExec("insert into posts").
		WithArgs(post.PostId, post.UserId, post.Title, post.Body, post.CategoryId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.CreatePost(post)
	assert.NoError(t, err)
	assert.Equal(t, post, result)
}

func TestGetPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	postId := &pb.ById{Id: "1"}

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "body", "category_id"}).
		AddRow("1", "123", "Test Post", "This is a test post", "456")

	mock.ExpectQuery("select post_id, user_id, title, body, category_id from posts where post_id = \\$1 and deleted_at = 0").
		WithArgs(postId.Id).
		WillReturnRows(rows)

	result, err := storage.GetPost(postId)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Post{PostId: "1", UserId: "123", Title: "Test Post", Body: "This is a test post", CategoryId: "456"}, result)
}

func TestGetPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	filter := &pb.PostFilter{
		CategoryId: "456",
		UserId:     "123",
		Title:      "Test Post",
		Body:       "This is a test post",
	}

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "body", "category_id"}).
		AddRow("1", "123", "Test Post", "This is a test post", "456").
		AddRow("2", "123", "Another Test Post", "This is another test post", "456")

	mock.ExpectQuery("select post_id, user_id, title, body, category_id from posts where deleted_at = 0 and category_id = \\$1 and user_id = \\$2 and title ILIKE \\$3 and body ILIKE \\$4").
		WithArgs(filter.CategoryId, filter.UserId, filter.Title, filter.Body).
		WillReturnRows(rows)

	result, err := storage.GetPosts(filter)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Posts{
		Posts: []*pb.Post{
			{PostId: "1", UserId: "123", Title: "Test Post", Body: "This is a test post", CategoryId: "456"},
			{PostId: "2", UserId: "123", Title: "Another Test Post", Body: "This is another test post", CategoryId: "456"},
		},
	}, result)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	post := &pb.Post{
		PostId:     "1",
		UserId:     "123",
		Title:      "Updated Post",
		Body:       "This is an updated post",
		CategoryId: "456",
	}

	mock.ExpectQuery("update posts set user_id = \\$1, title = \\$2, body = \\$3, category_id = \\$4 where post_id = \\$5 and deleted_at=0 returning post_id, user_id, title, body, category_id").
		WithArgs(post.UserId, post.Title, post.Body, post.CategoryId, post.PostId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "user_id", "title", "body", "category_id"}).
			AddRow(post.PostId, post.UserId, post.Title, post.Body, post.CategoryId))

	result, err := storage.UpdatePost(post)
	assert.NoError(t, err)
	assert.Equal(t, post, result)
}


func TestDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	postId := &pb.ById{Id: "1"}

	mock.ExpectExec("update posts set deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) where post_id = \\$1 and deleted_at = 0").
		WithArgs(postId.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeletePost(postId)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}


func TestGetPostsByTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewPostStorage(db)

	tagFilter := &pb.TagFilter{Tag: "tag1"}

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "body", "category_id"}).
		AddRow("1", "123", "Test Post", "This is a test post", "456").
		AddRow("2", "123", "Another Test Post", "This is another test post", "456")

	mock.ExpectQuery("SELECT p.post_id, p.user_id, p.title, p.body, p.category_id FROM posts p JOIN posts_tags pt ON p.post_id = pt.post_id JOIN tags t ON pt.tag_id = t.tag_id WHERE t.tag_id = \\$1 AND p.deleted_at = 0 AND pt.deleted_at = 0 AND t.deleted_at = 0").
		WithArgs(tagFilter.Tag).
		WillReturnRows(rows)

	result, err := storage.GetPostsByTag(tagFilter)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Posts{
		Posts: []*pb.Post{
			{PostId: "1", UserId: "123", Title: "Test Post", Body: "This is a test post", CategoryId: "456"},
			{PostId: "2", UserId: "123", Title: "Another Test Post", Body: "This is another test post", CategoryId: "456"},
		},
	}, result)
}

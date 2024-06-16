package postgres

import (
	"database/sql"
	"fmt"

	"github.com/saladin2098/forum_service/config"
	"github.com/saladin2098/forum_service/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	CategoryS storage.CategoryI
	PostS     storage.PostI
	CommentS  storage.CommentI
	TagS      storage.TagI
	PostTagS  storage.PostTagI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	categroryS := NewCategoryStorage(db)
	postS := NewPostStorage(db)
	commentS := NewCommentStorage(db)
	tagS := NewTagStorage(db)
	postTagS := NewPostTagStorage(db)
	return &Storage{
		db:        db,
        CategoryS: categroryS,
        PostS:     postS,
        CommentS:  commentS,
        TagS:      tagS,
        PostTagS:  postTagS,
	},nil
}
func (s *Storage) Category() storage.CategoryI {
	if s.CategoryS == nil {
		s.CategoryS = NewCategoryStorage(s.db)
	}
	return s.CategoryS
}
func (s *Storage) Post() storage.PostI {
	if s.PostS == nil {
        s.PostS = NewPostStorage(s.db)
    }
    return s.PostS
}
func (s *Storage) Comment() storage.CommentI {
	if s.CommentS == nil {
        s.CommentS = NewCommentStorage(s.db)
    }
    return s.CommentS
}
func (s *Storage) Tag() storage.TagI {
	if s.TagS == nil {
        s.TagS = NewTagStorage(s.db)
    }
    return s.TagS
}
func (s *Storage) PostTag() storage.PostTagI {
	if s.PostTagS == nil {
        s.PostTagS = NewPostTagStorage(s.db)
    }
    return s.PostTagS
}

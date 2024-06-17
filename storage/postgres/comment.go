package postgres

import (
	"database/sql"
	"strconv"
	"strings"

	pb "github.com/saladin2098/forum_service/genproto"
)

type CommentStorage struct {
	db *sql.DB
}

func NewCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{db: db}
}

func (s *CommentStorage) CreateComment(comment *pb.Comment) (*pb.Comment, error) {
	query := `insert into comments(
				comment_id,
				post_id,
                user_id,
                body) values($1,$2,$3,$4)`
	_, err := s.db.Exec(query,
		comment.CommentId,
		comment.PostId,
		comment.UserId,
		comment.Body)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
func (s *CommentStorage) GetComment(commentId *pb.ById) (*pb.Comment, error) {
	query := `select 
            comment_id,
            post_id,
            user_id,
            body
            from comments 
            where comment_id = $1 and deleted_at = 0`
	row := s.db.QueryRow(query, commentId.Id)
	comment := &pb.Comment{}
	err := row.Scan(
		&comment.CommentId,
		&comment.PostId,
		&comment.UserId,
		&comment.Body)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
func (s *CommentStorage) UpdateComment(comment *pb.Comment) (*pb.Comment, error) {
	///-------------WARNING--------------------------------////////////////////////////////------
	query := `UPDATE comments SET `
	var conditions []string
	var args []interface{}
	if comment.PostId != "" && comment.PostId != "string" {
		conditions = append(conditions, "post_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, comment.PostId)
	}
	if comment.UserId != "" && comment.UserId != "string" {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, comment.UserId)
	}
	if comment.Body != "" && comment.Body != "string" {
		conditions = append(conditions, "body = $"+strconv.Itoa(len(args)+1))
		args = append(args, comment.Body)
	}
	if len(conditions) > 0 {
		query += strings.Join(conditions, ", ")
	}
	query += ` WHERE comment_id = $` + strconv.Itoa(len(args)+1)
	args = append(args, comment.CommentId)
	_, err := s.db.Exec(query, args...)
	if err!= nil {
        return nil, err
    }
	return comment, nil
}
func (s *CommentStorage) DeleteComment(commentId *pb.ById) (*pb.Void, error) {
	query := `update comments set deleted_at = EXTRACT(EPOCH FROM NOW()) 
        	where comment_id = $1 and deleted_at = 0`
    _, err := s.db.Exec(query, commentId.Id)
    if err!= nil {
        return nil, err
    }
    return &pb.Void{}, nil
}
func (s *CommentStorage) GetComments(filter *pb.CommentFilter) (*pb.Comments, error) {
	query := `SELECT 
                comment_id,
                post_id,
                user_id,
                body
              FROM comments 
              WHERE deleted_at = 0`
	var args []interface{}
	
	if filter.PostId != "" && filter.PostId != "string" {
		query += ` AND post_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.PostId)
	}
	if filter.UserId != "" && filter.UserId != "string" {
		query += ` AND user_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.UserId)
	}
	if filter.Body != "" && filter.Body != "string" {
		query += ` AND body = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.Body)
	}
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	comments := &pb.Comments{}
	for rows.Next() {
		comment := &pb.Comment{}
		err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Body)
		if err != nil {
			return nil, err
		}
		comments.Comments = append(comments.Comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return comments, nil
}
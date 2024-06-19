package postgres

import (
	"database/sql"
	"strconv"
	"strings"

	pb "github.com/saladin2098/forum_service/genproto"
)

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{db: db}
}

func (s *PostStorage) CreatePost(post *pb.Post) (*pb.Post, error) {
	query := `insert into posts(
				post_id, 
				user_id,
				title, 
				body, 
				category_id) values($1,$2,$3,$4,$5)`
	_, err := s.db.Exec(query,
		post.PostId,
		post.UserId,
		post.Title,
		post.Body,
		post.CategoryId)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostStorage) GetPost(postId *pb.ById) (*pb.Post, error) {
	query := `select 
            post_id, 
            user_id,
            title, 
            body, 
            category_id 
            from posts 
            where post_id = $1 and deleted_at = 0`
	row := s.db.QueryRow(query, postId.Id)
	post := &pb.Post{}
	err := row.Scan(&post.PostId,
		&post.UserId,
		&post.Title,
		&post.Body,
		&post.CategoryId)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostStorage) GetPosts(filter *pb.PostFilter) (*pb.Posts, error) {
	query := `select 
            post_id, 
            user_id,
            title, 
            body, 
            category_id 
            from posts 
            where deleted_at = 0`
	var conditions []string
	var args []interface{}
	if filter.CategoryId != "" && filter.CategoryId != "string" {
		conditions = append(conditions, "category_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, filter.CategoryId)
	}
	if filter.UserId != "" && filter.UserId != "string" {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, filter.UserId)
	}
	if filter.Title != "" && filter.Title != "string" {
		conditions = append(conditions, "title ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, filter.Title)
	}
	if filter.Body != "" && filter.Body != "string" {
		conditions = append(conditions, "body ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, filter.Body)
	}
	if len(conditions) > 0 {
		query += " and " + strings.Join(conditions, " and ")
	}
	rows, err := s.db.Query(query, args...)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()
	posts := &pb.Posts{}
	for rows.Next() {
		post := &pb.Post{}
        err := rows.Scan(&post.PostId,
            &post.UserId,
            &post.Title,
            &post.Body,
            &post.CategoryId)
        if err!= nil {
            return nil, err
        }
        posts.Posts = append(posts.Posts, post)
	}
	return posts, nil
}
func (s *PostStorage) UpdatePost(post *pb.Post) (*pb.Post, error) {
	query := `update posts set `
	var conditions []string
	var args []interface{}
	if post.UserId != "" && post.UserId != "string" {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, post.UserId)
	}
	if post.Title != "" && post.Title != "string" {
		conditions = append(conditions, "title = $"+strconv.Itoa(len(args)+1))
		args = append(args, post.Title)
	}
	if post.Body != "" && post.Body != "string" {
		conditions = append(conditions, "body = $"+strconv.Itoa(len(args)+1))
		args = append(args, post.Body)
	}
	if post.CategoryId != "" && post.CategoryId != "string" {
		conditions = append(conditions, "category_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, post.CategoryId)
	}
	if len(conditions) > 0 {
		query += strings.Join(conditions, ", ")
	}
	query += " where post_id = $"+strconv.Itoa(len(args)+1)+" and deleted_at=0 returning post_id, user_id, title, body, category_id"
	args = append(args, post.PostId)
	row := s.db.QueryRow(query, args...)
	var res pb.Post
	err := row.Scan(&res.PostId, &res.UserId, &res.Title, &res.Body, &res.CategoryId)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *PostStorage) DeletePost(postId *pb.ById) (*pb.Void, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `update posts set deleted_at = EXTRACT(EPOCH FROM NOW())
		where post_id = $1 and deleted_at = 0`
	_, err = tx.Exec(query, postId.Id)
	if err != nil {
		return nil, err
	}

	query_com := `update comments set deleted_at = EXTRACT(EPOCH FROM NOW()) 
				where post_id = $1 and deleted_at = 0`
	_, err = tx.Exec(query_com, postId.Id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *PostStorage) GetPostsByTag(tag *pb.TagFilter) (*pb.Posts, error) {
	query := `SELECT 
				p.post_id,
				p.user_id,
				p.title, 
                p.body, 
                p.category_id
			FROM posts p
			JOIN posts_tags pt ON p.post_id = pt.post_id
			JOIN tags t ON pt.tag_id = t.tag_id
			WHERE t.tag_id = $1 AND p.deleted_at = 0 AND pt.deleted_at = 0 AND t.deleted_at = 0`
	
	rows, err := s.db.Query(query, tag.Tag)
	if err != nil {
        return nil, err
    }
	defer rows.Close()

	posts := &pb.Posts{}
	for rows.Next() {
		post := &pb.Post{}
        err := rows.Scan(&post.PostId,
            &post.UserId,
            &post.Title,
            &post.Body,
            &post.CategoryId)
        if err != nil {
            return nil, err
        }
        posts.Posts = append(posts.Posts, post)
	}
	if err = rows.Err(); err != nil {
        return nil, err
    }
	return posts, nil
}

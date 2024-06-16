package postgres

import (
	"database/sql"
	"strconv"
	"strings"

	pb "github.com/saladin2098/forum_service/genproto"
)

type PostTagStorage struct {
	db *sql.DB
}

func NewPostTagStorage(db *sql.DB) *PostTagStorage {
	return &PostTagStorage{db: db}
}

func (p *PostTagStorage) CreatePostTag(pt *pb.PostTag) (*pb.PostTag, error) {
	query := `insert into posts_tags(
			post_tag_id,
			post_id, 
			tag_id) values($1,$2,$3)
			`
	_, err := p.db.Exec(query, pt.PostTagId, pt.PostId, pt.TagId)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
func (p *PostTagStorage) DeletePostTag(pt *pb.ById) (*pb.Void, error) {
	query := `update posts_tags set deleted_at = EXTRACT(EPOCH FROM NOW()) 
			where post_tag_id = $1 and deleted_at = 0`
	_, err := p.db.Exec(query, pt.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (p *PostTagStorage) UpdatePostTag(pt *pb.PostTag) (*pb.PostTag, error) {
	query := `UPDATE posts_tags SET `
	var conditions []string
	var args []interface{}
	if pt.PostId != "" && pt.PostId != "string" {
		conditions = append(conditions, "post_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, pt.PostId)
	}
	if pt.TagId != "" && pt.TagId != "string" {
		conditions = append(conditions, "tag_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, pt.TagId)
	}
	if len(conditions) > 0 {
		query += strings.Join(conditions, ", ")
	}
	query += ` WHERE deleted_at = 0 and post_tag_id = $` + strconv.Itoa(len(args)+1)
	args = append(args, pt.PostTagId)
	_, err := p.db.Exec(query, args...)
	if err!= nil {
        return nil, err
    }
	return pt, nil
}
func (p *PostTagStorage) GetPostTag(id *pb.ById) (*pb.PostTag, error) {
	query := `select 
            post_tag_id, 
            post_id, 
            tag_id 
            from posts_tags 
            where post_tag_id = $1 and deleted_at = 0`
    row := p.db.QueryRow(query, id.Id)
    res := &pb.PostTag{}
    err := row.Scan(res.PostTagId, res.PostId, res.TagId)
    if err!= nil {
        return nil, err
    }
    return res, nil
}
func (p *PostTagStorage) GetPostTags(bp *pb.ByPost) (*pb.PostTags,error) {
	query := `select 
            post_tag_id, 
            post_id, 
            tag_id 
            from posts_tags 
            where post_id = $1 and deleted_at = 0`
    rows, err := p.db.Query(query, bp.PostId)
    if err!= nil {
        return nil, err
    }
    defer rows.Close()
    tags := &pb.PostTags{}
    for rows.Next() {
        tag := &pb.PostTag{}
        err := rows.Scan(&tag.PostTagId, &tag.PostId, &tag.TagId)
        if err!= nil {
            return nil, err
        }
        tags.PostTags = append(tags.PostTags, tag)
    }
    return tags, nil
}
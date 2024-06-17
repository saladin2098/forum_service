package postgres

import (
	"database/sql"

	pb "github.com/saladin2098/forum_service/genproto"
)

type TagStorage struct {
	db *sql.DB
}

func NewTagStorage(db *sql.DB) *TagStorage {
    return &TagStorage{db: db}
}

func (s *TagStorage) CreateTag(tag *pb.Tag) (*pb.Tag, error) {
	query := `insert into tags (tag_id,name) values($1,$2)`
	_, err := s.db.Exec(query, tag.TagId, tag.Name)
	if err!= nil {
        return nil, err
    }
	return tag, nil
}
func (s *TagStorage) GetTag(name *pb.ByName) (*pb.Tag, error) {
	query := `select 
            tag_id, 
            name 
            from tags 
            where name = $1 and deleted_at = 0`
    row := s.db.QueryRow(query, name.Name)
    tag := &pb.Tag{}
    err := row.Scan(&tag.TagId, &tag.Name)
    if err!= nil {
        return nil, err
    }
    return tag, nil
}
func (s *TagStorage) DeleteTag(id *pb.ById) (*pb.Void,error) {
	query := `update tags set deleted_at = EXTRACT(EPOCH FROM NOW())
        where tag_id = $1 and deleted_at = 0`
    _,err := s.db.Exec(query,id.Id)
    if err!= nil {
        return nil, err
    }
    return &pb.Void{}, nil
}
func (s *TagStorage) UpdateTag(tag *pb.Tag) (*pb.Tag, error) {
	query := `UPDATE tags 
            SET name = $1 WHERE tag_id = $2 and deleted_at = 0
            `
    _,err := s.db.Exec(query,tag.Name,tag.TagId)        
    if err!= nil {
        return nil, err
    }
    return tag, nil
}
func (s *TagStorage) GetTags(void *pb.Void) (*pb.TagList,error) {
	query := `select 
            tag_id, 
            name 
            from tags where deleted_at = 0`
    rows, err := s.db.Query(query)
    if err!= nil {
        return nil, err
    }
    defer rows.Close()
    tags := &pb.TagList{}
    for rows.Next() {
        tag := &pb.Tag{}
        err := rows.Scan(&tag.TagId, &tag.Name)
        if err!= nil {
            return nil, err
        }
        tags.Tags = append(tags.Tags, tag)
    }
    return tags, nil
}
func (s *TagStorage) GetPopularTags(*pb.Void) (*pb.TagList, error) {
        query := `SELECT 
                    t.tag_id, 
                    t.name,
                    COUNT(pt.tag_id) AS usage_count
                FROM tags t
                JOIN posts_tags pt ON t.tag_id = pt.tag_id
                WHERE t.deleted_at = 0 AND pt.deleted_at = 0
                GROUP BY t.tag_id, t.name
                ORDER BY usage_count DESC
                LIMIT 10`
                
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    tags := &pb.TagList{}
    var count int
    for rows.Next() {
        tag := &pb.Tag{}
        err := rows.Scan(&tag.TagId, &tag.Name,&count)
        if err != nil {
            return nil, err
    }
    tags.Tags = append(tags.Tags, tag)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return tags, nil
}
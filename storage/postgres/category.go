package postgres

import (
	"database/sql"

	pb "github.com/saladin2098/forum_service/genproto"
)

type CategoryStorage struct {
	db *sql.DB
}

func NewCategoryStorage(db *sql.DB) *CategoryStorage {
	return &CategoryStorage{db: db}
}

func (s *CategoryStorage) CreateCategory(cat *pb.Category) (*pb.Category,error) {
	query := `insert into categories(category_id, name) values($1,$2)`
	_, err := s.db.Exec(query, cat.CategoryId, cat.Name)
	if err!= nil {
        return nil, err
    }
	return cat, nil
}
func (s *CategoryStorage) GetCategory(name *pb.ByName) (*pb.Category, error) {
	query := `select 
			category_id, 
			name 
			from categories 
			where name = $1 and deleted_at = 0`
    row := s.db.QueryRow(query, name.Name)
    cat := &pb.Category{}
    err := row.Scan(&cat.CategoryId, &cat.Name)
    if err!= nil {
        return nil, err
    }
    return cat, nil
}
func (s *CategoryStorage) ListCategories(void *pb.Void) (*pb.Categories,error) {
	query := `select 
            category_id, 
            name 
            from categories where deleted_at = 0`
    rows, err := s.db.Query(query)
    if err!= nil {
        return nil, err
    }
    defer rows.Close()
    cats := &pb.Categories{}
    for rows.Next() {
        cat := &pb.Category{}
        err := rows.Scan(&cat.CategoryId, &cat.Name)
        if err!= nil {
            return nil, err
        }
        cats.Categories = append(cats.Categories, cat)
    }
    return cats, nil
}
func (s *CategoryStorage) DeleteCategory(name *pb.ByName) (*pb.Void, error) {
	query := `UPDATE categories 
			SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE name = $1 and deleted_at = 0
			`
	_,err := s.db.Exec(query,name.Name)		
	if err!= nil {
        return nil, err
    }
	return &pb.Void{}, nil
}
func (s *CategoryStorage) UpdateCategory(cat *pb.Category) (*pb.Category, error) {
	query := `UPDATE categories 
            SET name = $1 WHERE category_id = $2 and deleted_at = 0
            `
    _,err := s.db.Exec(query,cat.Name,cat.CategoryId)        
    if err!= nil {
        return nil, err
    }
    return cat, nil
}

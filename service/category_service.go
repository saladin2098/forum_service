package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/storage"
)

type CategoryService struct {
	stg storage.StorageI
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(stg storage.StorageI) *CategoryService {
	return &CategoryService{stg: stg}
}

func (s *CategoryService) CreateCategory(ctx context.Context, cat *pb.Category) (*pb.Category, error) {
	id := uuid.NewString()
	cat.CategoryId = id
	return s.stg.Category().CreateCategory(cat)
}

func (s *CategoryService) GetCategory(ctx context.Context, name *pb.ByName) (*pb.Category, error) {
    return s.stg.Category().GetCategory(name)
}

func (s *CategoryService) GetCategories(ctx context.Context, void *pb.Void) (*pb.Categories, error) {
    return s.stg.Category().ListCategories(void)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, name *pb.ByName) (*pb.Void, error) {
    return s.stg.Category().DeleteCategory(name)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, cat *pb.Category) (*pb.Category, error) {
    return s.stg.Category().UpdateCategory(cat)
}
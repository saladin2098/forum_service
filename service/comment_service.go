package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/storage"
)

type CommentService struct {
	stg storage.StorageI
	pb.UnimplementedCommentServiceServer
}

func NewCommentService(stg storage.StorageI) *CommentService {
	return &CommentService{stg: stg}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *pb.Comment) (*pb.Comment, error) {
	id := uuid.NewString()
	comment.CommentId = id
	return s.stg.Comment().CreateComment(comment)
}

func (s *CommentService) GetComment(ctx context.Context, id *pb.ById) (*pb.Comment, error) {
    return s.stg.Comment().GetComment(id)
}

func (s *CommentService) GetComments(ctx context.Context, filter *pb.CommentFilter) (*pb.Comments, error) {
    return s.stg.Comment().GetComments(filter)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *pb.Comment) (*pb.Comment, error) {
    return s.stg.Comment().UpdateComment(comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, id *pb.ById) (*pb.Void, error) {
    return s.stg.Comment().DeleteComment(id)
}

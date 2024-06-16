package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/storage"
)

type PostService struct {
	stg storage.StorageI
	pb.UnimplementedPostServiceServer
}

func NewPostService(stg storage.StorageI) *PostService {
	return &PostService{stg: stg}
}

func (s *PostService) CreatePost(ctx context.Context, post *pb.Post) (*pb.Post, error) {
	id := uuid.NewString()
	post.PostId = id
	return s.stg.Post().CreatePost(post)
}

func (s *PostService) GetPost(ctx context.Context, id *pb.ById) (*pb.Post, error) {
    return s.stg.Post().GetPost(id)
}

func (s *PostService) GetPostsByTag(ctx context.Context, tag *pb.TagFilter) (*pb.Posts, error) {
    return s.stg.Post().GetPostsByTag(tag)
}

func (s *PostService) GetPosts(ctx context.Context, filter *pb.PostFilter) (*pb.Posts, error) {
    return s.stg.Post().GetPosts(filter)
}

func (s *PostService) UpdatePost(ctx context.Context, post *pb.Post) (*pb.Post, error) {
    return s.stg.Post().UpdatePost(post)
}

func (s *PostService) DeletePost(ctx context.Context, id *pb.ById) (*pb.Void, error) {
    return s.stg.Post().DeletePost(id)
}

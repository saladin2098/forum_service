package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/storage"
)

type PostTagService struct {
	stg storage.StorageI
	pb.UnimplementedPostTagServiceServer
}

func NewPostTagService(stg storage.StorageI) *PostTagService {
	return &PostTagService{stg: stg}
}

func (ps *PostTagService) CreatePostTag(ctx context.Context,postTag *pb.PostTag) (*pb.PostTag,error) {
	id := uuid.NewString()
	postTag.PostTagId = id
	return ps.stg.PostTag().CreatePostTag(postTag)
}

func (ps *PostTagService) GetPostTag(ctx context.Context, id *pb.ById) (*pb.PostTag, error) {
    return ps.stg.PostTag().GetPostTag(id)
}

func (ps *PostTagService) DeletePostTagById(ctx context.Context, id *pb.ById) (*pb.Void, error) {
    return ps.stg.PostTag().DeletePostTag(id)
}

func (ps *PostTagService) GetPostTags(ctx context.Context, post *pb.ByPost) (*pb.PostTags, error) {
    return ps.stg.PostTag().GetPostTags(post)
}

func (ps *PostTagService) UpdatePostTag(ctx context.Context, postTag *pb.PostTag) (*pb.PostTag, error) {
    return ps.stg.PostTag().UpdatePostTag(postTag)
}

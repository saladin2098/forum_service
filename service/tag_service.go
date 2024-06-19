package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/storage"
)

type TagService struct {
	stg storage.StorageI
	pb.UnimplementedTagServiceServer
}

func NewTagService(stg storage.StorageI) *TagService {
	return &TagService{stg: stg}
}

func (s *TagService) CreateTag(ctx context.Context, tag *pb.Tag) (*pb.Tag, error) {
	if err := validateTagName(tag.Name); err != nil {
		return nil, err
	}
	id := uuid.NewString()
	tag.TagId = id
	return s.stg.Tag().CreateTag(tag)
}

func (s *TagService) GetTag(ctx context.Context, name *pb.ByName) (*pb.Tag, error) {
	return s.stg.Tag().GetTag(name)
}

func (s *TagService) GetTags(ctx context.Context, void *pb.Void) (*pb.TagList, error) {
	return s.stg.Tag().GetTags(void)
}

func (s *TagService) DeleteTag(ctx context.Context, id *pb.ById) (*pb.Void, error) {
	return s.stg.Tag().DeleteTag(id)
}

func (s *TagService) UpdateTag(ctx context.Context, tag *pb.Tag) (*pb.Tag, error) {
	if err := validateTagName(tag.Name); err != nil {
		return nil, err
	}
	return s.stg.Tag().UpdateTag(tag)
}

func (s *TagService) GetPopularTags(ctx context.Context, void *pb.Void) (*pb.TagList, error) {
	return s.stg.Tag().GetPopularTags(void)
}

func validateTagName(name string) error {
	const tagNamePattern = `^#[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(tagNamePattern, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid tag name: must start with # and contain only letters, numbers, and underscores")
	}
	return nil
}

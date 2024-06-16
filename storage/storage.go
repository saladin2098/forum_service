package storage

import pb "github.com/saladin2098/forum_service/genproto"

type StorageI interface {
	Category() CategoryI
	Post() PostI
	Comment() CommentI
	Tag() TagI
	PostTag() PostTagI
}
type CategoryI interface {
	CreateCategory(cat *pb.Category) (*pb.Category,error)
	GetCategory(name *pb.ByName) (*pb.Category, error)
	ListCategories(void *pb.Void) (*pb.Categories,error)
	DeleteCategory(name *pb.ByName) (*pb.Void, error)
	UpdateCategory(cat *pb.Category) (*pb.Category, error)
}
type PostI interface {
	CreatePost(post *pb.Post) (*pb.Post,error)
	GetPost(id *pb.ById) (*pb.Post, error)
	GetPosts(filter *pb.PostFilter) (*pb.Posts,error)
	UpdatePost(post *pb.Post) (*pb.Post, error)
	DeletePost(id *pb.ById) (*pb.Void, error)
	GetPostsByTag(tag *pb.TagFilter) (*pb.Posts, error)
}
type CommentI interface {
	CreateComment(comment *pb.Comment) (*pb.Comment,error)
	GetComment(id *pb.ById) (*pb.Comment, error)
	UpdateComment(comment *pb.Comment) (*pb.Comment, error)
	DeleteComment(id *pb.ById) (*pb.Void, error)
	GetComments(filter *pb.CommentFilter) (*pb.Comments,error)
}
type TagI interface {
	CreateTag(tag *pb.Tag) (*pb.Tag,error)
	GetTag(name *pb.ByName) (*pb.Tag, error)
	GetTags(void *pb.Void) (*pb.TagList,error)
	DeleteTag(id *pb.ById) (*pb.Void, error)
	UpdateTag(tag *pb.Tag) (*pb.Tag, error)
	GetPopularTags(*pb.Void) (*pb.TagList, error)
}
type PostTagI interface {
	CreatePostTag(postTag *pb.PostTag) (*pb.PostTag,error)
	GetPostTag(id *pb.ById) (*pb.PostTag, error)
	DeletePostTag(id *pb.ById) (*pb.Void, error)
	GetPostTags(post *pb.ByPost) (*pb.PostTags,error)
	UpdatePostTag(postTag *pb.PostTag) (*pb.PostTag, error)
}
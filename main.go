package main

import (
	"log"
	"net"

	"github.com/saladin2098/forum_service/config"
	pb "github.com/saladin2098/forum_service/genproto"
	"github.com/saladin2098/forum_service/service"
	"github.com/saladin2098/forum_service/storage/postgres"
	"google.golang.org/grpc"
)



func main() {
	db,err := postgres.ConnectDB()
	if err!= nil {
        log.Fatal(err)
    }
	cfg := config.Load()
	liss, err := net.Listen("tcp", cfg.HTTPPort)
	if err!= nil {
        log.Fatal(err)
    }
	s := grpc.NewServer()
	pb.RegisterCategoryServiceServer(s,service.NewCategoryService(db))
	pb.RegisterPostServiceServer(s,service.NewPostService(db))
	pb.RegisterCommentServiceServer(s,service.NewCommentService(db))
	pb.RegisterTagServiceServer(s,service.NewTagService(db))
	pb.RegisterPostTagServiceServer(s,service.NewPostTagService(db))
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err!= nil {
        log.Fatal(err)
    }
}
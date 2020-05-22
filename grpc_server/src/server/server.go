package server

import (
	"context"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/client"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/domain"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/proto/blogpb"
	"github.com/maei/shared_utils_go/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

type server struct{}

var (
	s = grpc.NewServer()
)

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	data := domain.BlogItem{
		AuthorID: req.GetBlog().GetAuthorId(),
		Content:  req.GetBlog().GetContent(),
		Title:    req.GetBlog().GetTitle(),
	}
	log.Println(data)
	insert, err := client.MongoCollection.InsertOne(context.Background(), data)
	if err != nil {
		logger.Error("gRPC-Server: Error while writing to Database", err)
		return nil, status.Errorf(codes.Internal, "gRPC-Server: Error while writing to Database")
	}
	oid, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "gRPC-Server: Cannot convert OID")
	}
	res := &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: req.GetBlog().GetAuthorId(),
			Content:  req.GetBlog().GetContent(),
			Title:    req.GetBlog().GetTitle(),
		},
	}

	return res, nil
}

func StartGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("error while listening gRPC Server", err)
	}

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		errServer := s.Serve(lis)
		if errServer != nil {
			logger.Error("error while serve gRPC Server", errServer)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	logger.Info("gRPC-Server: Stopping gRPC-Server")
	s.Stop()
	logger.Info("gRPC-Server: Closing gRPC-Server listener")
	lis.Close()

}

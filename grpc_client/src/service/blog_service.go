package service

import (
	"context"
	"fmt"
	"github.com/maei/golang_grpc_crud_mongo/grpc_client/src/client"
	"github.com/maei/golang_grpc_crud_mongo/grpc_client/src/proto/blogpb"
	"github.com/maei/shared_utils_go/logger"
)

var (
	BlogService blogServiceInterface = &blogService{}
)

type blogServiceInterface interface {
	CreatBlogItem()
}

type blogService struct{}

func (*blogService) CreatBlogItem() {
	cc, ccErr := client.GRPCClient.SetClient()
	if ccErr != nil {
		logger.Error("gRPC-Client: Error creating gRPC-Client", ccErr)
	}
	defer cc.Close()
	conn := blogpb.NewBlogServiceClient(cc)

	req := &blogpb.CreateBlogRequest{Blog: &blogpb.Blog{
		AuthorId: "maei",
		Content:  "Ich bin der beste, hoffentlich",
		Title:    "Mein Leben",
	}}

	res, resErr := conn.CreateBlog(context.Background(), req)
	if resErr != resErr {
		logger.Error("gRPC-Client: Error receiving Data from gRPC-Server", resErr)
	}
	fmt.Println(res)

}

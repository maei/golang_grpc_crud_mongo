package app

import "github.com/maei/golang_grpc_crud_mongo/grpc_server/src/client"

func StartApplication() {
	c := client.InitClient()
	db := client.InitDatabase(c)
	client.InitCollection(db)

}

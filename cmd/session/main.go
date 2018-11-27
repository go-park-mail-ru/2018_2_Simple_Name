package main

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/redis"
	"SimpleGame/internal/db/redis"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	listner, err := net.Listen("tcp", ":8081")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	redisService := db.NewRedisSessionService()

	if redisService == nil {
		return
	}

	server := grpc.NewServer()
	db.RegisterAuthCheckerServer(server, redisService)

	fmt.Println("Starting sess server at :8081")

	server.Serve(listner)
}

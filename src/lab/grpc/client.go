package main

import (
	"lab/grpc/inf"
	"log"
	"runtime"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	server   = "127.0.0.1"
	port     = ""
	parallel = 50
	times    = 100000
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	currTime := time.Now()

	//并行请求
	exe()
	upper()

	log.Printf("time taken: %.2f ", time.Now().Sub(currTime).Seconds())
}

func exe() {
	conn, _ := grpc.Dial(server+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := inf.NewDataClient(conn)

	getUser(client)
}

func getUser(client inf.DataClient) {
	var request inf.UserRq
	request.Id = int32(10000)

	response, err := client.GetUser(context.Background(), &request)
	if err != nil {
		log.Printf("client.GetUser failed. %v", err)
		return
	}

	log.Printf("response.Name: %s", response.Name)
}

func upper() {
	conn, err := grpc.Dial(server+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial failed. %v", err)
		return
	}
	defer conn.Close()
	client := inf.NewUpperClient(conn)

	resp, err := client.Upper(context.Background(), &inf.Req{
		Str: "hello",
	})
	if err != nil {
		log.Printf("client.Upper failed. %v", err)
		return
	}
	log.Printf("resp: %s", resp.Str)
}

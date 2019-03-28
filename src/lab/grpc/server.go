// grpc project main.go
package main

import (
	"flag"
	"fmt"
	"lab/grpc/inf"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	port int
)

type Data struct{}

type Upper struct{}

func args() {
	flag.IntVar(&port, "port", 5237, "service port")
	flag.Parse()
}

func main() {
	args()
	runtime.GOMAXPROCS(runtime.NumCPU())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}
	s := grpc.NewServer()
	inf.RegisterDataServer(s, &Data{})
	inf.RegisterUpperServer(s, &Upper{})

	s.Serve(lis)

	log.Println("grpc server in: %s", port)
}

// 定义方法
func (t *Data) GetUser(ctx context.Context, request *inf.UserRq) (response *inf.UserRp, err error) {
	response = &inf.UserRp{
		Name: strconv.Itoa(int(request.Id)) + ":test",
	}
	return response, err
}

func (u *Upper) Upper(ctx context.Context, req *inf.Req) (resp *inf.Resp, err error) {
	resp = &inf.Resp{
		Str: strings.ToUpper(req.Str),
	}
	return resp, err
}

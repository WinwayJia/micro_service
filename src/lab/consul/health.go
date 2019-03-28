package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"time"

	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthImpl 健康检查实现
type Health struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *Health) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *Health) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	port := 3000

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	entry := logrus.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(entry)

	server := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(server, &Health{})

	// 使用 consul 注册服务
	register := NewConsulRegister()
	register.Port = port
	if err := register.Register(); err != nil {
		panic(err)
	}

	address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		panic(err)
	}

	if err := server.Serve(address); err != nil {
		panic(err)
	}
}

// ConsulRegister consul service register
type ConsulRegister struct {
	Address                        string
	Service                        string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

// NewConsulRegister create a new consul register
func NewConsulRegister() *ConsulRegister {
	return &ConsulRegister{
		Address: "127.0.0.1:8500",
		Service: "unknown",
		Tag:     []string{},
		Port:    3000,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

// Register register service
func (r *ConsulRegister) Register() error {
	config := api.DefaultConfig()
	config.Address = r.Address
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	agent := client.Agent()

	IP := localIP()
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", r.Service, IP, r.Port), // 服务节点的名称
		Name:    fmt.Sprintf("grpc.health.v1.%v", r.Service),    // 服务名称
		Tags:    r.Tag,                                          // tag，可以为空
		Port:    r.Port,                                         // 服务端口
		Address: IP,                                             // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval: r.Interval.String(),                            // 健康检查间隔
			GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(), // 注销时间，相当于过期时间
		},
	}

	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}

	return nil
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

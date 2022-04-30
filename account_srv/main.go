package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-trainning-lessons/account_srv/biz"
	"mic-trainning-lessons/account_srv/proto/pb"
	"mic-trainning-lessons/internal"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	ip := flag.String("ip", "0.0.0.0", "输入ip")
	port := flag.Int("port", 9095, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv 启动异常" + err.Error())
		panic(err)
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	cheakAddr := fmt.Sprintf("%s:%d",
		internal.AppConf.AccountSrvConfig.Host,
		internal.AppConf.AccountSrvConfig.Port)
	check := api.AgentServiceCheck{
		GRPC:                           cheakAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.AccountSrvConfig.SrvName,
		Address: internal.AppConf.AccountSrvConfig.Host,
		Port:    internal.AppConf.AccountSrvConfig.Port,
		Tags:    internal.AppConf.AccountSrvConfig.Tags,
		Check:   &check,
	}
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		zap.S().Error("account_srv 启动异常" + err.Error())
		panic(err)
	}
	zap.S().Info("account_srv 启动成功")
}

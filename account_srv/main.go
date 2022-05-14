package main

import (
	"fmt"
	"github.com/go-basic/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-trainning-lessons/account_srv/biz"
	"mic-trainning-lessons/account_srv/proto/pb"
	"mic-trainning-lessons/internal"
	"mic-trainning-lessons/util"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	//ip := flag.String("ip", "0.0.0.0", "输入ip")
	//port := flag.Int("port", 9095, "输入端口")
	//flag.Parse()
	//addr := fmt.Sprintf("%s:%d", *ip, *port)

	port := util.GenRandomPort()
	addr := fmt.Sprintf("%s:%d", internal.AppConf.AccountSrvConfig.Host, port)
	// 将定义的对象注册grpc服务
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	// 启动服务监听
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv 启动异常" + err.Error())
		panic(err)
	}

	// 在consul 注册grpc 服务。
	// grpc 服务的健康检查  注册服务健康检查  启动的grpc  健康检查方法
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// consul的相关配置
	defaultConfig := api.DefaultConfig()
	// 配置地址
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	// 创建consul的客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	// 生成健康检查对象
	cheakAddr := fmt.Sprintf("%s:%d",
		internal.AppConf.AccountSrvConfig.Host,
		port)
	check := api.AgentServiceCheck{
		GRPC:                           cheakAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	randUUID := uuid.New()
	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.AccountSrvConfig.SrvName,
		Address: internal.AppConf.AccountSrvConfig.Host,
		ID:      randUUID,
		Port:    port,
		Tags:    internal.AppConf.AccountSrvConfig.Tags,
		Check:   &check,
	}
	// 注册grpc服务
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port))
	err = server.Serve(listen)
	if err != nil {
		zap.S().Error("account_srv 启动异常" + err.Error())
		panic(err)
	}
	zap.S().Info("account_srv 启动成功")
}

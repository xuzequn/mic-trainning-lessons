package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mic-trainning-lessons/account_web/handler"
	"mic-trainning-lessons/internal"
	"net"
)

func init() {
	fmt.Println(internal.AppConf.AccountWebConfig)
	err := internal.Reg(internal.AppConf.AccountWebConfig.Host,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.Port,
		internal.AppConf.AccountWebConfig.Tags)
	if err != nil {
		panic(err)
	}

}

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	fmt.Println(addrs)
	ip := flag.String("ip", "0.0.0.0", "输入Ip")
	port := flag.Int("port", 8081, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	r := gin.Default()
	accountGroup := r.Group("/v1/account")
	{
		accountGroup.GET("/list", handler.AccountListHandler)
		accountGroup.POST("/login", handler.LoginByPasswordHandler)
		accountGroup.GET("/captcha", handler.CaptchaHandler)
	}
	r.GET("/health", handler.HealthHandler)
	r.Run(addr)
}

package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" //
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mic-trainning-lessons/account_srv/proto/pb"
	"mic-trainning-lessons/account_web/req"
	"mic-trainning-lessons/account_web/res"
	"mic-trainning-lessons/custom_error"
	"mic-trainning-lessons/internal"
	"mic-trainning-lessons/jwt_op"
	"mic-trainning-lessons/log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func HandleError(err error) string {
	if err != nil {
		switch err.Error() {
		case custom_error.AccountExists:
			return custom_error.AccountExists
		case custom_error.AccountNotFound:
			return custom_error.AccountNotFound
		case custom_error.SaltError:
			return custom_error.SaltError
		default:
			return custom_error.InternalError

		}
	}
	return ""
}

//var accountSrvHost string
//var accountSrvPort int
var client pb.AccountServiceClient

//func initConsul() error {
//	// 创建注册中心客户端
//	defaultConfig := api.DefaultConfig()
//	consuleAddr := fmt.Sprintf("%s:%d",
//		internal.AppConf.ConsulConfig.Host,
//		internal.AppConf.ConsulConfig.Port)
//	defaultConfig.Address = consuleAddr
//	consulClient, err := api.NewClient(defaultConfig)
//	if err != nil {
//		zap.S().Error("AccountListHandler, 创建Consul的client失败:" + err.Error())
//		//c.JSON(http.StatusOK, gin.H{
//		//	"msg": "服务端内部错误",
//		//})
//		return err
//	}
//
//	// 调用注册中心上注册的服务的调用地址并调用。
//	accountSrvHost = ""
//	accountSrvPort = 0
//	// 通过服务名Service 过滤到服务，获取服务的提供者地址列表
//	serviceList, err := consulClient.Agent().ServicesWithFilter(`Service=="account_srv"`)
//	if err != nil {
//		zap.S().Error("AccountListHandler, 创建Consul获取服务列表失败:" + err.Error())
//		//c.JSON(http.StatusOK, gin.H{
//		//	"msg": "服务端内部错误",
//		//})
//		return err
//	}
//	// 多个配置负载均衡
//	for _, v := range serviceList {
//		accountSrvHost = v.Address
//		accountSrvPort = v.Port
//	}
//	return nil
//}

func initGrpcClient() error {

	// 从consul 配置负载均衡策略 获取服务调用地址。
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port)
	dialAddr := fmt.Sprintf("consul://%s/account_srv?wait=14", addr)
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		zap.S().Fatal(err)
	}

	// 通过获取的服务提供者地址信息进行grpc 调用
	//grpcAddr := fmt.Sprintf("%s:%d", accountSrvHost, accountSrvPort)
	//conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-Grpc拨号失败:%s", err.Error())
		log.Logger.Info(s)
		//e := HandleError(err)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": e,
		//})
		return err
	}
	client = pb.NewAccountServiceClient(conn)
	return nil
}

func init() {
	//err := initConsul()
	//if err != nil {
	//	panic(err)
	//}
	err := initGrpcClient()
	if err != nil {
		panic(err)
	}
}

func AccountListHandler(c *gin.Context) {
	// 解析参数
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "3")
	fmt.Println(pageNoStr, pageNoStr)
	// 调用业务逻辑
	// 将参数从string穿换成10进制 uint32
	pageNo, _ := strconv.ParseUint(pageNoStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	fmt.Println(pageNoStr, pageNoStr)
	r, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   uint32(pageNo),
		PageSize: uint32(pageSize),
	})
	// 错误处理
	if err != nil {
		s := fmt.Sprintf("AccountListHandler调用失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	//封装数据
	var resList []res.Account4Res
	// 通过grpc服务获取的protobuf格式返回信息，转换成accountres类型的结构体对象切片
	for _, item := range r.AccountList {
		resList = append(resList, pb2res(item))
	}
	// 拼装gin http response 报文
	c.JSON(http.StatusOK, gin.H{
		"msg":   "",
		"total": r.Total,
		"data":  resList,
	})
}

func pb2res(accountRes *pb.AccountRes) res.Account4Res {
	return res.Account4Res{
		Mobile:   accountRes.Mobile,
		NikeName: accountRes.Nikename,
		Gender:   accountRes.Gender,
	}

}

func LoginByPasswordHandler(c *gin.Context) {
	var loginByPassword req.LoginByPassword
	err := c.ShouldBind(&loginByPassword)
	if err != nil {
		log.Logger.Error("LoginByPassword出错:" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "解析参数错误",
		})
		return
	}
	// Down 校验手机号码格式
	// loginByPassword.Mobile不匹配正则表达式，就报错
	if !CheckMobile(loginByPassword.Mobile) {
		log.Logger.Error("LoginByPassword 手机号不合法")
		c.JSON(http.StatusOK, gin.H{
			"msg": "LoginByPassword 手机号不合法",
		})
		return
	}
	//conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithInsecure())
	//if err != nil {
	//	log.Logger.Error("LoginByPassword 拨号错误:" + err.Error())
	//	e := HandleError(err)
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": e,
	//	})
	//	return
	//}
	//client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountByMobile(context.Background(), &pb.MobileRequest{
		Mobile: loginByPassword.Mobile,
	})
	if err != nil {
		log.Logger.Error("LoginByPassword 错误:" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	cheRes, err := client.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       loginByPassword.Password,
		HashedPassword: r.Password,
		AccountId:      uint32(r.Id),
	})
	if err != nil {
		log.Logger.Error("LoginByPassword 错误:" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	checkResult := "登陆失败"
	if cheRes.Result {
		checkResult = "登陆成功"
		j := jwt_op.NewJWT()
		now := time.Now()
		claims := jwt_op.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				NotBefore: now.Unix(),
				ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
			},
			Id:          r.Id,
			NikeName:    r.Nikename,
			AuthorityId: int32(r.Role),
		}
		token, err := j.GenerateJWT(claims)
		if err != nil {
			log.Logger.Error("LoginByPassword 错误:" + err.Error())
			e := HandleError(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": e,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":    "",
			"result": checkResult,
			"token":  token,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "",
		"result": checkResult,
		"token":  "",
	})
}

func CheckMobile(mobile string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[3456789]{1}\\d{9}$"

	reg := regexp.MustCompile(regRuler)

	return reg.MatchString(mobile)

	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	//regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
}

func HealthHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

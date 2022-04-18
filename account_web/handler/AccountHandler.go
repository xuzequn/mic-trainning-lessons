package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"mic-trainning-lessons/account_srv/proto/pb"
	"mic-trainning-lessons/custom_error"
	"mic-trainning-lessons/log"
	"mic-trainning-lessons/res"
	"net/http"
	"strconv"
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

func AccountListHandler(c *gin.Context) {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "3")
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithInsecure())
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-Grpc拨号失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}

	pageNo, _ := strconv.ParseUint(pageNoStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   uint32(pageNo),
		PageSize: uint32(pageSize),
	})
	//r, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
	//	PageNo:   1,
	//	PageSize: 3,
	//})
	if err != nil {
		s := fmt.Sprintf("AccountListHandler调用失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	var resList []res.Account4Res
	for _, item := range r.AccountList {
		resList = append(resList, pb2res(item))
	}
	//log.Logger.Info("AccountListHandler调试通过")
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

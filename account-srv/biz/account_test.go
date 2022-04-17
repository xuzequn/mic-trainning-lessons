package biz

import (
	"context"
	"fmt"
	"mic-trainning-lessons/account-srv/internal"
	"mic-trainning-lessons/account-srv/proto/pb"
	"testing"
)

func init() {
	internal.InitDB()
}

func TestAccountServer_AddAccount(t *testing.T) {
	accountServer := AccountServer{}
	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("1300000000%d", i)
		res, err := accountServer.AddAccount(context.Background(), &pb.AddAccountRequest{
			Mobile:   s,
			Password: s,
			NikeName: s,
			Gender:   "male",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res.Id)
	}
}

func TestAccountServer_GetAccountList(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}

	res, err = accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   2,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}
}

func TestAccountServer_GetAccountByMobile(t *testing.T) {
	mobile := "13000000000"
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: mobile})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Id)

}

func TestAccountServer_GetAccountById(t *testing.T) {
	id := 3
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountById(context.Background(), &pb.IdRequest{Id: uint32(id)})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Mobile)

}

func TestAccountServer_UpdateAccount(t *testing.T) {
	var accountServer AccountServer
	res, err := accountServer.UpdateAccount(context.Background(), &pb.UpdateAccountRequest{
		Id:       1,
		Mobile:   "13000000100",
		Password: "13000000100",
		NikeName: "13000000100",
		Gender:   "female",
		Role:     2,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)

}

func TestAccountServer_CheckPassword(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       "13000000004",
		HashedPassword: "53fbfdbe4a43056ba2f2ad5cd522cfa72a14362ed7f07ca95fe8f5ec68fa717b",
		AccountId:      5,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)

}

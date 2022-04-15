package biz

import (
	"context"
	"gorm.io/gorm"
	"mic-trainning-lessons/account-srv/internal"
	"mic-trainning-lessons/account-srv/model"
	"mic-trainning-lessons/account-srv/proto/pb"
)

type AccountServer struct {
}

func Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo == 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (a *AccountServer) GetAccountList(ctx context.Context, req *pb.PagingRequest) (*pb.AccountListRes, error) {
	var accountList []model.Account
	//result := internal.DB.Find(&accountList)
	result := internal.DB.Scopes(Paginate(int(req.PageNo), int(req.PageSize))).Find(&accountList)
	if result.Error != nil {
		return nil, result.Error
	}
	accountListRes := &pb.AccountListRes{}

	accountListRes.Total = int32(result.RowsAffected)

	for _, account := range accountList {
		accountRes := Model2Pb(account)
		accountListRes.AccountList = append(accountListRes.AccountList, accountRes)
	}

	return accountListRes, nil
}

func Model2Pb(account model.Account) *pb.AccountRes {
	accountRes := &pb.AccountRes{
		Id:       int32(account.ID),
		Mobile:   account.Mobile,
		Password: account.Password,
		Nikename: account.NikeName,
		Gender:   account.Gender,
		Role:     uint32(account.Role),
	}
	return accountRes
}

func (a *AccountServer) GetAccountByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.AccountRes, error) {
	return &pb.AccountRes{}, nil
}

func (a *AccountServer) GetAccountById(ctx context.Context, req *pb.IdRequest) (*pb.AccountRes, error) {
	return &pb.AccountRes{}, nil
}

func (a *AccountServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.AccountRes, error) {
	return &pb.AccountRes{}, nil
}

func (a *AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountRes, error) {
	return &pb.UpdateAccountRes{Result: true}, nil
}

func (a *AccountServer) CheckPassword(ctx context.Context, req *pb.CheckPasswordRequest) (*pb.CheckPasswordRes, error) {
	return &pb.CheckPasswordRes{Result: true}, nil
}

func main() {

}

package biz

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"mic-trainning-lessons/account_srv/model"
	"mic-trainning-lessons/account_srv/proto/pb"
	"mic-trainning-lessons/custom_error"
	"mic-trainning-lessons/internal"
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
	fmt.Println(req.PageNo, req.PageSize)
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
	var account model.Account
	result := internal.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 0 { // 没找到
		return nil, errors.New(custom_error.AccountNotFound)
	}
	res := Model2Pb(account)
	return res, nil
}

func (a *AccountServer) GetAccountById(ctx context.Context, req *pb.IdRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := internal.DB.First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New(custom_error.AccountNotFound)
	}
	res := Model2Pb(account)
	return res, nil
}

func (a *AccountServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.AccountRes, error) {

	var account model.Account
	result := internal.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 1 { // 是否存在
		return nil, errors.New(custom_error.AccountExists)
	}
	// 创建账户
	account.Mobile = req.Mobile
	account.NikeName = req.NikeName
	account.Role = 1
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodePwd := password.Encode(req.Password, &options)
	account.Salt = salt
	account.Password = encodePwd
	r := internal.DB.Create(&account)
	if r.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}
	accountRes := Model2Pb(account)

	return accountRes, nil
}

func (a *AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountRes, error) {
	var account model.Account
	result := internal.DB.First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New(custom_error.AccountNotFound)
	}

	account.Mobile = req.Mobile
	account.NikeName = req.NikeName
	account.Gender = req.Gender

	r := internal.DB.Save(&account)
	if r.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}

	return &pb.UpdateAccountRes{Result: true}, nil
}

func (a *AccountServer) CheckPassword(ctx context.Context, req *pb.CheckPasswordRequest) (*pb.CheckPasswordRes, error) {
	var account model.Account
	result := internal.DB.First(&account, req.AccountId)
	if result.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}
	if account.Salt == "" {
		return nil, errors.New(custom_error.SaltError)
	}
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	r := password.Verify(req.Password, account.Salt, account.Password, &options)
	return &pb.CheckPasswordRes{Result: r}, nil
}

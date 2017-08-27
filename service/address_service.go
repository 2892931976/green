package service

import (
	"database/sql"
	"errors"
	"runtime"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
)

type AddressService struct {
	Service
}

// @desc 查找地址
func (this *AddressService) Search(user *model.User, search *SearchAddress) (*SearchData, error) {
	bean := new(model.Address)
	return GetSearchData(this.Db, user, bean, search, func(session *xorm.Session) {
		session.Cols("id", "name", "parent_id", "value")
		if search.Keyword != "" {
			session.Where("value like ?", search.Keyword+"%")
		}
		if search.ParentId != 0 {
			session.Where("parent_id=?", search.ParentId)
		}
	})
}

func NewAddressServ() *AddressService {
	config.Db.Sync2(new(model.Address))
	LoadSql()
	return new(AddressService)
}

/*****************************************************************************
 *                                 api above                                 *
 *****************************************************************************/

// @desc 导入地址数据
func LoadSql() ([]sql.Result, error) {
	Db := config.Db
	if ok, err := Db.IsTableEmpty(new(model.Address)); err == nil {
		if ok {
			res, err := Db.ImportFile(pkg_path() + "/address.sql")
			return res, err
		} else {
			return nil, errors.New("数据已经导入address表")
		}
	} else {
		return nil, err
	}
}

func pkg_path() string {
	var path string
	_, file, _, _ := runtime.Caller(0)
	if index := strings.LastIndex(file, "/"); index > 0 {
		path = file[:index]
	}
	return path
}

// 查找地址
type SearchAddress struct {
	SearchPage
	Keyword  string `json:"keyword,omitempty" gev:"地址 如:'北京|'"`
	ParentId int    `json:"parent_id,omitempty" gev:"父地址id"`
}

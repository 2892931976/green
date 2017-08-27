package service

import (
	"github.com/go-xorm/xorm"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
)

type IItem interface {
	GetId() int
}

type IOwnerItem interface {
	IItem
	GetOwnerId() int
}

type permissionFunc func() bool

func ItemSelect(Db *xorm.Session, user *model.User, bean IItem, canRead permissionFunc) (interface{}, error) {
	ok, err := Db.Where("id=?", bean.GetId()).Get(bean)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, config.ItemNotExistError
	}
	if canRead != nil && !canRead() {
		return nil, config.PermissionDenied
	}
	if data, ok := bean.(IDataDetail); ok {
		data.GetDetail(user)
		return data, err
	}
	return bean, nil
}

func itemCreate(Db *xorm.Session, user *model.User, bean IItem, canCreate permissionFunc) (interface{}, error) {
	if canCreate != nil && !canCreate() {
		return nil, config.PermissionDenied
	}
	_, err := Db.InsertOne(bean)
	if err != nil {
		return nil, err
	}
	if data, ok := bean.(IDataDetail); ok {
		data.GetDetail(user)
		return data, err
	}
	return bean, nil
}

func itemCanUpdate(Db *xorm.Session, user *model.User, bean IItem, canUpdate permissionFunc) (interface{}, error) {
	ok, err := Db.Where("id=?", bean.GetId()).Get(bean)
	if !ok {
		return nil, config.UserNotExistError
	}
	if err != nil {
		return nil, err
	}
	if canUpdate != nil && !canUpdate() {
		return nil, config.PermissionDenied
	}
	return bean, nil
}

func itemUpdate(Db *xorm.Session, user *model.User, bean IItem) (interface{}, error) {
	_, err := Db.Where("id=?", bean.GetId()).Update(bean)
	if err != nil {
		return nil, err
	}
	if data, ok := bean.(IDataDetail); ok {
		data.GetDetail(user)
		return data, nil
	}
	return bean, nil
}

func ItemSave(Db *xorm.Session, user *model.User, bean IItem, copyTo func() error, canCreate permissionFunc, canUpdate permissionFunc) (interface{}, error) {
	isCreate := bean.GetId() < 1
	if isCreate {
		err := copyTo()
		if err != nil {
			return nil, err
		}
		return itemCreate(Db, user, bean, canCreate)
	} else {
		_, err := itemCanUpdate(Db, user, bean, canUpdate)
		if err != nil {
			return nil, err
		}
		err = copyTo()
		if err != nil {
			return nil, err
		}
		return itemUpdate(Db, user, bean)
	}
}

func ItemDelete(Db *xorm.Session, user *model.User, bean IItem, canDelete permissionFunc) (int64, error) {
	ok, err := Db.Where("id=?", bean.GetId()).Get(bean)
	if !ok {
		return 0, config.ItemNotExistError
	}
	if err != nil {
		return 0, err
	}
	if canDelete != nil && !canDelete() {
		return 0, config.PermissionDenied
	}
	return Db.Where("id=?", bean.GetId()).Delete(bean)
}

func ItemDeleteByIds(Db *xorm.Session, user *model.User, bean IItem, ids []int) (int64, error) {
	if user == nil {
		return 0, config.NeedLoginError
	}
	if len(ids) < 1 {
		return 0, config.EmptyArrayError
	}
	if user.IsAdmin() {
		return Db.In("id", ids).Delete(bean)
	}
	return Db.In("id", ids).Where("owner_id=?", user.Id).Delete(bean)
}

func UserPermission(user *model.User, bean IOwnerItem) permissionFunc {
	return func() bool {
		if user != nil {
			return true
		}
		return false
	}
}

func OwnerPermission(user *model.User, bean IOwnerItem) permissionFunc {
	return func() bool {
		if user != nil && user.Id == bean.GetOwnerId() {
			return true
		}
		return false
	}
}

func AdminPermission(user *model.User, bean IOwnerItem) permissionFunc {
	return func() bool {
		if user != nil && user.IsAdmin() {
			return true
		}
		return false
	}
}

func CommonPermission(user *model.User, bean IOwnerItem) permissionFunc {
	return func() bool {
		if user == nil {
			return false
		}
		if user.IsAdmin() || user.Id == bean.GetOwnerId() {
			return true
		}
		return false
	}
}

func SelfPermission(user *model.User, bean *model.User) permissionFunc {
	return func() bool {
		if user == nil {
			return false
		}
		if user.Id == bean.Id || user.IsAdmin() || user.Id == bean.OwnerId {
			return true
		}
		return false
	}
}

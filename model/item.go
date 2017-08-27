package model

type Item struct {
	CreateAt Time `json:"create_at,omitempty" xorm:"created"`
	UpdateAt Time `json:"-" xorm:"updated"`
}

func (this *Item) GetDetail(user *User) {}

type IdItem struct {
	Id   int `json:"id,omitempty" xorm:"pk autoincr"`
	Item `xorm:"extends"`
}

func (this *IdItem) GetId() int {
	return this.Id
}

type OwnerItem struct {
	IdItem  `xorm:"extends"`
	OwnerId int `gev:"-" json:"-" xorm:""`
}

func (this *OwnerItem) GetOwnerId() int {
	return this.OwnerId
}

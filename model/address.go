package model

// Address Entity
type Address struct {
	Id       int    `json:"id,omitempty" xorm:"pk autoincr"`
	Center   string `json:"center,omitempty" xorm:"" gev:"中心经纬度"`
	Citycode string `json:"citycode,omitempty" xorm:"" gev:"城市区号"`
	Level    string `json:"level,omitempty" xorm:"" gev:"级别"`
	Name     string `json:"name,omitempty" xorm:"" gev:"城市名"`
	ParentId int    `json:"parent_id,omitempty" xorm:"" gev:"父地址"`
	Value    string `json:"value,omitempty" xorm:"" gev:"地址全称,如: 北京|北京市市辖区|东城区"`
}

func (this *Address) TableName() string {
	return "address"
}

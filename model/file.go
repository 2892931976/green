package model

type File struct {
	OwnerItem `xorm:"extends"`
	Ext       string `json:"ext,omitempty" xorm:"" gev:"文件后缀"`
	Place     string `json:"-" xorm:""`
	Filename  string `json:"filename,omitempty" xorm:"" gev:""`
	MD5       string `json:"-" xorm:"" gev:""`
	Url       string `json:"url" xorm:"" gev:"文件地址,需加上host,如http://www.tederen.com:8017/"`
}

func (this *File) TableName() string {
	return "file"
}

package service

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/inu1255/green/model"
)

// 查找时通过 GetSearch 筛选返回字段
type IDataSearch interface {
	GetSearch(user *model.User)
}

func GetSearchData(session *xorm.Session, user *model.User, bean interface{}, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	sessionFunc(session)
	total, _ := session.Count(bean)
	sessionFunc(session)
	session.Limit(condition.GetSize(), condition.GetBegin())
	SetOrderBy(session, condition.GetOrder())
	data := make([]interface{}, condition.GetSize())
	n := 0
	var err error
	if _, ok := bean.(IDataSearch); ok {
		err = session.Iterate(bean, func(i int, item interface{}) error {
			item.(IDataSearch).GetSearch(user)
			data[i] = item
			n++
			return nil
		})
	} else {
		err = session.Iterate(bean, func(i int, item interface{}) error {
			data[i] = item
			n++
			return nil
		})
	}
	return &SearchData{Content: data[:n], Total: total}, err
}

type SearchData struct {
	Content []interface{} `json:"content" xorm:"" gev:"数据数组"`
	Total   int64         `json:"total" xorm:"" gev:"数据总量"`
	Ext     interface{}   `json:"ext,omitempty" xorm:"" gev:"附加数据"`
}

type ISearch interface {
	GetBegin() int
	GetSize() int
	GetOrder() string
}

// 分页查询
type SearchPage struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	OrderBy string `json:"order_by,omitempty" gev:"排序规则:-id"`
}

func (this *SearchPage) GetSize() int {
	if this.Size < 1 {
		return 10
	}
	return this.Size
}

func (this *SearchPage) GetBegin() int {
	return this.Page * this.GetSize()
}

func (this *SearchPage) GetOrder() string {
	return this.OrderBy
}

func SetOrderBy(session *xorm.Session, order string) {
	if order != "" {
		orders := strings.Split(order, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	}
}

// 关键词查询
type SearchKeyword struct {
	SearchPage
	Keyword string `json:"keyword,omitempty" gev:"关键词"`
}

// %keyword%
func (this *SearchKeyword) GetWordLike() string {
	key := this.Keyword
	return strings.Join([]string{"%", key, "%"}, "")
}

// %k%e%y%w%o%r%d%
func (this *SearchKeyword) GetCharLike() string {
	ss := strings.Split(this.Keyword, "")
	key := strings.Join(ss, "%")
	return strings.Join([]string{"%", key, "%"}, "")
}

package config

import (
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	gomail "gopkg.in/gomail.v2"
	ini "gopkg.in/ini.v1"
)

var (
	cfg, _      = ini.LooseLoad("config/goweb.ini", "goweb.ini")
	Db          *xorm.Engine
	EmailDialer *gomail.Dialer
	Cfg         = &config{
		Mode:       "debug",
		UploadPath: "upload",
		DataBase: database{
			DriverName:     "sqlite3",
			DataSourceName: "./test.db",
		},
		Cookie: cookie{
			TokenExpire: 86400, // one day
		},
		Secure: secure{
			NoiseLength:      5,
			CaptchaLength:    6,
			VerifyCodeLengh:  4,
			VerifyCodeRetry:  10,
			CodeExpireSecond: 600,
		},
		Email: email{
			Username: "uniwise@aliyun.com",
			Password: "uniwise87",
			Title:    "云央科技",
			Host:     "smtp.aliyun.com",
			Port:     25,
		},
		User: user{
			Avatar: "http://xorm.io/img/favicon.png",
		},
	}
)

type config struct {
	Mode       string
	UploadPath string
	DataBase   database
	Cookie     cookie
	Secure     secure
	Email      email
	User       user
}

type database struct {
	DriverName     string
	DataSourceName string
}

type cookie struct {
	TokenExpire int
}

type secure struct {
	NoiseLength      int
	VerifyCodeLengh  int
	VerifyCodeRetry  int
	CaptchaLength    int
	CodeExpireSecond int
}

type email struct {
	Username string
	Password string
	Title    string
	Host     string
	Port     int
}

type user struct {
	Avatar string
}

func init() {
	cfg.MapTo(Cfg)
	InitDb()
	InitVerify()
}

func InitDb() {
	var err error
	Db, err = xorm.NewEngine(Cfg.DataBase.DriverName, Cfg.DataBase.DataSourceName)
	if err != nil {
		log.Println("数据库初始化失败", err)
	} else {
		if Cfg.Mode == "debug" {
			Db.ShowSQL(true)
		}
	}
}

func InitVerify() {
	username := Cfg.Email.Username
	password := Cfg.Email.Password
	host := Cfg.Email.Host
	port := Cfg.Email.Port
	EmailDialer = gomail.NewDialer(host, port, username, password)
}

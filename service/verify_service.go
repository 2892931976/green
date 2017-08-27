package service

import (
	"regexp"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
	"github.com/labstack/gommon/random"
	gomail "gopkg.in/gomail.v2"
)

var (
	email_regex = regexp.MustCompile(`[.\d\w]+@[\d\w]+\.[\d\w]+`)
)

type VerifyService struct {
	Service
}

// @desc 发送邮箱验证码
func (this *VerifyService) SendEmailCode(email string) (bool, error) {
	return this.sendCode(email, func(bean *model.Verify) error {
		if !email_regex.Match([]byte(email)) {
			return config.EmailFormatError
		}
		username := config.Cfg.Email.Username
		title := config.Cfg.Email.Title
		code := RandCode()
		m := gomail.NewMessage()
		m.SetAddressHeader("From", username, title)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "<"+title+">验证码")
		m.SetBody("text/html", "<html>您的验证码是[ "+code+" ]，请勿告诉他人</html>")
		bean.SendCode(email, code)
		return config.EmailDialer.DialAndSend(m)
	})
}

// @desc 发送短信验证码
func (this *VerifyService) SendTelphoneCode(ctx *gin.Context, telphone, imagecode string) (bool, error) {
	if err := this.JudgeImageCode(ctx, imagecode); err != nil {
		return false, err
	}
	return this.sendCode(telphone, func(bean *model.Verify) error {
		// TODO
		code := RandCode()
		bean.SendCode(telphone, code)
		return nil
	})
}

// @desc 图片验证码
func (this *VerifyService) Image(ctx *gin.Context) (interface{}, error) {
	length := config.Cfg.Secure.CaptchaLength
	key := captcha.NewLen(length)
	maxAge := config.Cfg.Secure.CodeExpireSecond
	ctx.SetCookie("imagecode", key, maxAge, "/", "", false, false)
	err := captcha.WriteImage(ctx.Writer, key, 40*length, 80)
	return nil, err
}

func NewVerifyServ() *VerifyService {
	config.Db.Sync2(new(model.Verify))
	return new(VerifyService)
}

/*****************************************************************************
 *                                 api above                                 *
 *****************************************************************************/

func (this *VerifyService) sendCode(title string, send func(*model.Verify) error) (bool, error) {
	Db := this.Db
	bean := new(model.Verify)
	ok, _ := Db.Where("title=?", title).Get(bean)
	if bean.CanSend() {
		return false, config.CodeSendTooOffenError
	}
	err := send(bean)
	if err != nil {
		return false, err
	}
	// Log.Printf("class:%T - %v", bean, bean)
	if ok {
		_, err = Db.Where("title=?", title).Update(bean)
	} else {
		_, err = Db.InsertOne(bean)
	}
	return true, err
}

func (this *VerifyService) JudgeEmailCode(email, code string) error {
	Db := this.Db
	bean := new(model.Verify)
	ok, _ := Db.Where("title=?", email).Get(bean)
	if !ok {
		return config.CodeNotSendError
	}
	if bean.IsExpired() {
		return config.CodeExpiredError
	}
	if bean.Rest < 1 {
		return config.CodeNoRestError
	}
	if bean.Code != code {
		bean.Rest--
		Db.ID(bean.Id).Cols("rest").Update(bean)
		return config.CodeWrongError
	}
	// 验证成功时 使 验证码失效
	bean.Rest = -1
	Db.ID(bean.Id).Cols("rest").Update(bean)
	return nil
}

func (this *VerifyService) JudgeTelphoneCode(telphone, code string) error {
	return this.JudgeEmailCode(telphone, code)
}

func (this *VerifyService) JudgeImageCode(ctx *gin.Context, code string) error {
	id, err := ctx.Cookie("imagecode")
	if err != nil {
		return config.ImageCodeExpiredError
	}
	if !captcha.VerifyString(id, code) {
		return config.ImageCodeWrongError
	}
	return nil
}

func RandCode() string {
	length := config.Cfg.Secure.VerifyCodeLengh
	r := random.New()
	r.SetCharset(random.Numeric)
	return r.String(uint8(length))
}

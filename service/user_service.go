package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
	"github.com/labstack/gommon/random"
)

type UserService struct {
	Service
}

// @desc 注销
func (this *UserService) Logout(ctx *gin.Context, user *model.User) (interface{}, error) {
	if user == nil {
		return "没有登录", nil
	}
	device := ctx.Request.Header.Get("X-DEVICE")
	if device == "" {
		return this.Db.Exec("update access_token set expired_at='1993-03-07' where user_id=?", user.Id)
	}
	return this.Db.Exec("update access_token set expired_at='1993-03-07' where user_id=? and device=?", user.Id, device)
}

// @desc 登录
func (this *UserService) Login(ctx *gin.Context, key, password string) (*model.UserAccess, error) {
	if key == "" {
		return nil, config.ParamLackError
	}
	user := new(model.User)
	var ok bool
	var err error
	if strings.ContainsRune(key, '@') {
		// 邮箱登录
		ok, err = this.Db.Where("email=?", key).Get(user)
	} else {
		// 手机登录
		ok, err = this.Db.Where("telphone=?", key).Get(user)
	}
	// 通过手机号查用户
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, config.UserNotExistError
	}
	// 匹配密码
	if user.CheckPassword(password) {
		return this.GetUserAccess(ctx, user, "login")
	}
	return nil, config.PwdWrongError
}

// @desc email 修改密码
func (this *UserService) ChangePasswordEmail(ctx *gin.Context, rbody *EmailRegisterBody) (*model.UserAccess, error) {
	if rbody.Email == "" {
		return nil, config.EmailEmptyError
	}
	if !this.JudgePwdLength(len(rbody.Password)) {
		return nil, config.PwdLengthError
	}
	user, ok := this.GetByEmail(rbody.Email)
	if !ok {
		return nil, config.UserNotExistError
	}
	if err := this.Verify().JudgeEmailCode(rbody.Email, rbody.Code); err != nil {
		return nil, config.CodeWrongError
	}
	user.Password = model.EncodePwd(user.Noise, rbody.Password)
	_, err := this.Db.ID(user.Id).Cols("password").Update(user)
	if err != nil {
		return nil, err
	}
	return this.GetUserAccess(ctx, user, "chpwd")
}

// @desc telphone 修改密码
func (this *UserService) ChangePasswordTelphone(ctx *gin.Context, rbody *TelphoneRegisterBody) (*model.UserAccess, error) {
	if rbody.Telphone == "" {
		return nil, config.TelphoneEmptyError
	}
	if !this.JudgePwdLength(len(rbody.Password)) {
		return nil, config.PwdLengthError
	}
	user, ok := this.GetByTelphone(rbody.Telphone)
	if !ok {
		return nil, config.UserNotExistError
	}
	if err := this.Verify().JudgeTelphoneCode(rbody.Telphone, rbody.Code); err != nil {
		return nil, config.CodeWrongError
	}
	user.Password = model.EncodePwd(user.Noise, rbody.Password)
	_, err := this.Db.ID(user.Id).Cols("password").Update(user)
	if err != nil {
		return nil, err
	}
	return this.GetUserAccess(ctx, user, "chpwd")
}

// @desc 旧密码 修改密码
func (this *UserService) ChangePasswordOld(ctx *gin.Context, rbody *OldpwdBody) (*model.UserAccess, error) {
	if !this.JudgePwdLength(len(rbody.Password)) {
		return nil, config.PwdLengthError
	}
	var user *model.User
	var ok bool
	if rbody.Telphone != "" {
		if user, ok = this.GetByTelphone(rbody.Telphone); !ok {
			return nil, config.UserNotExistError
		}
	} else if rbody.Email != "" {
		if user, ok = this.GetByEmail(rbody.Email); !ok {
			return nil, config.UserNotExistError
		}
	} else {
		return nil, config.ParamLackError
	}
	if !user.CheckPassword(rbody.Oldpwd) {
		return nil, config.PwdWrongError
	}
	user.Password = model.EncodePwd(user.Noise, rbody.Password)
	_, err := this.Db.ID(user.Id).Cols("password").Update(user)
	if err != nil {
		return nil, err
	}
	return this.GetUserAccess(ctx, user, "chpwd")
}

// @desc 邮箱注册
func (this *UserService) RegisterEmail(ctx *gin.Context, rbody *EmailRegisterBody) (*model.UserAccess, error) {
	if rbody.Email == "" {
		return nil, config.EmailEmptyError
	}
	// 密码长度
	if !this.JudgePwdLength(len(rbody.Password)) {
		return nil, config.PwdLengthError
	}
	user := new(model.User)
	// 邮箱注册
	ok, _ := this.Db.Where("email=?", rbody.Email).Get(user)
	if ok {
		return nil, config.UserExistError
	}
	if err := this.Verify().JudgeEmailCode(rbody.Email, rbody.Code); err != nil {
		return nil, err
	}
	user.Email = rbody.Email
	user.Noise = RandNoise()
	user.Password = model.EncodePwd(user.Noise, rbody.Password)
	_, err := this.Db.InsertOne(user)
	if err != nil {
		return nil, err
	}
	return this.GetUserAccess(ctx, user, "regist")
}

// @desc 手机注册
func (this *UserService) RegisterTelphone(ctx *gin.Context, rbody *TelphoneRegisterBody) (*model.UserAccess, error) {
	if rbody.Telphone == "" {
		return nil, config.TelphoneEmptyError
	}
	// 密码长度
	if !this.JudgePwdLength(len(rbody.Password)) {
		return nil, config.PwdLengthError
	}
	user := new(model.User)
	ok, _ := this.Db.Where("telphone=?", rbody.Telphone).Get(user)
	if ok {
		return nil, config.UserExistError
	}
	if err := this.Verify().JudgeTelphoneCode(rbody.Telphone, rbody.Code); err != nil {
		return nil, err
	}
	user.Telphone = rbody.Telphone
	user.Noise = RandNoise()
	user.Password = model.EncodePwd(user.Noise, rbody.Password)
	_, err := this.Db.InsertOne(user)
	if err != nil {
		return nil, err
	}
	return this.GetUserAccess(ctx, user, "regist")
}

// @desc 绑定电话号码
func (this *UserService) BindTelphone(user *model.User, telphone string, code string) (interface{}, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	if user.Telphone != "" {
		return nil, config.NeedUnbindedError
	}
	if ok, _ := this.Db.Where("telphone=?", telphone).Exist(new(model.User)); ok {
		return nil, config.TelphoneBindedError
	}
	if err := this.Verify().JudgeTelphoneCode(telphone, code); err != nil {
		return nil, err
	}
	user.Telphone = telphone
	_, err := this.Db.Where("id=?", user.Id).Cols("telphone").Update(user)
	return user, err
}

// @desc 绑定邮箱
func (this *UserService) BindEmail(user *model.User, email string, code string) (interface{}, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	if user.Email != "" {
		return nil, config.NeedUnbindedError
	}
	if ok, _ := this.Db.Where("Email=?", email).Exist(user); ok {
		return nil, config.EmailBindedError
	}
	if err := this.Verify().JudgeEmailCode(email, code); err != nil {
		return nil, err
	}
	user.Email = email
	_, err := this.Db.Where("id=?", user.Id).Cols("email").Update(user)
	return user, err
}

// @desc 解绑手机号
func (this *UserService) UnbindTelphone(user *model.User) (interface{}, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	if user.Email == "" {
		return nil, config.NeedTelphoneOrEmail
	}
	user.Telphone = ""
	_, err := this.Db.Where("id=?", user.Id).Cols("telphone").Update(user)
	return user, err
}

// @desc 解绑邮箱
func (this *UserService) UnbindEmail(user *model.User) (interface{}, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	if user.Telphone == "" {
		return nil, config.NeedTelphoneOrEmail
	}
	user.Email = ""
	_, err := this.Db.Where("id=?", user.Id).Cols("email").Update(user)
	return user, err
}

// @desc 我的信息
func (this *UserService) MineInfo(user *model.User) (*model.User, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	user.GetDetail(user)
	return user, nil
}

// @desc 管理员删除用户
func (this *UserService) Delete(user *model.User, ids []int) (int64, error) {
	bean := new(model.User)
	return ItemDeleteByIds(this.Db, user, bean, ids)
}

// @desc 管理员创建用户
func (this *UserService) Create(user *model.User, title, password string) (interface{}, error) {
	bean := new(model.User)
	return ItemSave(this.Db, user, bean, func() error {
		bean.Telphone = title
		bean.Noise = RandNoise()
		bean.Password = model.EncodePwd(bean.Noise, password)
		return nil
	}, AdminPermission(user, bean), AdminPermission(user, bean))
}

// @desc 修改自己头像
func (this *UserService) ChangeAvatar(user *model.User, file *multipart.FileHeader) (string, error) {
	if user == nil {
		return "", config.NeedLoginError
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	bean, err := FileUpload(this.Db, file.Filename, src, user)
	if err != nil {
		return "", err
	}
	user.Avatar = bean.Place
	_, err = this.Db.ID(user.Id).Update(user)
	return user.Avatar, err
}

// @desc 修改个人信息
func (this *UserService) ChangeInfo(user *model.User, body *UserBody) (interface{}, error) {
	if user == nil {
		return nil, config.NeedLoginError
	}
	permission := SelfPermission(user, user)
	return ItemSave(this.Db, user, user, func() error {
		user.Nickname = body.Nickname
		user.Avatar = body.Avatar
		return nil
	}, permission, permission)
}

func NewUserServ() *UserService {
	config.Db.Sync2(new(model.User))
	config.Db.Sync2(new(model.AccessToken))
	return new(UserService)
}

/*****************************************************************************
 *                                 api above                                 *
 *****************************************************************************/

// 限制密码长度
func (this *UserService) JudgePwdLength(length int) bool {
	if length < 6 || length > 32 {
		return false
	}
	return true
}

// 通过telphone获取取用户
func (this *UserService) GetByTelphone(telphone string) (*model.User, bool) {
	user := new(model.User)
	ok, _ := this.Db.Where("telphone=?", telphone).Get(user)
	return user, ok
}

// 通过email获取取用户
func (this *UserService) GetByEmail(email string) (*model.User, bool) {
	user := new(model.User)
	ok, _ := this.Db.Where("email=?", email).Get(user)
	return user, ok
}

// 用户是否存在
func (this *UserService) Exist(telphone string) bool {
	ok, _ := this.Db.Where("telphone=?", telphone).Exist((*model.User)(nil))
	return ok
}

// 获取验证 service
func (this *UserService) Verify() *VerifyService {
	verify := new(VerifyService)
	verify.Db = this.Db
	return verify
}

// 记录登录状态
// @path
func (this *UserService) GetUserAccess(ctx *gin.Context, user *model.User, action string) (*model.UserAccess, error) {
	access := model.NewAccessToken(user.Id, ctx)
	access.Action = action
	if _, err := this.Db.InsertOne(access); err != nil {
		return nil, err
	}
	switch action {
	case "login":
		// 下线同种设备
		this.Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=? and device=?", access.Id, access.UserId, access.Device)
	case "chpwd":
		// 注销所有
		this.Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=?", access.Id, access.UserId)
	}
	user.GetDetail(user)
	return model.NewUserAccess(user, access), nil
}

func CurrentUserMW() gin.HandlerFunc {
	fmt.Println("推荐使用 maker.AddParamManager(service.UserManager) 来替代 CurrentUserMW")
	return func(c *gin.Context) {
		// 当前登录用户数据
		token := c.Query("access_token")
		if c.Request.Method != http.MethodOptions && token != "" {
			now := time.Now()
			user := new(model.User)
			ok, _ := config.Db.Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
			if ok {
				c.Set("user", user)
			}
		}
	}
}

// 随机生成密码噪音
func RandNoise() string {
	return random.String(uint8(config.Cfg.Secure.NoiseLength))
}

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/yinrenxin/joeblog/models"
	"errors"
	"github.com/yinrenxin/joeblog/syserror"
	_ "github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
)

const SESSION_USER_KEY = "SESSION_USER_KEY"
type MAP_H = map[string]interface{}

type BaseController struct {
	beego.Controller
	User models.User
	IsLogin bool
	IsAdmin bool
}

func (this *BaseController) Prepare() {
	this.Data["Path"] = this.Ctx.Request.RequestURI
	u, ok := this.GetSession(SESSION_USER_KEY).(models.User)
	//logs.Info("登录id：",u.Id, "登录邮箱：",u.Email)
	this.IsLogin = false
	this.IsAdmin = false
	if ok {
		this.User = u
		this.IsLogin = true
		if u.Role == 0 {
			this.IsAdmin = true
		}
		this.Data["User"] = this.User
	}
	this.Data["islogin"] = this.IsLogin
	this.Data["isadmin"] = this.IsAdmin
}


func (this *BaseController) Abort500(err error) {
	this.Data["error"] = err
	this.Abort("500")
}

func (this *BaseController) Abort401(err error){
	this.Data["error"] = err
	this.Abort("401")
}


func (this *BaseController) GetMushString(key, msg string) string {
	k := this.GetString(key)
	if len(k) == 0 {
		this.Abort500(errors.New(msg))
	}
	return k
}

func (this *BaseController) MustLogin() {
	if !this.IsLogin {
		this.Abort500(syserror.NoUserError{})
	}
}

func (this *BaseController) JsonOK(msg, action string) {
	this.Data["json"] = MAP_H{
		"code": 0,
		"msg": msg,
		"action": "/",
	}
	this.ServeJSON()
}

func (this *BaseController) JsonOKH(msg string, data MAP_H) {
	data["code"] = 0
	data["msg"] = msg
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) UUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		this.Abort500(syserror.New("系统错误", nil))
	}
	return u.String()
}
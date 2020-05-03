package api

import (
	"fmt"
	"strconv"

	"layonserver/models"
	"layonserver/types"
	"layonserver/util/log"
	"layonserver/util/token"
	"layonserver/util/verifycode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) OnClose() {}

func (h *Handlers) UserAuthrize(c *gin.Context) {
	var (
		err    error
		userId string
	)
	//权限过滤器，如果验证成功继续，失败则goto errDeal
	if userId, err = token.TokenValidate(c.GetHeader("token")); err != nil {
		goto errDeal
	}
	c.Request.Header.Add("userid", userId)
	c.Next()
	return
errDeal:
	HandlerFailed(c, "登录超时，请重新登录！", err.Error())
	c.Abort()
	return
}

func (h Handlers) GetUserId(c *gin.Context) (uid uint, err error) {
	sid := c.Request.Header.Get("userid")
	id, err := strconv.ParseInt(sid, 10, 64)
	return uint(id), err
}

func (h Handlers) HandlerDelDevice(c *gin.Context) {
	var (
		err error
		req models.Devices
	)

	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}

	if req.UId, err = h.GetUserId(c); err != nil {
		goto errDeal
	}
	req.DelDevice()
	HandlerSuccess(c, "HandlerDelDevice", "success")
	return
errDeal:
	HandlerFailed(c, "删除设备失败：", err.Error())
}

func (h Handlers) HandlerGetDevices(c *gin.Context) {
	var (
		err   error
		req   models.Devices
		devs  []models.Devices
		total int64
	)
	type Rsp struct {
		Total int64            `json:"total"`
		Devs  []models.Devices `json:"devs"`
	}
	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}

	if req.UId, err = h.GetUserId(c); err != nil {
		goto errDeal
	}
	if devs, total, err = req.GetDevices(req.Page, req.Limit, req.UId); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGetDevices", Rsp{Total: total, Devs: devs})
	return
errDeal:
	HandlerFailed(c, "获取设备失败：", err.Error())
}

func (h Handlers) HandlerAddDevice(c *gin.Context) {
	var (
		err  error
		req  models.Devices
		user models.Users
	)

	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}

	if req.UId, err = h.GetUserId(c); err != nil {
		goto errDeal
	}

	if user, err = user.GetUserById(req.UId); err != nil {
		goto errDeal
	}

	if req.Count(req.UId) >= user.Count {
		err = fmt.Errorf("设备数量配额不足，请去官网增加配额")
		goto errDeal
	}

	if err = req.InsertDevice(); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerAddDevice", "success")
	return
errDeal:
	HandlerFailed(c, "添加设备失败：", err.Error())
}

func (h Handlers) HandlerGetUser(c *gin.Context) {
	var (
		err error
		req models.Users
	)
	if req.ID, err = h.GetUserId(c); err != nil {
		goto errDeal
	}
	if req, err = req.GetUserById(req.ID); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGetUser", req)
	return
errDeal:
	HandlerFailed(c, "获取用户信息失败：", err.Error())
}

func (h Handlers) HandlerGetUsers(c *gin.Context) {
	var (
		err   error
		req   models.Users
		users []models.Users
	)
	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}
	if req.ID, err = h.GetUserId(c); err != nil {
		goto errDeal
	}
	if req, err = req.GetUserById(req.ID); err != nil {
		goto errDeal
	}
	if req.Type != 2 {
		err = fmt.Errorf("用户无权限")
		goto errDeal
	}
	if users, err = req.GetUsers(req.Page, req.Limit, req.Type, req.FsId); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGetUsers", users)
	return
errDeal:
	HandlerFailed(c, "获取用户失败：", err.Error())
}

func (h Handlers) HandlerGetUsersByFather(c *gin.Context) {
	var (
		err    error
		req    models.Users
		users  []models.Users
		userId uint
	)
	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}
	if userId, err = h.GetUserId(c); err != nil {
		goto errDeal
	}
	if req, err = req.GetUserById(userId); err != nil {
		goto errDeal
	}
	if users, err = req.GetUsers(req.Page, req.Limit, 0, req.FsId); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGetUsersByFather", users)
	return
errDeal:
	HandlerFailed(c, "获取用户失败：", err.Error())
}

func (h Handlers) HandlerLogin(c *gin.Context) {
	var (
		err error
		req models.Users
		tok string
	)
	type RspToken struct {
		Token string `json:"token"`
	}
	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}
	if !verifycode.VCodeValidate(req.VerifyId, req.VerifyCode) {
		err = fmt.Errorf("验证码错误")
		goto errDeal
	}
	if req, err = req.GetUserByNamePass(); err != nil {
		err = fmt.Errorf("用户名或密码错误")
		goto errDeal
	}
	if tok, err = token.TokenGenerate(int64(req.ID)); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerLogin", RspToken{Token: tok})
	return
errDeal:
	HandlerFailed(c, "用户登录失败：", err.Error())
	return
}

func (h Handlers) HandlerRegist(c *gin.Context) {
	var (
		err error
		req models.Users
	)
	if err = c.ShouldBindWith(&req, binding.JSON); err != nil {
		goto errDeal
	}
	if !verifycode.VCodeValidate(req.VerifyId, req.VerifyCode) {
		err = fmt.Errorf("验证码错误")
		goto errDeal
	}
	if err = req.InsertUser(); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerRegist", "用户注册成功")
	return
errDeal:
	HandlerFailed(c, "用户注册失败：", err.Error())
	return
}

func (h *Handlers) HandlerGetVerifyCode(c *gin.Context) {
	var (
		err error
		rsp types.RspVerify
	)
	if rsp.VerifyId, rsp.Imag, err = verifycode.VCodeGenerate(4); err != nil {
		goto errDeal
	}
	HandlerSuccess(c, "HandlerGetVerifyCode", rsp)
	return
errDeal:
	HandlerFailed(c, "验证码生成失败：", err.Error())
	return
}

func HandlerSuccess(c *gin.Context, requestType, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "success",
		"data":    data,
	})
	logMsg := fmt.Sprintf("From [%s] result success", c.Request.RemoteAddr)
	log.GetLog().LogInfo(requestType, logMsg)
}

func HandlerFailed(c *gin.Context, requestType, errMsg string) {
	c.JSON(200, gin.H{
		"success": false,
		"msg":     requestType + errMsg,
	})
	logMsg := fmt.Sprintf("From [%s] result error [%s]", c.Request.RemoteAddr, errMsg)
	log.GetLog().LogError(requestType, logMsg)
}

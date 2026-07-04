package sdk

import (
	"strconv"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

// QuickLoginRequest 新版客户端 quick-login 请求
// 字段名从 APK global-metadata.dat 推断,支持多种命名
type QuickLoginRequest struct {
	UID2     int64  `json:"UID2,omitempty" form:"uid"`
	Uid      int64  `json:"uid,omitempty" form:"uid"`
	Token    string `json:"Token,omitempty" form:"token"`
	TransCode string `json:"TransCode,omitempty" form:"transcode"`
	DeviceID string `json:"DeviceID,omitempty" form:"deviceId"`
	Platform string `json:"Platform,omitempty" form:"platform"`
}

// quickLoginUserInfo 新版响应中的 UserInfo
type quickLoginUserInfo struct {
	ID         string `json:"ID"`
	UID2       int64  `json:"UID2"`
	PID        string `json:"PID"`
	Token      string `json:"Token"`
	Birthday   string `json:"Birthday"`
	RegChannel string `json:"RegChannel"`
	TransCode  string `json:"TransCode"`
	State      int32  `json:"State"`
	DeviceID   string `json:"DeviceID"`
	CreatedAt  int64  `json:"CreatedAt"`
}

// quickLoginYostar 新版响应中的 Yostar
type quickLoginYostar struct {
	ID        string `json:"ID"`
	Country   string `json:"Country"`
	Nickname  string `json:"Nickname"`
	Picture   string `json:"Picture"`
	State     int32  `json:"State"`
	AgreeAd   int32  `json:"AgreeAd"`
	CreatedAt int64  `json:"CreatedAt"`
}

// quickLoginKey 新版响应中的 Keys 元素
type quickLoginKey struct {
	ID        string `json:"ID"`
	Type      string `json:"Type"`
	Key       string `json:"Key"`
	NickName  string `json:"NickName"`
	CreatedAt int64  `json:"CreatedAt"`
}

// quickLoginData 新版响应中的 Data
type quickLoginData struct {
	AgeVerifyMethod int32               `json:"AgeVerifyMethod"`
	Destroy         interface{}         `json:"Destroy"`
	IsTestAccount   bool                `json:"IsTestAccount"`
	Keys            []quickLoginKey     `json:"Keys"`
	ServerNowAt     int64               `json:"ServerNowAt"`
	UserInfo        quickLoginUserInfo  `json:"UserInfo"`
	Yostar          quickLoginYostar    `json:"Yostar"`
	YostarDestroy   interface{}         `json:"YostarDestroy"`
}

// quickLoginResponse 新版 quick-login 响应
type quickLoginResponse struct {
	Code int32          `json:"Code"`
	Data quickLoginData `json:"Data"`
	Msg  string         `json:"Msg"`
}

// QuickLogin 新版客户端快速登录
// 对应 POST /user/quick-login
// 真实响应格式: {"Code":200,"Data":{"UserInfo":{"Token":"...","UID2":...},"Yostar":{...},"Keys":[...],"ServerNowAt":...},"Msg":"OK"}
func (s *SDK) QuickLogin(c *gin.Context) {
	// 读取并记录请求体(调试用)
	reqBody, _ := c.GetRawData()
	logger.Debug("quick-login 请求体: %s", string(reqBody))

	req := &QuickLoginRequest{}
	// 尝试 JSON 解析
	if len(reqBody) > 0 {
		json.Unmarshal(reqBody, req)
	}
	// 如果 JSON 解析失败,尝试 form 解析
	if req.UID2 == 0 && req.Uid == 0 {
		c.ShouldBind(req)
	}

	rsp := &quickLoginResponse{
		Code: 200,
		Msg:  "OK",
	}
	rsp.Data.ServerNowAt = time.Now().Unix()
	rsp.Data.AgeVerifyMethod = 0
	rsp.Data.IsTestAccount = false

	var yostarUid int64
	var yostarLoginToken string

	// 尝试用请求中的 uid 查找已有用户
	if req.UID2 != 0 || req.Uid != 0 {
		uid := req.UID2
		if uid == 0 {
			uid = req.Uid
		}
		yostarUser := db.GetDBGame().GetYostarUserByYostarUid(uid)
		if yostarUser != nil {
			yostarUid = yostarUser.YostarUid
			// 查找或创建 YostarUserLogin
			yoStarUserLogin := db.GetDBGame().GetYoStarUserLoginByYostarUid(yostarUid)
			if yoStarUserLogin == nil {
				yoStarUserLogin, _ = db.GetDBGame().AddYoStarUserLoginByYostarUid(yostarUid)
			}
			if yoStarUserLogin != nil {
				yostarLoginToken = alg.RandStr(40)
				yoStarUserLogin.YostarLoginToken = yostarLoginToken
				db.GetDBGame().UpdateYoStarUserLogin(yoStarUserLogin)
			}
		}
	}

	// 如果没有找到用户,自动注册新用户
	if yostarUid == 0 {
		// 创建 YostarAccount
		yostarAccount, err := db.GetDBGame().AddYostarAccountByYostarAccount("guest_" + strconv.FormatInt(time.Now().UnixNano(), 10))
		if err != nil || yostarAccount == nil {
			logger.Debug("quick-login 自动注册 YostarAccount 失败: %s", err)
			rsp.Code = 500
			rsp.Msg = "server error"
			c.JSON(200, rsp)
			return
		}
		yostarUid = yostarAccount.YostarUid

		// 创建 YostarUser
		yostarUser, err := db.GetDBGame().AddYostarUserByYostarUid(yostarUid)
		if err != nil || yostarUser == nil {
			logger.Debug("quick-login 自动注册 YostarUser 失败: %s", err)
		}

		// 创建 YostarUserLogin
		yoStarUserLogin, err := db.GetDBGame().AddYoStarUserLoginByYostarUid(yostarUid)
		if err != nil || yoStarUserLogin == nil {
			logger.Debug("quick-login 自动注册 YostarUserLogin 失败: %s", err)
		} else {
			yostarLoginToken = alg.RandStr(40)
			yoStarUserLogin.YostarLoginToken = yostarLoginToken
			db.GetDBGame().UpdateYoStarUserLogin(yoStarUserLogin)
		}

		logger.Debug("quick-login 自动注册新用户 YostarUid: %d", yostarUid)
	}

	// 填充响应
	rsp.Data.UserInfo = quickLoginUserInfo{
		ID:         strconv.FormatInt(yostarUid, 10),
		UID2:       yostarUid,
		PID:        "JP-BA",
		Token:      yostarLoginToken,
		Birthday:   "",
		RegChannel: "guest",
		TransCode:  alg.RandStr(15),
		State:      1,
		DeviceID:   req.DeviceID,
		CreatedAt:  time.Now().Unix(),
	}
	rsp.Data.Yostar = quickLoginYostar{
		ID:        "Y" + strconv.FormatInt(yostarUid, 10),
		Country:   "JP",
		Nickname:  "guest_" + strconv.FormatInt(yostarUid, 10),
		Picture:   "",
		State:     1,
		AgreeAd:   0,
		CreatedAt: time.Now().Unix(),
	}
	rsp.Data.Keys = []quickLoginKey{
		{
			ID:        strconv.FormatInt(yostarUid, 10),
			Type:      "yostar",
			Key:       "guest_" + strconv.FormatInt(yostarUid, 10) + "@baps.local",
			NickName:  "guest_" + strconv.FormatInt(yostarUid, 10),
			CreatedAt: time.Now().Unix(),
		},
	}

	logger.Debug("quick-login 成功 YostarUid: %d, Token: %s", yostarUid, yostarLoginToken)
	c.JSON(200, rsp)
}

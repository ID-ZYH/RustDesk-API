package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
)

type UserDevice struct{}

func (ct *UserDevice) List(c *gin.Context) {
	query := &admin.UserDeviceListQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	res := service.AllService.UserService.TokenList(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.UserId > 0 {
			tx.Where("user_id = ?", query.UserId)
		}
		tx.Order("id desc")
	})

	items := make([]gin.H, 0, len(res.UserTokens))
	if len(res.UserTokens) > 0 {
		ids := make([]uint, 0, len(res.UserTokens))
		for _, token := range res.UserTokens {
			ids = append(ids, token.Id)
		}
		logs := make([]model.LoginLog, 0)
		global.DB.Where("user_token_id in ?", ids).Order("id desc").Find(&logs)
		logMap := make(map[uint]model.LoginLog)
		for _, logItem := range logs {
			if _, ok := logMap[logItem.UserTokenId]; !ok {
				logMap[logItem.UserTokenId] = logItem
			}
		}

		for _, token := range res.UserTokens {
			logItem := logMap[token.Id]
			items = append(items, gin.H{
				"id":          token.Id,
				"user_id":     token.UserId,
				"device_uuid": token.DeviceUuid,
				"device_id":   token.DeviceId,
				"token":       token.Token,
				"expired_at":  token.ExpiredAt,
				"created_at":  token.CreatedAt,
				"updated_at":  token.UpdatedAt,
				"client":      logItem.Client,
				"platform":    logItem.Platform,
				"ip":          logItem.Ip,
				"login_at":    logItem.CreatedAt,
			})
		}
	}

	response.Success(c, gin.H{
		"list":      items,
		"total":     res.Total,
		"page":      res.Page,
		"page_size": res.PageSize,
	})
}

func (ct *UserDevice) SetLimit(c *gin.Context) {
	f := &admin.UserDeviceSetLimitForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}

	u := service.AllService.UserService.InfoById(f.UserId)
	if u.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}

	u.MaxDevices = f.MaxDevices
	if err := global.DB.Model(u).Update("max_devices", f.MaxDevices).Error; err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, gin.H{
		"user_id":     u.Id,
		"max_devices": f.MaxDevices,
	})
}

func (ct *UserDevice) Unbind(c *gin.Context) {
	f := &admin.UserDeviceUnbindForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}

	token := service.AllService.UserService.TokenInfoById(f.Id)
	if token.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	if err := service.AllService.UserService.DeleteToken(token); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

func (ct *UserDevice) BatchUnbind(c *gin.Context) {
	f := &admin.UserDeviceBatchUnbindForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	if err := service.AllService.UserService.BatchDeleteUserToken(f.Ids); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

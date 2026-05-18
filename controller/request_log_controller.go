package controller

import (
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

// GetRequestDetail 获取请求/响应详情
// GET /api/log/request_detail?request_id=xxx
func GetRequestDetail(c *gin.Context) {
	requestId := c.Query("request_id")
	if requestId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "request_id is required",
		})
		return
	}

	userId := c.GetInt("id")

	entry, err := service.SearchRequestLog(requestId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "request log not found",
		})
		return
	}

	// 权限检查：非管理员只能查看自己 Token 的日志
	if !model.IsAdmin(userId) {
		if entry.UserId != userId {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "no permission to view this log",
			})
			return
		}
	}

	common.ApiSuccess(c, entry)
}

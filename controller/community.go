package controller

import (
	"errors"
	"strconv"

	"bluebell/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ----社区相关----
// CommunityHandler 处理社区列表请求
// @Summary 社区列表
// @Description 查询所有社区（id、name）列表
// @Tags 社区相关接口
// @Produce json
// @Success 200 {object} controller.Response{data=[]models.Community}
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	//查询到所有社区（community_id,community_name）列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler 处理社区详情请求
// @Summary 社区详情
// @Description 根据社区 ID 查询社区详细信息
// @Tags 社区相关接口
// @Produce json
// @Param id path int true "社区 ID"
// @Success 200 {object} controller.Response{data=models.CommunityDetail}
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
	communityID := c.Param("id")
	//做一个参数校验
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	//查询到社区详细信息（community_id,community_name）返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		// 区分业务错误和系统错误
		if errors.Is(err, logic.ErrCommunityNotExist) {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

package controller

import (
	"errors"

	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PostVoteController 投票接口 POST /api/v1/vote
// @Summary 帖子投票
// @Description 对帖子进行投票，direction 取值：1 赞成 | 0 取消 | -1 反对，需登录后携带 Token
// @Tags 投票相关接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param object body models.VoteData true "投票参数"
// @Success 200 {object} controller.Response
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	// 1. 参数校验
	p := new(models.VoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 获取当前用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 3. 投票（帖子存在性、投票时限校验在redis层完成）
	if err := logic.VoteForPost(userID, p.PostID, p.Direction); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		if errors.Is(err, logic.ErrPostNotExist) {
			ResponseError(c, CodePostNotExist)
			return
		}
		if errors.Is(err, logic.ErrVoteTimeExpire) {
			ResponseError(c, CodeVoteTimeExpire)
			return
		}
		if errors.Is(err, logic.ErrRepeatedVote) {
			ResponseError(c, CodeRepeatedVote)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 4. 返回
	ResponseSuccess(c, nil)
}

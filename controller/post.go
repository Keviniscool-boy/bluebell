package controller

import (
	"errors"
	"strconv"

	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新帖子，需登录后携带 Token
// @Tags 帖子相关接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param object body models.Post true "帖子参数"
// @Success 200 {object} controller.Response
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数校验
	p := new(models.Post)
	err := c.ShouldBindJSON(p)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	//从context中获取当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 帖子详情
// @Summary 帖子详情
// @Description 根据帖子 ID 查询帖子详情（含作者名、社区信息）
// @Tags 帖子相关接口
// @Produce json
// @Param id path int true "帖子 ID"
// @Success 200 {object} controller.Response{data=models.ApiPostDetail}
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	//获取帖子id
	postID := c.Param("id")
	//做一个参数校验
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	//查询到帖子详细信息返回
	data, err := logic.GetPostDetail(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() failed", zap.Error(err))
		// 区分业务错误和系统错误
		if errors.Is(err, logic.ErrPostNotExist) {
			ResponseError(c, CodePostNotExist)
			return
		}
		if errors.Is(err, logic.ErrCommunityNotExist) {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数 GET /api/v1/posts?page=1&size=10
// @Summary 帖子列表
// @Description 分页查询帖子列表
// @Tags 帖子相关接口
// @Produce json
// @Param page query int false "页码（从 1 开始，默认 1）"
// @Param size query int false "每页数量（默认 10）"
// @Success 200 {object} controller.Response{data=[]models.ApiPostDetail}
// @Router /posts/ [get]
func GetPostListHandler(c *gin.Context) {
	p, ok := getPostListParam(c)
	if !ok {
		return
	}
	//获取数据
	data, err := logic.GetPostList(p.Page, p.Size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回
	ResponseSuccess(c, data)
}

// GetPostList2Handler 按热度/时间排序的帖子列表 GET /api/v1/posts2?order=score|time&page=1&size=10
// @Summary 帖子列表（排序）
// @Description 按时间或热度排序分页查询帖子列表
// @Tags 帖子相关接口
// @Produce json
// @Param order query string false "排序方式：time 按时间（默认）| score 按热度"
// @Param page query int false "页码（从 1 开始，默认 1）"
// @Param size query int false "每页数量（默认 10）"
// @Success 200 {object} controller.Response{data=[]models.ApiPostDetail}
// @Router /posts2/ [get]
func GetPostList2Handler(c *gin.Context) {
	p, ok := getPostListParam(c)
	if !ok {
		return
	}
	data, err := logic.GetPostList2(p.Page, p.Size, p.Order)
	if err != nil {
		zap.L().Error("logic.GetPostList2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// getPostListParam 绑定并填充帖子列表查询参数默认值；绑定失败时已写错误响应并返回 ok=false
func getPostListParam(c *gin.Context) (*models.ParamPostList, bool) {
	p := new(models.ParamPostList)
	if err := c.ShouldBindQuery(p); err != nil {
		ResponseError(c, CodeInvalidParams)
		return p, false
	}
	// 填充默认值
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Order != "score" {
		p.Order = "time"
	}
	return p, true
}

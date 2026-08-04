package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidCredentials
	CodePasswordMismatch
	CodeServerBusy
	CodePostNotExist
	CodeCommunityNotExist

	CodeInvalidAuth
	CodeNeedLogin

	CodeVoteTimeExpire
	CodeRepeatedVote
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "成功",
	CodeInvalidParams:      "请求参数错误",
	CodeUserExist:          "用户名已存在",
	CodeUserNotExist:       "用户不存在",
	CodeInvalidCredentials: "用户名或密码错误",
	CodePasswordMismatch:   "密码不匹配",
	CodeServerBusy:         "服务器繁忙",
	CodePostNotExist:       "帖子不存在",
	CodeCommunityNotExist:  "社区不存在",
	CodeInvalidAuth:        "认证失败",
	CodeNeedLogin:          "需要登录",
	CodeVoteTimeExpire:     "投票时间已过",
	CodeRepeatedVote:       "请勿重复投票",
}

// 创建一个函数，传入code返回msg
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return "未知错误"
	}
	return msg
}

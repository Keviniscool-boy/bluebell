package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"
)

// 业务错误定义
var ErrCommunityNotExist = errors.New("社区不存在")

func GetCommunityList() (data []*models.Community, err error) {
	//查数据库，查到所有社区（community_id,community_name）列表形式返回
	//返回给前端
	return mysql.GetCommunityLList()

}
func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	community, err := mysql.GetCommunityDetailByID(id)
	if err != nil {
		if errors.Is(err, mysql.ErrCommunityNotExist) {
			return nil, ErrCommunityNotExist
		}
		return nil, err
	}
	return community, nil
}

package mysql

import (
	"database/sql"
	"errors"

	"bluebell/models"

	"go.uber.org/zap"
)

var ErrCommunityNotExist = errors.New("社区不存在")

func GetCommunityLList() (data []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err = db.Select(&data, sqlStr); err != nil {
		zap.L().Error("query community failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := "select community_id,community_name,introduction,create_time from community where community_id=?"
	if err = db.Get(community, sqlStr, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("community not found", zap.Int64("id", id))
			return nil, ErrCommunityNotExist
		}
		zap.L().Error("query community detail failed", zap.Error(err))
		return nil, err
	}
	return community, nil

}

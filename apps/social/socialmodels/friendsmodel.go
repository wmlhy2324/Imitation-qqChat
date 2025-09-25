package socialmodels

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface {
		friendsModel
		// 搜索好友：根据用户ID和关键词搜索好友列表
		SearchFriends(ctx context.Context, userId, keyword string) ([]*Friends, error)
	}

	customFriendsModel struct {
		*defaultFriendsModel
	}
)

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c),
	}
}

// SearchFriends 搜索好友：根据用户ID和关键词搜索好友列表
// 这里我们通过用户的好友关系表进行搜索，然后需要结合用户信息进行筛选
func (m *customFriendsModel) SearchFriends(ctx context.Context, userId, keyword string) ([]*Friends, error) {
	// 获取当前用户的所有好友
	friends, err := m.ListByUserid(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 这里返回所有好友，具体的搜索筛选在上层逻辑中通过用户服务完成
	// 因为需要根据好友的昵称、手机号等用户信息进行搜索
	return friends, nil
}

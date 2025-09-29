package socialmodels

import (
	"context"
	"database/sql"

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
		// 删除好友：删除双向好友关系
		DeleteFriend(ctx context.Context, userId, friendUid string) error
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

// DeleteFriend 删除好友：删除双向好友关系
func (m *customFriendsModel) DeleteFriend(ctx context.Context, userId, friendUid string) error {
	// 删除双向好友关系，需要删除两条记录：userId->friendUid 和 friendUid->userId
	query := `DELETE FROM friends WHERE (user_id = ? AND friend_uid = ?) OR (user_id = ? AND friend_uid = ?)`
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, userId, friendUid, friendUid, userId)
	})
	return err
}

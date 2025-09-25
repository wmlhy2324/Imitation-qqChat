package socialmodels

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupsModel = (*customGroupsModel)(nil)

type (
	// GroupsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupsModel.
	GroupsModel interface {
		groupsModel
		// 搜索群聊：根据群名关键词搜索群列表
		SearchGroups(ctx context.Context, keyword string) ([]*Groups, error)
		// 根据用户ID获取用户加入的群列表
		ListByUserId(ctx context.Context, userId string) ([]*Groups, error)
	}

	customGroupsModel struct {
		*defaultGroupsModel
	}
)

// NewGroupsModel returns a model for the database table.
func NewGroupsModel(conn sqlx.SqlConn, c cache.CacheConf) GroupsModel {
	return &customGroupsModel{
		defaultGroupsModel: newGroupsModel(conn, c),
	}
}

// SearchGroups 搜索群聊：根据群名关键词搜索群列表
func (m *customGroupsModel) SearchGroups(ctx context.Context, keyword string) ([]*Groups, error) {
	var resp []*Groups
	query := fmt.Sprintf("select %s from %s where name like ? order by created_at desc", groupsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, fmt.Sprintf("%%%s%%", keyword))
	return resp, err
}

// ListByUserId 根据用户ID获取用户加入的群列表
// 需要关联group_members表查询用户加入的群
func (m *customGroupsModel) ListByUserId(ctx context.Context, userId string) ([]*Groups, error) {
	var resp []*Groups
	query := fmt.Sprintf(`
		select %s from %s g 
		inner join group_members gm on g.id = gm.group_id 
		where gm.user_id = ? 
		order by gm.join_time desc
	`, groupsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	return resp, err
}

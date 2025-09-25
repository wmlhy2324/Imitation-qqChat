package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/social/rpc/internal/config"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	socialmodels.FriendsModel
	socialmodels.FriendRequestsModel
	socialmodels.GroupsModel
	socialmodels.GroupRequestsModel
	socialmodels.GroupMembersModel

	// RPC客户端
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		FriendsModel:        socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(sqlConn, c.Cache),

		// 初始化RPC客户端
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

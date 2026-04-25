package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/im/api/internal/config"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	imclient.Im
	userclient.User
	socialclient.Social

	immodels.UserMessageDeletesModel
	immodels.ConversationsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Im:                      imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		User:                    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social:                  socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		UserMessageDeletesModel: immodels.MustUserMessageDeletesModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel:      immodels.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
	}
}

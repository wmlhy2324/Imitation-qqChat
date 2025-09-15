/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package handler

import (
	"imooc.com/easy-chat/apps/im/ws/internal/handler/conversation"
	"imooc.com/easy-chat/apps/im/ws/internal/handler/push"
	"imooc.com/easy-chat/apps/im/ws/internal/handler/user"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "conversation.markChat",
			Handler: conversation.MarkRead(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}

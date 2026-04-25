package conversation

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/task/mq/mq"
	"imooc.com/easy-chat/pkg/constants"
)

func Revoke(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Revoke
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 轻量校验：有 conversationId
		if data.ConversationId == "" || data.MsgId == "" {
			srv.Send(websocket.NewErrMessage(errors.New("msgId and conversationId are required")), conn)
			return
		}

		// 遵循 IM-WS 异步解耦原则：只写 Kafka，不查询 MongoDB
		err := svc.MsgRevokeTransferClient.Push(&mq.MsgRevokeTransfer{
			MsgId:          data.MsgId,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			ChatType:       int32(constants.GroupChatType), // 由 Task-MQ 查 ChatLog 时获取真实类型
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}

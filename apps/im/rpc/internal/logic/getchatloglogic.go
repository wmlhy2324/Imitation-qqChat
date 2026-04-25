package logic

import (
	"context"

	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/pkg/ctxdata"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// todo: add your logic here and delete this line

	// 根据id
	if in.MsgId != "" {
		chatlog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog by msgId err %v, req %v", err, in.MsgId)
		}

		return &im.GetChatLogResp{
			List: []*im.ChatLog{{
				Id:             chatlog.ID.Hex(),
				ConversationId: chatlog.ConversationId,
				SendId:         chatlog.SendId,
				RecvId:         chatlog.RecvId,
				MsgType:        int32(chatlog.MsgType),
				MsgContent:     chatlog.MsgContent,
				ChatType:       int32(chatlog.ChatType),
				SendTime:       chatlog.SendTime,
				ReadRecords:    chatlog.ReadRecords,
			}},
		}, nil
	}

	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog list by SendTime err %v, req %v", err, in)
	}

	// 过滤用户已删除的消息
	userId := ctxdata.GetUId(l.ctx)
	if userId != "" {
		data = l.filterDeletedMessages(l.ctx, userId, in.ConversationId, data)
	}

	res := make([]*im.ChatLog, 0, len(data))
	for _, datum := range data {
		res = append(res, &im.ChatLog{
			Id:             datum.ID.Hex(),
			ConversationId: datum.ConversationId,
			SendId:         datum.SendId,
			RecvId:         datum.RecvId,
			MsgType:        int32(datum.MsgType),
			MsgContent:     datum.MsgContent,
			ChatType:       int32(datum.ChatType),
			SendTime:       datum.SendTime,
			ReadRecords:    datum.ReadRecords,
		})
	}

	return &im.GetChatLogResp{
		List: res,
	}, nil
}

// filterDeletedMessages 过滤掉用户已删除的消息
func (l *GetChatLogLogic) filterDeletedMessages(ctx context.Context, userId, conversationId string, chatLogs []*immodels.ChatLog) []*immodels.ChatLog {
	if len(chatLogs) == 0 {
		return chatLogs
	}

	// 查询被删除的消息 ID 列表
	deletedRecords, err := l.svcCtx.UserMessageDeletesModel.ListByUserIdAndConversation(ctx, userId, conversationId)
	if err != nil {
		l.Errorf("filterDeletedMessages ListByUserIdAndConversation err %v", err)
		return chatLogs
	}

	// 构建删除 ID 集合
	deletedSet := make(map[string]bool, len(deletedRecords))
	for _, dr := range deletedRecords {
		deletedSet[dr.MsgId] = true
	}

	// 查询 ClearUpTo
	var clearUpTo int64
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(ctx, userId)
	if err == nil && conversations != nil {
		if conv, ok := conversations.ConversationList[conversationId]; ok {
			clearUpTo = conv.ClearUpTo
		}
	}

	// 过滤
	filtered := make([]*immodels.ChatLog, 0, len(chatLogs))
	for _, cl := range chatLogs {
		if deletedSet[cl.ID.Hex()] {
			continue
		}
		if clearUpTo > 0 && cl.SendTime <= clearUpTo {
			continue
		}
		filtered = append(filtered, cl)
	}

	return filtered
}

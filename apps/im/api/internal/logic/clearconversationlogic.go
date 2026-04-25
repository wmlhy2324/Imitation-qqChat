package logic

import (
	"context"
	"errors"
	"time"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearConversationLogic {
	return &ClearConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ClearConversation 清空会话聊天记录（仅本端隐藏，数据库保留）
func (l *ClearConversationLogic) ClearConversation(req *types.ClearConversationReq) (resp *types.ClearConversationResp, err error) {
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		return nil, errors.New("用户未登录")
	}

	// 查询用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("ClearConversation FindByUserId err %v", err)
		return nil, err
	}

	if conversations.ConversationList == nil {
		conversations.ConversationList = make(map[string]*immodels.Conversation)
	}

	// 更新 clearUpTo 时间戳
	now := time.Now().UnixMilli()
	if conv, ok := conversations.ConversationList[req.ConversationId]; ok {
		conv.ClearUpTo = now
		conv.IsShow = false
	} else {
		conversations.ConversationList[req.ConversationId] = &immodels.Conversation{
			ConversationId: req.ConversationId,
			ClearUpTo:      now,
			IsShow:         false,
		}
	}

	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		l.Errorf("ClearConversation Update err %v", err)
		return nil, err
	}

	// 清理单条删除记录（ClearUpTo 已覆盖）
	err = l.svcCtx.UserMessageDeletesModel.DeleteByConversation(l.ctx, userId, req.ConversationId)
	if err != nil {
		l.Errorf("ClearConversation DeleteByConversation err %v", err)
	}

	return &types.ClearConversationResp{}, nil
}

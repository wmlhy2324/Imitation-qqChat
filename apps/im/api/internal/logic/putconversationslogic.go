package logic

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PutConversations 更新用户会话信息
// 接收前端传来的会话更新数据，包括已读消息数、显示状态等
// 转换数据格式后调用RPC服务更新用户的会话记录
func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (resp *types.PutConversationsResp, err error) {
	// 从JWT上下文中获取当前用户ID
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		l.Error("用户未登录")
		return nil, errors.New("用户未登录")
	}

	// 构建RPC请求
	rpcReq := &imclient.PutConversationsReq{
		UserId:           userId,
		ConversationList: make(map[string]*imclient.Conversation),
	}

	// 转换API请求到RPC请求格式
	if req.ConversationList != nil {
		for conversationId, conversation := range req.ConversationList {
			rpcConversation := &imclient.Conversation{}
			err = copier.Copy(rpcConversation, conversation)
			if err != nil {
				l.Errorf("复制会话数据失败: %v", err)
				return nil, err
			}
			rpcReq.ConversationList[conversationId] = rpcConversation
		}
	}

	// 调用RPC服务更新会话
	_, err = l.svcCtx.Im.PutConversations(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("更新会话失败: %v", err)
		return nil, err
	}

	// 返回成功响应
	resp = &types.PutConversationsResp{}
	return resp, nil
}

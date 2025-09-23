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

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetConversations 获取用户会话列表
// 从JWT token中获取当前用户ID，然后调用RPC服务获取该用户的所有会话
// 包括会话ID、聊天类型、是否显示、消息统计等信息
func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	// 从JWT上下文中获取当前用户ID
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		l.Error("用户未登录")
		return nil, errors.New("用户未登录")
	}

	// 调用RPC服务获取用户会话列表
	rpcResp, err := l.svcCtx.Im.GetConversations(l.ctx, &imclient.GetConversationsReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("获取会话列表失败: %v", err)
		return nil, err
	}

	// 构建API响应
	resp = &types.GetConversationsResp{
		UserId:           userId,
		ConversationList: make(map[string]*types.Conversation),
	}

	// 转换RPC响应到API响应格式
	if rpcResp.ConversationList != nil {
		for conversationId, conversation := range rpcResp.ConversationList {
			apiConversation := &types.Conversation{}
			err = copier.Copy(apiConversation, conversation)
			if err != nil {
				l.Errorf("复制会话数据失败: %v", err)
				return nil, err
			}
			resp.ConversationList[conversationId] = apiConversation
		}
	}

	return resp, nil
}

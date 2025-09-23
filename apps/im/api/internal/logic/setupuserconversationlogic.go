package logic

import (
	"context"
	"errors"

	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SetUpUserConversation 建立用户会话（私聊或群聊）
// 根据聊天类型创建相应的会话，对于私聊需要发送者和接收者ID
// 对于群聊需要群组ID，调用RPC服务在数据库中建立会话记录
func (l *SetUpUserConversationLogic) SetUpUserConversation(req *types.SetUpUserConversationReq) (resp *types.SetUpUserConversationResp, err error) {
	// 从JWT上下文中获取当前用户ID，作为发送者
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		l.Error("用户未登录")
		return nil, errors.New("用户未登录")
	}

	// 参数验证
	if req.ChatType == 0 {
		l.Error("聊天类型不能为空")
		return nil, errors.New("聊天类型不能为空")
	}

	// 构建RPC请求
	rpcReq := &imclient.SetUpUserConversationReq{
		SendId:   userId, // 当前登录用户作为发送者
		ChatType: req.ChatType,
	}

	// 根据聊天类型设置不同的参数
	switch req.ChatType {
	case 1: // 私聊
		if req.RecvId == "" {
			l.Error("私聊接收者ID不能为空")
			return nil, errors.New("私聊接收者ID不能为空")
		}
		if req.RecvId == userId {
			l.Error("不能和自己聊天")
			return nil, errors.New("不能和自己聊天")
		}
		rpcReq.RecvId = req.RecvId
	case 2: // 群聊
		if req.RecvId == "" {
			l.Error("群聊ID不能为空")
			return nil, errors.New("群聊ID不能为空")
		}
		rpcReq.RecvId = req.RecvId // 对于群聊，RecvId表示群组ID
	default:
		l.Error("不支持的聊天类型")
		return nil, errors.New("不支持的聊天类型")
	}

	// 如果请求中指定了SendId，使用请求中的SendId（通常用于管理员操作）
	if req.SendId != "" {
		rpcReq.SendId = req.SendId
	}

	// 调用RPC服务建立会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("建立会话失败: %v", err)
		return nil, err
	}

	// 返回成功响应
	resp = &types.SetUpUserConversationResp{}
	return resp, nil
}

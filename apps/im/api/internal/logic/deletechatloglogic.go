package logic

import (
	"context"
	"errors"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteChatLogLogic {
	return &DeleteChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteChatLog 删除单条聊天记录（仅本端隐藏，数据库保留）
func (l *DeleteChatLogLogic) DeleteChatLog(req *types.DeleteChatLogReq) (resp *types.DeleteChatLogResp, err error) {
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		return nil, errors.New("用户未登录")
	}

	// 插入删除记录
	err = l.svcCtx.UserMessageDeletesModel.Insert(l.ctx, &immodels.UserMessageDelete{
		UserId:         userId,
		ConversationId: req.ConversationId,
		MsgId:          req.MsgId,
	})
	if err != nil {
		l.Errorf("DeleteChatLog Insert err %v", err)
		return nil, err
	}

	return &types.DeleteChatLogResp{}, nil
}

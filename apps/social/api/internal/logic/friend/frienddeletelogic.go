package friend

import (
	"context"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除好友
func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendDeleteLogic) FriendDelete(req *types.FriendDeleteReq) (resp *types.FriendDeleteResp, err error) {
	// 获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)
	if uid == "" {
		return nil, err
	}

	// 调用RPC服务删除好友
	_, err = l.svcCtx.Social.FriendDelete(l.ctx, &socialclient.FriendDeleteReq{
		UserId:    uid,
		FriendUid: req.FriendUid,
	})
	if err != nil {
		return nil, err
	}

	return &types.FriendDeleteResp{}, nil
}

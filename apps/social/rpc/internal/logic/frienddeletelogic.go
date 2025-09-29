package logic

import (
	"context"

	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendDeleteLogic) FriendDelete(in *social.FriendDeleteReq) (*social.FriendDeleteResp, error) {
	// 参数验证
	if in.UserId == "" || in.FriendUid == "" {
		return nil, errors.Wrapf(xerr.NewDBErr(), "invalid params: userId=%s, friendUid=%s", in.UserId, in.FriendUid)
	}

	// 删除好友关系
	err := l.svcCtx.FriendsModel.DeleteFriend(l.ctx, in.UserId, in.FriendUid)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "delete friend failed: %v", err)
	}

	return &social.FriendDeleteResp{}, nil
}

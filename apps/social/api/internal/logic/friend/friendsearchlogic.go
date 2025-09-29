package friend

import (
	"context"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索好友
func NewFriendSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendSearchLogic {
	return &FriendSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendSearchLogic) FriendSearch(req *types.FriendSearchReq) (resp *types.FriendSearchResp, err error) {
	// 获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)
	if uid == "" {
		return nil, err
	}

	// 调用RPC服务搜索好友
	searchResp, err := l.svcCtx.Social.FriendSearch(l.ctx, &socialclient.FriendSearchReq{
		UserId:  uid,
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	var resultList []*types.SearchFriendInfo
	for _, friend := range searchResp.List {
		resultList = append(resultList, &types.SearchFriendInfo{
			Id:       friend.Id,
			Nickname: friend.Nickname,
			Avatar:   friend.Avatar,
			Phone:    friend.Phone,
			Remark:   friend.Remark,
		})
	}

	return &types.FriendSearchResp{
		List: resultList,
	}, nil
}

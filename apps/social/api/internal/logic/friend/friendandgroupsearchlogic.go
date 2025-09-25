package friend

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendAndGroupSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索好友和群聊
func NewFriendAndGroupSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendAndGroupSearchLogic {
	return &FriendAndGroupSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendAndGroupSearch 搜索好友和群聊
// 根据输入的关键词，在用户的好友列表和群聊中进行模糊搜索
// 好友搜索基于昵称和手机号，群聊搜索基于群名称
// 返回匹配的好友列表和群聊列表，包括加入状态信息
func (l *FriendAndGroupSearchLogic) FriendAndGroupSearch(req *types.FriendAndGroupSearchReq) (resp *types.FriendAndGroupSearchResp, err error) {
	// 参数验证
	if req.Keyword == "" {
		return &types.FriendAndGroupSearchResp{
			FriendList: []*types.SearchFriendInfo{},
			GroupList:  []*types.SearchGroupInfo{},
		}, nil
	}

	// 从JWT上下文中获取当前用户ID
	userId := ctxdata.GetUId(l.ctx)
	if userId == "" {
		l.Error("用户未登录")
		return nil, errors.New("用户未登录")
	}

	// 调用RPC服务进行搜索
	rpcResp, err := l.svcCtx.Social.FriendAndGroupSearch(l.ctx, &socialclient.FriendAndGroupSearchReq{
		UserId:  userId,
		Keyword: req.Keyword,
	})
	if err != nil {
		l.Errorf("搜索好友和群聊失败: %v", err)
		return nil, err
	}

	// 构建API响应
	resp = &types.FriendAndGroupSearchResp{
		FriendList: make([]*types.SearchFriendInfo, 0),
		GroupList:  make([]*types.SearchGroupInfo, 0),
	}

	// 转换好友搜索结果
	if rpcResp.FriendList != nil {
		for _, rpcFriend := range rpcResp.FriendList {
			apiFriend := &types.SearchFriendInfo{}
			err = copier.Copy(apiFriend, rpcFriend)
			if err != nil {
				l.Errorf("复制好友数据失败: %v", err)
				continue
			}
			resp.FriendList = append(resp.FriendList, apiFriend)
		}
	}

	// 转换群聊搜索结果
	if rpcResp.GroupList != nil {
		for _, rpcGroup := range rpcResp.GroupList {
			apiGroup := &types.SearchGroupInfo{}
			err = copier.Copy(apiGroup, rpcGroup)
			if err != nil {
				l.Errorf("复制群聊数据失败: %v", err)
				continue
			}
			resp.GroupList = append(resp.GroupList, apiGroup)
		}
	}

	l.Infof("搜索完成，找到 %d 个好友，%d 个群聊", len(resp.FriendList), len(resp.GroupList))
	return resp, nil
}

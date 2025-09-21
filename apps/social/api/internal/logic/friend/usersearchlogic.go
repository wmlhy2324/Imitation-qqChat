package friend

import (
	"context"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
	"imooc.com/easy-chat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索用户
func NewUserSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSearchLogic {
	return &UserSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSearchLogic) UserSearch(req *types.UserSearchReq) (resp *types.UserSearchResp, err error) {
	// 参数校验
	if req.Keyword == "" {
		return &types.UserSearchResp{
			List: []*types.UserInfo{},
		}, nil
	}

	// 用map去重，防止同一用户在昵称和手机号搜索中重复出现
	userMap := make(map[string]*userclient.UserEntity)

	// 1. 按昵称搜索
	nameUsers, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Name: req.Keyword,
	})
	if err != nil {
		l.Errorf("按昵称搜索用户失败: %v", err)
	} else {
		for _, user := range nameUsers.User {
			userMap[user.Id] = user
		}
	}

	// 2. 按手机号搜索
	phoneUsers, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Phone: req.Keyword,
	})
	if err != nil {
		l.Errorf("按手机号搜索用户失败: %v", err)
	} else {
		for _, user := range phoneUsers.User {
			userMap[user.Id] = user
		}
	}

	// 如果两次搜索都失败了，返回错误
	if len(userMap) == 0 && err != nil {
		return nil, err
	}

	// 构造响应数据
	respList := make([]*types.UserInfo, 0, len(userMap))
	for _, user := range userMap {
		userInfo := &types.UserInfo{
			Id:       user.Id,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
		}
		respList = append(respList, userInfo)
	}

	return &types.UserSearchResp{
		List: respList,
	}, nil
}

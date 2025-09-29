package logic

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"imooc.com/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendSearchLogic {
	return &FriendSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendSearchLogic) FriendSearch(in *social.FriendSearchReq) (*social.FriendSearchResp, error) {
	// 参数验证
	if in.UserId == "" {
		return nil, errors.Wrapf(xerr.NewDBErr(), "invalid params: userId=%s", in.UserId)
	}

	// 如果没有搜索关键词，返回空结果
	if strings.TrimSpace(in.Keyword) == "" {
		return &social.FriendSearchResp{List: []*social.SearchFriendInfo{}}, nil
	}

	// 获取当前用户的所有好友
	friends, err := l.svcCtx.FriendsModel.ListByUserid(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get friends failed: %v", err)
	}

	if len(friends) == 0 {
		return &social.FriendSearchResp{List: []*social.SearchFriendInfo{}}, nil
	}

	// 获取好友的用户信息
	friendUids := make([]string, 0, len(friends))
	friendRemarkMap := make(map[string]string) // 好友备注映射
	for _, friend := range friends {
		friendUids = append(friendUids, friend.FriendUid)
		remark := ""
		if friend.Remark.Valid {
			remark = friend.Remark.String
		}
		friendRemarkMap[friend.FriendUid] = remark
	}

	// 调用用户服务获取用户详细信息
	userResp, err := l.svcCtx.UserRpc.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: friendUids,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get user info failed: %v", err)
	}

	// 根据关键词筛选好友
	keyword := strings.ToLower(strings.TrimSpace(in.Keyword))
	var resultList []*social.SearchFriendInfo

	for _, user := range userResp.User {
		// 检查是否匹配搜索条件（昵称、手机号、用户ID）
		if l.matchSearchCondition(user, keyword) {
			result := &social.SearchFriendInfo{
				Id:       user.Id,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Phone:    user.Phone,
				Remark:   friendRemarkMap[user.Id], // 添加好友备注
			}
			resultList = append(resultList, result)
		}
	}

	return &social.FriendSearchResp{List: resultList}, nil
}

// matchSearchCondition 检查用户是否匹配搜索条件
func (l *FriendSearchLogic) matchSearchCondition(user *userclient.UserEntity, keyword string) bool {
	// 转换为小写进行比较
	nickname := strings.ToLower(user.Nickname)
	phone := strings.ToLower(user.Phone)
	userId := strings.ToLower(user.Id)

	// 检查昵称、手机号、用户ID是否包含关键词
	return strings.Contains(nickname, keyword) ||
		strings.Contains(phone, keyword) ||
		strings.Contains(userId, keyword)
}

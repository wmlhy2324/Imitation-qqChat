package logic

import (
	"context"
	"strings"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendAndGroupSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendAndGroupSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendAndGroupSearchLogic {
	return &FriendAndGroupSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索好友和群聊
func (l *FriendAndGroupSearchLogic) FriendAndGroupSearch(in *social.FriendAndGroupSearchReq) (*social.FriendAndGroupSearchResp, error) {
	if in.Keyword == "" {
		return &social.FriendAndGroupSearchResp{
			FriendList: []*social.SearchFriendInfo{},
			GroupList:  []*social.SearchGroupInfo{},
		}, nil
	}

	// 并发搜索好友和群聊
	friendsChan := make(chan []*social.SearchFriendInfo, 1)
	groupsChan := make(chan []*social.SearchGroupInfo, 1)
	errChan := make(chan error, 2)

	// 搜索好友
	go func() {
		friends, err := l.searchFriends(in.UserId, in.Keyword)
		if err != nil {
			errChan <- err
			return
		}
		friendsChan <- friends
	}()

	// 搜索群聊
	go func() {
		groups, err := l.searchGroups(in.UserId, in.Keyword)
		if err != nil {
			errChan <- err
			return
		}
		groupsChan <- groups
	}()

	// 收集结果
	var friends []*social.SearchFriendInfo
	var groups []*social.SearchGroupInfo

	for i := 0; i < 2; i++ {
		select {
		case friends = <-friendsChan:
		case groups = <-groupsChan:
		case err := <-errChan:
			l.Errorf("搜索失败: %v", err)
			// 即使一个搜索失败，也继续返回另一个的结果
		}
	}

	return &social.FriendAndGroupSearchResp{
		FriendList: friends,
		GroupList:  groups,
	}, nil
}

// searchFriends 搜索好友
func (l *FriendAndGroupSearchLogic) searchFriends(userId, keyword string) ([]*social.SearchFriendInfo, error) {
	// 获取用户的好友列表
	friends, err := l.svcCtx.FriendsModel.SearchFriends(l.ctx, userId, keyword)
	if err != nil {
		return nil, err
	}

	if len(friends) == 0 {
		return []*social.SearchFriendInfo{}, nil
	}

	// 获取好友的用户信息用于搜索匹配
	friendIds := make([]string, 0, len(friends))
	for _, friend := range friends {
		friendIds = append(friendIds, friend.FriendUid)
	}

	// 调用用户服务获取用户信息
	userResp, err := l.svcCtx.UserRpc.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: friendIds,
	})
	if err != nil {
		l.Errorf("获取用户信息失败: %v", err)
		return []*social.SearchFriendInfo{}, nil
	}

	// 创建用户信息映射
	userMap := make(map[string]*userclient.UserEntity)
	for _, user := range userResp.User {
		userMap[user.Id] = user
	}

	// 根据关键词过滤并构建结果
	result := make([]*social.SearchFriendInfo, 0)
	for _, friend := range friends {
		user, exists := userMap[friend.FriendUid]
		if !exists {
			continue
		}

		// 检查是否匹配关键词（昵称或手机号）
		if l.matchKeyword(user.Nickname, user.Phone, keyword) {
			remark := ""
			if friend.Remark.Valid {
				remark = friend.Remark.String
			}

			result = append(result, &social.SearchFriendInfo{
				Id:       user.Id,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Phone:    user.Phone,
				Remark:   remark,
			})
		}
	}

	return result, nil
}

// searchGroups 搜索群聊
func (l *FriendAndGroupSearchLogic) searchGroups(userId, keyword string) ([]*social.SearchGroupInfo, error) {
	// 获取用户加入的群列表
	userGroups, err := l.svcCtx.GroupsModel.ListByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("获取用户群列表失败: %v", err)
		userGroups = []*socialmodels.Groups{} // 如果获取用户群失败，继续搜索公开群
	}

	// 搜索匹配关键词的群（包括用户未加入的群）
	allMatchedGroups, err := l.svcCtx.GroupsModel.SearchGroups(l.ctx, keyword)
	if err != nil {
		return []*social.SearchGroupInfo{}, err
	}

	// 创建用户已加入群组的映射
	joinedGroupMap := make(map[string]bool)
	for _, group := range userGroups {
		joinedGroupMap[group.Id] = true
	}

	// 构建结果（包括已加入和未加入的群）
	result := make([]*social.SearchGroupInfo, 0)
	for _, group := range allMatchedGroups {
		// 获取群成员数量
		memberCount, err := l.getGroupMemberCount(group.Id)
		if err != nil {
			l.Errorf("获取群成员数量失败: %v", err)
			memberCount = 0
		}

		notification := ""
		if group.Notification.Valid {
			notification = group.Notification.String
		}

		result = append(result, &social.SearchGroupInfo{
			Id:           group.Id,
			Name:         group.Name,
			Icon:         group.Icon,
			MemberCount:  int32(memberCount),
			Notification: notification,
			IsJoined:     joinedGroupMap[group.Id],
		})
	}

	return result, nil
}

// matchKeyword 检查是否匹配关键词
func (l *FriendAndGroupSearchLogic) matchKeyword(nickname, phone, keyword string) bool {
	keyword = strings.ToLower(keyword)
	return strings.Contains(strings.ToLower(nickname), keyword) ||
		strings.Contains(strings.ToLower(phone), keyword)
}

// getGroupMemberCount 获取群成员数量
func (l *FriendAndGroupSearchLogic) getGroupMemberCount(groupId string) (int64, error) {
	// 这里需要查询group_members表来获取成员数量
	// 由于当前没有直接的方法，我们先返回0，后续可以优化
	return 0, nil
}

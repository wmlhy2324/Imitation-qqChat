package logic

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"imooc.com/easy-chat/pkg/constants"
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
// 包括会话ID、聊天类型、是否显示、消息统计等信息，以及用户头像、名字、群人数等前端渲染所需信息
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
		// 收集需要查询的用户ID和群组ID
		userIds := make([]string, 0)
		groupIds := make([]string, 0)

		for _, conversation := range rpcResp.ConversationList {
			if conversation.ChatType == int32(constants.GroupChatType) {
				// 群聊：需要获取群组信息
				if conversation.TargetId != "" {
					groupIds = append(groupIds, conversation.TargetId)
				}
			} else if conversation.ChatType == int32(constants.SingleChatType) {
				// 私聊：需要获取对方用户信息
				if conversation.TargetId != "" && conversation.TargetId != userId {
					userIds = append(userIds, conversation.TargetId)
				}
			}
		}

		// 批量获取用户信息
		userInfoMap := make(map[string]*userclient.UserEntity)
		if len(userIds) > 0 {
			userResp, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
				Ids: userIds,
			})
			if err != nil {
				l.Errorf("获取用户信息失败: %v", err)
			} else {
				for _, user := range userResp.User {
					userInfoMap[user.Id] = user
				}
			}
		}

		// 批量获取群组信息和成员数量
		groupInfoMap := make(map[string]*socialclient.Groups)
		groupMemberCountMap := make(map[string]int32)

		if len(groupIds) > 0 {
			// 一次性获取用户的所有群组信息
			groupResp, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
				UserId: userId,
			})
			if err != nil {
				l.Errorf("获取群组列表失败: %v", err)
			} else {
				// 建立群组ID到群组信息的映射
				for _, group := range groupResp.List {
					for _, groupId := range groupIds {
						if group.Id == groupId {
							groupInfoMap[groupId] = group
							break
						}
					}
				}
			}

			// 为每个群组获取成员数量
			for _, groupId := range groupIds {
				membersResp, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
					GroupId: groupId,
				})
				if err != nil {
					l.Errorf("获取群组[%s]成员失败: %v", groupId, err)
					// 即使获取失败，也设置为0，避免前端显示问题
					groupMemberCountMap[groupId] = 0
				} else {
					groupMemberCountMap[groupId] = int32(len(membersResp.List))
				}
			}
		}

		// 填充会话数据
		for conversationId, conversation := range rpcResp.ConversationList {
			apiConversation := &types.Conversation{}
			err = copier.Copy(apiConversation, conversation)
			if err != nil {
				l.Errorf("复制会话数据失败: %v", err)
				return nil, err
			}

			// 根据聊天类型填充额外信息
			if conversation.ChatType == int32(constants.GroupChatType) {
				// 群聊：填充群组信息
				if groupInfo, exists := groupInfoMap[conversation.TargetId]; exists {
					apiConversation.Name = groupInfo.Name
					apiConversation.Avatar = groupInfo.Icon
				}
				if memberCount, exists := groupMemberCountMap[conversation.TargetId]; exists {
					apiConversation.MemberCount = memberCount
				}
			} else if conversation.ChatType == int32(constants.SingleChatType) {
				// 私聊：填充对方用户信息
				if userInfo, exists := userInfoMap[conversation.TargetId]; exists {
					apiConversation.Name = userInfo.Nickname
					apiConversation.Avatar = userInfo.Avatar
				}
			}

			resp.ConversationList[conversationId] = apiConversation
		}
	}

	return resp, nil
}

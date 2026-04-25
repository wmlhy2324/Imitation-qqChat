package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/apps/task/mq/mq"
	"imooc.com/easy-chat/pkg/constants"
)

// 撤回时间窗口：2 分钟
const RevokeTimeWindowMs = 120000

type MsgRevokeTransfer struct {
	*baseMsgTransfer
}

func NewMsgRevokeTransfer(svc *svc.ServiceContext) *MsgRevokeTransfer {
	return &MsgRevokeTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgRevokeTransfer) Consume(key, value string) error {
	fmt.Println("MsgRevokeTransfer key:", key, "value:", value)

	var (
		data mq.MsgRevokeTransfer
		ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// Step 1: 查询 ChatLog
	chatLog, err := m.svcCtx.ChatLogModel.FindOne(ctx, data.MsgId)
	if err != nil {
		m.Errorf("MsgRevokeTransfer FindOne err %v, msgId %v", err, data.MsgId)
		return err
	}
	if chatLog.Status == 4 {
		// 已经撤回，幂等处理
		m.Infof("MsgRevokeTransfer msg already revoked, msgId %v", data.MsgId)
		return nil
	}

	// Step 2: 权限校验
	if !m.checkPermission(data.SendId, chatLog) {
		// 推送权限不足通知给操作者
		m.pushErrorToUser(data.SendId, data.ConversationId, chatLog.ChatType, "无权撤回他人消息")
		return nil
	}

	// Step 3: 时间校验
	if time.Now().UnixMilli()-chatLog.SendTime > RevokeTimeWindowMs {
		m.pushErrorToUser(data.SendId, data.ConversationId, chatLog.ChatType, "已超过2分钟无法撤回")
		return nil
	}

	// Step 4: 执行撤回
	chatLog.Status = 4 // 已撤回
	if err = m.svcCtx.ChatLogModel.Update(ctx, chatLog); err != nil {
		m.Errorf("MsgRevokeTransfer Update chatLog err %v", err)
		return err
	}

	// 如果撤回的是最后一条消息，更新会话摘要
	m.updateLastMsgIfNeeded(ctx, chatLog)

	// Step 5: 推送撤回通知给全体参与者
	push := &ws.Push{
		ConversationId: chatLog.ConversationId,
		ChatType:       chatLog.ChatType,
		SendId:         data.SendId,
		RecvId:         chatLog.RecvId,
		MsgId:          data.MsgId,
		ContentType:    constants.ContentRevoke,
		Content:        fmt.Sprintf("%s撤回了一条消息", data.SendId),
	}

	return m.Transfer(ctx, push)
}

// checkPermission 检查撤回权限
// 返回 true 表示允许撤回
func (m *MsgRevokeTransfer) checkPermission(userId string, chatLog *immodels.ChatLog) bool {
	// 发送者本人可以撤回
	if userId == chatLog.SendId {
		return true
	}

	// 群聊中，群主和管理员也可以撤回
	if chatLog.ChatType == constants.GroupChatType {
		users, err := m.svcCtx.Social.GroupUsers(context.Background(), &socialclient.GroupUsersReq{
			GroupId: chatLog.RecvId,
		})
		if err != nil {
			m.Errorf("checkPermission GroupUsers err %v", err)
			return false
		}
		for _, member := range users.List {
			if member.UserId == userId {
				roleLevel := constants.GroupRoleLevel(member.RoleLevel)
				if roleLevel == constants.CreatorGroupRoleLevel ||
					roleLevel == constants.ManagerGroupRoleLevel {
					return true
				}
			}
		}
	}

	return false
}

// pushErrorToUser 推送错误消息给特定用户
func (m *MsgRevokeTransfer) pushErrorToUser(userId, conversationId string, chatType constants.ChatType, errMsg string) {
	m.Errorf("MsgRevokeTransfer revoke denied: userId=%s, err=%s", userId, errMsg)
	// 只推送给操作者本人
	push := &ws.Push{
		ConversationId: conversationId,
		ChatType:       chatType,
		RecvId:         userId,
		ContentType:    constants.ContentRevoke,
		Content:        errMsg,
	}
	if err := m.Transfer(context.Background(), push); err != nil {
		m.Errorf("pushErrorToUser err %v", err)
	}
}

// updateLastMsgIfNeeded 如果撤回的是最后一条消息，更新会话摘要
func (m *MsgRevokeTransfer) updateLastMsgIfNeeded(ctx context.Context, chatLog *immodels.ChatLog) {
	conversation, err := m.svcCtx.ConversationModel.FindOne(ctx, chatLog.ConversationId)
	if err != nil {
		m.Errorf("updateLastMsgIfNeeded FindOne err %v", err)
		return
	}

	if conversation.Msg != nil && conversation.Msg.ID == chatLog.ID {
		// 撤回的是最后一条消息，更新内容
		chatLog.MsgContent = "消息已撤回"
		chatLog.Status = 4
		if err := m.svcCtx.ConversationModel.UpdateMsg(ctx, chatLog); err != nil {
			m.Errorf("updateLastMsgIfNeeded UpdateMsg err %v", err)
		}
	}
}

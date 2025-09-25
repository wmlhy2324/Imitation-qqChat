package immodels

import (
	"time"

	"imooc.com/easy-chat/pkg/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversation 会话结构体，用于存储聊天会话的基本信息
type Conversation struct {
	// MongoDB文档ID，作为主键
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// 会话唯一标识ID，用于区分不同的聊天会话
	// 对于私聊：由两个用户ID组合生成（如：user1_user2）
	// 对于群聊：直接使用群组ID
	ConversationId string `bson:"conversationId,omitempty"`

	// 聊天类型，1=私聊，2=群聊
	ChatType constants.ChatType `bson:"chatType,omitempty"`

	// 目标ID（已注释，可能用于存储对方用户ID或群组ID）
	TargetId string `bson:"targetId,omitempty"`

	// 是否在会话列表中显示此会话，true=显示，false=隐藏
	IsShow bool `bson:"isShow,omitempty"`

	// 该会话的消息总数量
	Total int `bson:"total,omitempty"`

	// 消息序列号，用于消息排序和去重，数字越大表示消息越新
	Seq int64 `bson:"seq"`

	// 会话中最新的一条消息，用于在会话列表中显示预览
	Msg *ChatLog `bson:"msg,omitempty"`

	// 会话最后更新时间，每当有新消息时更新
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`

	// 会话创建时间，首次建立会话时设置
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

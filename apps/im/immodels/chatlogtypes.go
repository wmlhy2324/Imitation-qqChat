package immodels

import (
	"time"

	"imooc.com/easy-chat/pkg/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DefaultChatLogLimit int64 = 100

// ChatLog 聊天记录结构体，存储所有聊天消息的详细信息
type ChatLog struct {
	// ID MongoDB文档的唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// ConversationId 会话ID，用于标识这条消息属于哪个会话
	// 单聊：由两个用户ID组成的唯一标识
	// 群聊：群组的唯一标识符
	ConversationId string `bson:"conversationId"`

	// SendId 发送者的用户ID
	SendId string `bson:"sendId"`

	// RecvId 接收者的用户ID
	// 单聊：接收方用户ID
	// 群聊：可能为空或群组ID
	RecvId string `bson:"recvId"`

	// MsgFrom 消息来源标识
	// 用于区分消息是从哪个客户端或渠道发送的
	MsgFrom int `bson:"msgFrom"`

	// ChatType 聊天类型
	// 1: GroupChatType-群聊, 2: SingleChatType-单聊
	ChatType constants.ChatType `bson:"chatType"`

	// MsgType 消息类型
	// 0: TextMType-文本消息（当前只支持文本消息）
	MsgType constants.MType `bson:"msgType"`

	// MsgContent 消息内容
	// 当前支持文本消息，直接存储文本内容
	// 未来扩展：媒体消息可存储文件路径或URL
	MsgContent string `bson:"msgContent"`

	// SendTime 消息发送时间戳（Unix时间戳，毫秒）
	SendTime int64 `bson:"sendTime"`

	// Status 消息状态
	// 如：0-发送中, 1-已发送, 2-已送达, 3-已读, -1-发送失败
	Status int `bson:"status"`

	// ReadRecords 消息已读记录
	// 存储哪些用户已读此消息的信息（序列化后的字节数组）
	// 主要用于群聊中追踪消息的阅读状态
	ReadRecords []byte `bson:"readRecords"`

	// UpdateAt 记录最后更新时间
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`

	// CreateAt 记录创建时间
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

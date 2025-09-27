/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package websocket

import "time"

type FrameType uint8

const (
	FrameData      FrameType = 0x0
	FramePing      FrameType = 0x1
	FrameAck       FrameType = 0x2
	FrameNoAck     FrameType = 0x3
	FrameErr       FrameType = 0x9
	FrameTranspond FrameType = 0x6

	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// Message WebSocket消息结构体，定义了客户端与服务器之间通信的消息格式
type Message struct {
	// FrameType 消息帧类型
	// 0x0: FrameData-数据帧, 0x1: FramePing-心跳帧, 0x2: FrameAck-确认帧
	// 0x3: FrameNoAck-无需确认帧, 0x6: FrameTranspond-转发帧, 0x9: FrameErr-错误帧
	FrameType `json:"frameType"`

	// Id 消息唯一标识符，用于消息去重和确认机制
	Id string `json:"id"`

	// TranspondUid 转发目标用户ID，用于分布式环境下的消息转发
	// 当消息需要转发到其他服务器上的用户时使用
	TranspondUid string `json:"transpondUid"`

	// AckSeq 确认序列号，用于消息确认机制
	// 在RigorAck模式下，客户端需要回复对应的序列号来确认消息已收到
	AckSeq int `json:"ackSeq"`

	// ackTime 确认时间戳，记录消息确认的时间（内部使用，不序列化给客户端）
	ackTime time.Time `json:"ackTime"`

	// errCount 错误计数，记录消息处理失败的次数（内部使用，不序列化给客户端）
	errCount int `json:"errCount"`

	// Method 请求方法名，用于路由到对应的处理器
	// 如: "user.online", "conversation.chat", "conversation.markChat", "push"
	Method string `json:"method"`

	// FormId 发送方用户ID，标识消息的发起者
	// 对于系统消息，可能使用特殊的系统ID
	FormId string `json:"formId"`

	// Data 消息数据载荷，可以是任意JSON结构
	// 根据Method的不同，包含不同的业务数据
	Data interface{} `json:"data"` // map[string]interface{}
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}

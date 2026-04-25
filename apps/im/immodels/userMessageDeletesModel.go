package immodels

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserMessageDelete 用户删除消息的记录
// 用于记录用户删除了哪些消息（仅本端隐藏，不影响其他用户）
type UserMessageDelete struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserId         string             `bson:"userId"`
	ConversationId string             `bson:"conversationId"`
	MsgId          string             `bson:"msgId"` // ChatLog._id 的 hex
	CreateAt       time.Time          `bson:"createAt"`
}

var _ UserMessageDeletesModel = (*defaultUserMessageDeletesModel)(nil)

type UserMessageDeletesModel interface {
	Insert(ctx context.Context, data *UserMessageDelete) error
	ListByUserIdAndConversation(ctx context.Context, userId, conversationId string) ([]*UserMessageDelete, error)
	DeleteByConversation(ctx context.Context, userId, conversationId string) error
}

type defaultUserMessageDeletesModel struct {
	conn *mon.Model
}

func NewUserMessageDeletesModel(url, db string) UserMessageDeletesModel {
	conn := mon.MustNewModel(url, db, "user_message_deletes")
	return &defaultUserMessageDeletesModel{conn: conn}
}

func MustUserMessageDeletesModel(url, db string) UserMessageDeletesModel {
	return NewUserMessageDeletesModel(url, db)
}

func (m *defaultUserMessageDeletesModel) Insert(ctx context.Context, data *UserMessageDelete) error {
	data.ID = primitive.NewObjectID()
	data.CreateAt = time.Now()
	_, err := m.conn.InsertOne(ctx, data)
	return err
}

func (m *defaultUserMessageDeletesModel) ListByUserIdAndConversation(ctx context.Context, userId, conversationId string) ([]*UserMessageDelete, error) {
	var data []*UserMessageDelete
	filter := bson.M{
		"userId":         userId,
		"conversationId": conversationId,
	}
	err := m.conn.Find(ctx, &data, filter)
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserMessageDeletesModel) DeleteByConversation(ctx context.Context, userId, conversationId string) error {
	filter := bson.M{
		"userId":         userId,
		"conversationId": conversationId,
	}
	_, err := m.conn.DeleteMany(ctx, filter)
	return err
}

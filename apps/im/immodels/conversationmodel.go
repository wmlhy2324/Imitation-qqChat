package immodels

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		FindOneByConversationIdAndTargetId(ctx context.Context, conversationId, targetId string) (*Conversation, error)
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the mongo.
func NewConversationModel(url, db, collection string) ConversationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

func MustConversationModel(url, db string) ConversationModel {
	return NewConversationModel(url, db, "conversation")
}

// FindOneByConversationIdAndTargetId 根据 conversationId 和 targetId 查询会话
func (m *customConversationModel) FindOneByConversationIdAndTargetId(ctx context.Context, conversationId, targetId string) (*Conversation, error) {
	var data Conversation

	err := m.conn.FindOne(ctx, &data, bson.M{
		"conversationId": conversationId,
		"targetId":       targetId,
	})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

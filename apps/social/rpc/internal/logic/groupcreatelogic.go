package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/wuid"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群要求
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	groups := &socialmodels.Groups{
		Id:          wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:        in.Name,
		Icon:        in.Icon,
		CreatorUid:  in.CreatorUid,
		Description: in.Description,
		IsVerify:    false,
	}

	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 创建群
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		// 添加群主
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int(constants.CreatorGroupRoleLevel),
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group creator member err %v req %v", err, in)
		}

		// 添加初始成员
		for _, memberId := range in.MemberIds {
			// 跳过群主，避免重复添加
			if memberId == in.CreatorUid {
				continue
			}
			_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
				GroupId:   groups.Id,
				UserId:    memberId,
				RoleLevel: int(constants.AtLargeGroupRoleLevel), // 普通成员
			})
			if err != nil {
				// 记录日志但不中断事务，允许部分成员添加失败
				logx.Errorf("insert group member failed, groupId: %s, userId: %s, err: %v", groups.Id, memberId, err)
			}
		}
		return nil
	})

	return &social.GroupCreateResp{
		Id: groups.Id,
	}, err
}

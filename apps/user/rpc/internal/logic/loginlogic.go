package logic

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/pkg/ctxdata"
	"imooc.com/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"imooc.com/easy-chat/apps/user/models"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/encrypt"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 添加调试日志 - 开始RPC登录流程
	l.Infof("=== 开始RPC登录流程 ===")
	l.Infof("输入参数: Phone=%s, Password=%s", in.Phone, maskPassword(in.Password))

	startTime := time.Now()

	// 1. 验证用户是否注册，根据手机号码验证
	l.Infof("步骤1: 正在根据手机号查询用户信息...")
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			l.Errorf("用户不存在: 手机号 %s 未注册", in.Phone)
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		l.Errorf("数据库查询失败: %v", err)
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	l.Infof("用户查询成功: ID=%s, Nickname=%s, Phone=%s",
		userEntity.Id, userEntity.Nickname, userEntity.Phone)

	// 2. 密码验证
	l.Infof("步骤2: 正在进行密码验证...")
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		l.Errorf("密码验证失败: 用户 %s 密码不正确", userEntity.Id)
		return nil, errors.WithStack(ErrUserPwdError)
	}
	l.Infof("密码验证成功")

	// 3. 生成token
	l.Infof("步骤3: 正在生成JWT Token...")
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id)
	if err != nil {
		l.Errorf("Token生成失败: %v", err)
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	l.Infof("Token生成成功: %s", maskToken(token))
	l.Infof("Token配置: Secret=%s, Expire=%d",
		maskSecret(l.svcCtx.Config.Jwt.AccessSecret), l.svcCtx.Config.Jwt.AccessExpire)

	// 4. 数据转换
	l.Infof("步骤4: 正在进行数据转换...")
	var u user.UserEntity
	if err := copier.Copy(&u, userEntity); err != nil {
		l.Errorf("数据转换失败: %v", err)
		return nil, err
	}
	l.Infof("数据转换成功")

	// 5. 构建响应
	expireTime := now + l.svcCtx.Config.Jwt.AccessExpire
	resp := &user.LoginResp{
		User:   &u,
		Id:     userEntity.Id,
		Token:  token,
		Expire: expireTime,
	}

	totalDuration := time.Since(startTime)
	l.Infof("=== RPC登录流程完成，总耗时: %v ===", totalDuration)
	l.Infof("返回结果: UserID=%s, Token=%s, Expire=%d",
		resp.Id, maskToken(resp.Token), resp.Expire)

	return resp, nil
}

// 辅助函数：掩码密码
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

// 辅助函数：掩码Token
func maskToken(token string) string {
	if len(token) <= 10 {
		return "****"
	}
	return token[:5] + "****" + token[len(token)-5:]
}

// 辅助函数：掩码Secret
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

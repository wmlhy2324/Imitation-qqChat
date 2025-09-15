package user

import (
	"context"

	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/constants"

	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/apps/user/api/internal/types"

	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 添加调试日志 - 开始登录流程
	l.Infof("=== 开始登录流程 ===")
	l.Infof("请求参数: Phone=%s, Password=%s", req.Phone, maskPassword(req.Password))

	// 记录请求时间
	startTime := time.Now()

	// 调用RPC服务进行登录验证
	l.Infof("正在调用user-rpc服务进行登录验证...")
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})

	// 记录RPC调用耗时
	rpcDuration := time.Since(startTime)
	l.Infof("RPC调用耗时: %v", rpcDuration)

	if err != nil {
		l.Errorf("RPC登录失败: %v", err)
		return nil, err
	}

	l.Infof("RPC登录成功，用户ID: %s", loginResp.Id)
	l.Infof("Token: %s", maskToken(loginResp.Token))
	l.Infof("过期时间: %d", loginResp.Expire)

	// 数据转换
	l.Infof("正在进行数据转换...")
	var res types.LoginResp
	if err := copier.Copy(&res, loginResp); err != nil {
		l.Errorf("数据转换失败: %v", err)
		return nil, err
	}
	l.Infof("数据转换成功")

	// 处理登入的业务 - 记录用户在线状态到Redis
	l.Infof("正在记录用户在线状态到Redis...")
	redisKey := constants.REDIS_ONLINE_USER
	redisErr := l.svcCtx.Redis.HsetCtx(l.ctx, redisKey, loginResp.Id, "1")
	if redisErr != nil {
		l.Errorf("Redis操作失败: %v", redisErr)
		// 注意：这里可以选择继续执行，不因为Redis失败而中断登录流程
	} else {
		l.Infof("Redis在线状态记录成功，Key: %s, UserID: %s", redisKey, loginResp.Id)
	}

	// 记录总耗时
	totalDuration := time.Since(startTime)
	l.Infof("=== 登录流程完成，总耗时: %v ===", totalDuration)

	return &res, nil
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

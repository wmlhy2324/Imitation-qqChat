/**
 * @Author: dn-jinmin
 * @File:  websocket
 * @Version: 1.0.0
 * @Date: 2024/3/27
 * @Description:
 */

package websocket

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
)

// 服务发现机制【该方式是去中心化，自己在内部实现服务发现整套机制】
// 该机制主要针对用户而设立
//
//	用户连接后，会将用户信息与服务器ip一起绑定注册到某一个位置比如redis
//	当用户发送信息的时候，根据发送的目标从记录位置中获取绑定关系，查找相应服务并发送
type Discover interface {
	// 注册服务
	Register(serverAddr string) error
	// 绑定用户
	BoundUser(uid string) error
	// 解除与用户绑定
	RelieveUser(uid string) error
	// 转发
	Transpond(msg interface{}, uid ...string) error
}

// 默认的
type nopDiscover struct {
	serverAddr string
}

// 注册服务
func (d *nopDiscover) Register(serverAddr string) error { return nil }

// 绑定用户
func (d *nopDiscover) BoundUser(uid string) error { return nil }

func (d *nopDiscover) RelieveUser(uid string) error { return nil }

// 转发消息
func (d *nopDiscover) Transpond(msg interface{}, uid ...string) error { return nil }

// 默认的
type redisDiscover struct {
	serverAddr string
	auth       http.Header
	//含义：Redis 中用于存储服务列表的键名
	//命名规则通常是业务相关的固定字符串（如 "chat_service_nodes"）
	//服务注册时将此实例地址写入这个 Key
	srvKey string
	//含义：存储用户-服务绑定关系的 Redis 键名
	//格式如: "user_service_mapping" 或 "chat_service:user_bindings"
	//用途：维护用户 ID 与所在服务地址的映射关系
	//存储结构：通常使用 Redis Hash 或 Sorted Set
	boundUserKey string
	redis        *redis.Redis
	//key：其他服务实例的地址（如 "192.168.1.101:8080"）
	//value：指向该服务实例的客户端连接对象
	clients map[string]Client
}

func NewRedisDiscover(auth http.Header, srvKey string, redisCfg redis.RedisConf) *redisDiscover {
	return &redisDiscover{
		srvKey:       srvKey,
		boundUserKey: fmt.Sprintf("%s.%s", srvKey, "boundUserKey"),
		redis:        redis.MustNewRedis(redisCfg),
		clients:      make(map[string]Client),
		auth:         auth,
	}
}

// 注册服务
func (d *redisDiscover) Register(serverAddr string) (err error) {
	d.serverAddr = serverAddr

	// 服务列表：redis存储用set
	go d.redis.Set(d.srvKey, serverAddr)

	return
}

// 绑定用户
func (d *redisDiscover) BoundUser(uid string) (err error) {
	// 用户绑定
	exists, err := d.redis.Hexists(d.boundUserKey, uid)
	if err != nil {
		return err
	}
	if exists {
		// 存在绑定关系
		return nil
	}

	// 绑定
	return d.redis.Hset(d.boundUserKey, uid, d.serverAddr)
}

func (d *redisDiscover) RelieveUser(uid string) (err error) {
	_, err = d.redis.Hdel(d.boundUserKey, uid)
	return
}

// 转发消息
func (d *redisDiscover) Transpond(msg interface{}, uids ...string) (err error) {

	for _, uid := range uids {
		srvAddr, err := d.redis.Hget(d.boundUserKey, uid)
		if err != nil {
			return err
		}
		srvClient, ok := d.clients[srvAddr]
		if !ok {
			srvClient = d.createClient(srvAddr)
		}

		fmt.Println("redis transpand -》 ", srvAddr, " uid ", uid)

		if err := d.send(srvClient, msg, uid); err != nil {
			return err
		}
	}

	return
}

func (d *redisDiscover) send(srvClient Client, msg interface{}, uid string) error {
	return srvClient.Send(Message{
		FrameType:    FrameTranspond,
		TranspondUid: uid,
		Data:         msg,
	})
}

func (d *redisDiscover) createClient(srvAddr string) Client {
	return NewClient(srvAddr, WithClientHeader(d.auth))
}

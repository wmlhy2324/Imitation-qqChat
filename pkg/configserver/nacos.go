/**
 * @author: dn-jinmin/dn-jinmin
 * @doc: Nacos 配置中心适配器
 */

package configserver

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

type NacosConfig struct {
	Addr      string `toml:"addr"`      // Nacos 服务器地址，例如：127.0.0.1:8848
	Namespace string `toml:"namespace"` // 命名空间
	Group     string `toml:"group"`     // 配置分组，默认 DEFAULT_GROUP
	DataId    string `toml:"data_id"`   // 配置文件 DataId
	Username  string `toml:"username"`  // 用户名
	Password  string `toml:"password"`  // 密码
	LogLevel  string `toml:"log_level"` // 日志级别
}

type Nacos struct {
	configClient config_client.IConfigClient
	onChange     OnChange
	config       *NacosConfig
	content      string
}

func NewNacos(cfg *NacosConfig) *Nacos {
	// 设置默认值
	if cfg.Group == "" {
		cfg.Group = "DEFAULT_GROUP"
	}
	if cfg.Username == "" {
		cfg.Username = "nacos"
	}
	if cfg.Password == "" {
		cfg.Password = "nacos"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "warn"
	}

	return &Nacos{
		config: cfg,
	}
}

func (n *Nacos) Build() error {
	// 解析服务器地址
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(n.config.Addr, 8848),
	}

	// 创建客户端配置
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(n.config.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel(n.config.LogLevel),
		constant.WithUsername(n.config.Username),
		constant.WithPassword(n.config.Password),
	)

	// 创建配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return fmt.Errorf("创建 Nacos 配置客户端失败: %w", err)
	}

	n.configClient = configClient

	// 监听配置变化
	if n.onChange != nil {
		err = n.configClient.ListenConfig(vo.ConfigParam{
			DataId: n.config.DataId,
			Group:  n.config.Group,
			OnChange: func(namespace, group, dataId, data string) {
				fmt.Printf("配置发生变化 - Namespace: %s, Group: %s, DataId: %s\n", namespace, group, dataId)
				n.content = data
				if n.onChange != nil {
					// 将 YAML 转换为 JSON 格式
					jsonData, err := n.yamlToJson([]byte(data))
					if err != nil {
						fmt.Printf("配置转换失败: %v\n", err)
						return
					}
					if err := n.onChange(jsonData); err != nil {
						fmt.Printf("配置更新回调失败: %v\n", err)
					}
				}
			},
		})
		if err != nil {
			return fmt.Errorf("监听配置变化失败: %w", err)
		}
	}

	return nil
}

func (n *Nacos) SetOnChange(f OnChange) {
	n.onChange = f
}

func (n *Nacos) FromJsonBytes() ([]byte, error) {
	// 获取配置内容
	content, err := n.configClient.GetConfig(vo.ConfigParam{
		DataId: n.config.DataId,
		Group:  n.config.Group,
	})
	if err != nil {
		return nil, fmt.Errorf("获取 Nacos 配置失败: %w", err)
	}

	n.content = content

	// 将 YAML 内容转换为 JSON
	return n.yamlToJson([]byte(content))
}

// yamlToJson 将 YAML 格式转换为 JSON 格式
func (n *Nacos) yamlToJson(yamlData []byte) ([]byte, error) {
	var data interface{}

	// 解析 YAML
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}

	// 转换 map[interface{}]interface{} 为 map[string]interface{}
	data = convertMapInterface(data)

	// 转换为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("转换为 JSON 失败: %w", err)
	}

	return jsonData, nil
}

// convertMapInterface 递归转换 map[interface{}]interface{} 为 map[string]interface{}
func convertMapInterface(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[fmt.Sprint(k)] = convertMapInterface(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertMapInterface(v)
		}
	}
	return i
}

// GetContent 获取原始配置内容
func (n *Nacos) GetContent() string {
	return n.content
}

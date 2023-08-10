package driver

import (
	"fmt"
	"github.com/polarismesh/polaris-go/pkg/config"
	"os"
	"strconv"
)

// 北极星配置
type PolarisConf struct {
	Namespace    string //命名空间
	Host         string //服务监听host
	ProviderPort int    //服务实例监听port-服务注册
	ServiceToken string //服务访问令牌
}

// 初始化北极星配置
func GetPolarisConfiguration(cfg PolarisConf) config.Configuration {
	if len(cfg.Host) < 1 {
		return nil
	}
	if cfg.ProviderPort < 1 {
		return nil
	}
	providerConf := config.NewDefaultConfigurationWithDomain()
	if cfg.ProviderPort > 0 {
		providerConf.GetGlobal().GetServerConnector().SetAddresses([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.ProviderPort)})
	}
	return providerConf
}

// 获取默认北极星配置
func GetPolarisConf() PolarisConf {
	//从环境变量获取
	return PolarisConf{
		Namespace:    GetEnv("DTM-MICRO-SERVICE-POLARIS-NAMESPACE", "dtm"),
		Host:         GetEnv("DTM-MICRO-SERVICE-POLARIS-HOST", "127.0.0.1"),
		ProviderPort: GetEnvInt("DTM-MICRO-SERVICE-POLARIS-PROVIDERPORT", "8091"),
		ServiceToken: GetEnv("DTM-MICRO-SERVICE-POLARIS-SERVICETOKEN", ""),
	}
}

// 获取环境变量
func GetEnv(name string, defaultValue string) string {
	ret := os.Getenv(name)
	if ret == "" {
		ret = defaultValue
	}
	return ret
}

// 默认int型的环境变量
func GetEnvInt(name string, defaultValue string) int {
	ret := GetEnv(name, defaultValue)
	if len(ret) > 0 {
		i, _ := strconv.Atoi(ret)
		return i
	}
	return 0
}

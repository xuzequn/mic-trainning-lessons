package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"mic-trainning-lessons/internal"
)

func main() {
	nacosconfig := internal.ViperConf.NacosConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosconfig.Host,
			Port:   nacosconfig.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosconfig.NameSpace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosconfig.DataId,
		Group:  nacosconfig.Group,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}

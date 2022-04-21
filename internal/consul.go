package internal

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int32  `mapstructure:"port"`
}

func Reg(host, name, id string, port int, tags []string) error {
	defaultConfig := api.DefaultConfig()
	h := ViperConf.ConsulConfig.Host
	p := ViperConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}

	agentServiceRegistration := new(api.AgentServiceRegistration)
	agentServiceRegistration.Address = defaultConfig.Address
	agentServiceRegistration.Port = port
	agentServiceRegistration.ID = id
	agentServiceRegistration.Name = name
	agentServiceRegistration.Tags = tags
	serverAddr := fmt.Sprintf("http://%s:%d/health", host, port)
	check := api.AgentServiceCheck{
		HTTP:                           serverAddr,
		Timeout:                        "3S",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5S",
	}
	agentServiceRegistration.Check = &check
	return client.Agent().ServiceRegister(agentServiceRegistration)

}

func GetServiceList() error {
	defaultConfig := api.DefaultConfig()
	h := ViperConf.ConsulConfig.Host
	p := ViperConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}

	services, err := client.Agent().Services()
	if err != nil {
		return err
	}
	for k, v := range services {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("----------------------")
	}
	return nil
}

func FilterService() error {
	defaultConfig := api.DefaultConfig()
	h := ViperConf.ConsulConfig.Host
	p := ViperConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}

	services, err := client.Agent().ServicesWithFilter("Service==account_web")
	if err != nil {
		return err
	}
	for k, v := range services {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("----------------------")
	}
	return nil
}
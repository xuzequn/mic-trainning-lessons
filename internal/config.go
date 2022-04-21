package internal

type AccountSrvConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srv_name"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

type AccountWebConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srv_name"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

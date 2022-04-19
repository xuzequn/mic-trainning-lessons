package conf

type JWTConfig struct {
	SingingKey string `mapstructure:"signing_key"`
}

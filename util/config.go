package util

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Username  string `envconfig:"username",default:"nikita",required:"true"`
	Password  string `envconfig:"password",default:"Nikita123!",required:"true"`
	DataPath  string `envconfig:"datapath",default:"data",required:"true"`
	JwtSecret string `envconfig:"jwtsecret",required:"true"`
}

var AppConfig *Config

func InitConfig() (err error) {
	AppConfig = &Config{}
	err = envconfig.Process("hd", AppConfig)
	if err != nil {
		return
	}
	return nil
}

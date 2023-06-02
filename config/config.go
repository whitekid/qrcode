package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/whitekid/goxp/flags"
)

const (
	keyBind      = "bind_addr"
	keyRateLimit = "rate_limit"

	keyGrpcBind = "bind_addr"
)

var configs = map[string][]flags.Flag{
	"qrcodeapi": {
		{keyBind, "B", "127.0.0.1:8000", "bind address"},
		{keyRateLimit, "", "20", "rate limit"},
	},
	"grpc-server": {
		{keyGrpcBind, "B", "127.0.0.1:9000", "bind address"},
	},
}

func init() {
	viper.SetEnvPrefix("qr")
	viper.AutomaticEnv()

	flags.InitDefaults(nil, configs)
}

func InitFlagSet(use string, fs *pflag.FlagSet) { flags.InitFlagSet(nil, configs, use, fs) }

func BindAddr() string { return viper.GetString(keyBind) }
func RateLimit() int   { return viper.GetInt(keyRateLimit) }

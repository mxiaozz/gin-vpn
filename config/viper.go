package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"openvpn.funcworks.net/log"
)

var Viper = viper.New()

func init() {
	Viper.AutomaticEnv()
	Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	path, err := os.Executable()
	if err == nil {
		Viper.AddConfigPath(filepath.Dir(path))
	}
	Viper.AddConfigPath("/etc/openvpn")
	Viper.AddConfigPath(".")
	Viper.SetConfigName("app")
	Viper.SetConfigType("yaml")

	err = Viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}
}

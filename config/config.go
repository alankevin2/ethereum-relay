package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

const Debug = true

func InitializeConfiguration() {
	// viper.AddConfigPath("config")
	// viper.SetConfigType("yaml")
}

type ProviderInfo struct {
	URL         string
	IsUsingGeth bool
	IsUsingWS   bool
}

func GetProviderInfo(fileName string) *ProviderInfo {

	currentPath, _ := os.Getwd()
	// viper.AddConfigPath("gitlab.inlive7.com/crypto/ethereum-relay/config/")
	// viper.AddConfigPath(currentPath)
	fullpath := path.Join(currentPath, "config", fileName)
	_, err := os.Stat(fullpath)
	if err != nil {
		fullpath = path.Join(currentPath, "vendor/gitlab.inlive7.com/crypto/ethereum-relay/config/", fileName)
	}
	viper.SetConfigFile(fullpath)
	viper.SetConfigType("yml")
	// viper.SetConfigName(fileName)
	// viper.SetConfigType("yml")
	// viper.SetConfigFile(fileName)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return &ProviderInfo{
		viper.GetString("root.url"),
		viper.GetBool("root.is-using-geth"),
		viper.GetBool("root.is-using-ws"),
	}
}

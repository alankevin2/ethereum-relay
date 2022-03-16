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

	currentPath, err := os.Getwd()
	fmt.Println(currentPath)
	// viper.AddConfigPath("gitlab.inlive7.com/crypto/ethereum-relay/config/")
	// viper.AddConfigPath(currentPath)
	viper.SetConfigFile(path.Join(currentPath, "config", fileName))
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

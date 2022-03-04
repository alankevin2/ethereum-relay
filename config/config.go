package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const Debug = true

func InitializeConfiguration() {
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
}

type ProviderInfo struct {
	URL         string
	IsUsingGeth bool
	IsUsingWS   bool
}

func GetProviderInfo(fileName string) *ProviderInfo {
	viper.SetConfigName(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return &ProviderInfo{
		viper.GetString("root.url"),
		viper.GetBool("root.is-using-geth"),
		viper.GetBool("root.is-using-ws"),
	}
}

package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"

	"github.com/spf13/viper"
)

const Debug = true

type ProviderInfo struct {
	URL         string
	IsUsingGeth bool
	IsUsingWS   bool
	Tokens      map[string]string
}

func GetProviderInfo(fileName string) *ProviderInfo {
	currentPath, _ := os.Getwd()
	fullpath := path.Join(currentPath, "config", fileName)
	_, err := os.Stat(fullpath)
	// 如果找不到，代表當前執行環境不是以此pkg為主，而是被別人vendor引用
	if err != nil {
		pkgPath := reflect.TypeOf(ProviderInfo{}).PkgPath()
		fullpath = path.Join(currentPath, "vendor", pkgPath, fileName)
	}
	viper.SetConfigFile(fullpath)
	viper.SetConfigType("yml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return &ProviderInfo{
		viper.GetString("root.url"),
		viper.GetBool("root.is-using-geth"),
		viper.GetBool("root.is-using-ws"),
		viper.GetStringMapString("root.tokens"),
	}
}

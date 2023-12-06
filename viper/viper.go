package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

var Config = struct {
	Server ServerConfig `mapstructure:"server"`
}{}

func main() {
	viper.SetConfigFile("./conf/config.ini")
	viper.SetConfigType("ini")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(viper.GetString("server.name"))
	fmt.Println(viper.GetString("server.host"))
	fmt.Println(viper.GetInt("server.port"))

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err = viper.Unmarshal(&Config)
		if err != nil {
			fmt.Printf("unmarshal on file changed, err: %v\n", err)
			return
		}
		fmt.Printf("file changed: %+v\n", Config)
	})

	// 保存到结构体中
	err = viper.Unmarshal(&Config)
	if err != nil {
		fmt.Printf("unmarshal err: %v\n", err)
		return
	}
	fmt.Printf("Conifg: %+v\n", Config)

	select {}
}

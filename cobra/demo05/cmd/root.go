package cmd

import (
	"fmt"
	"github.com/lgc202/go-example/cobra/demo05/app/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile   string            // 配置文件
	appConfig *config.AppConfig // 配置对应的结构体

	rootCmd = &cobra.Command{
		Use:   "",
		Short: "命令行的简要描述....",
		Long:  `学习使用Cobra,开发cli项目,app: 指的是编译后的文件名。`,
	}
)

func initConfig() {
	// 接收指定的配置文件
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// 设置配置文件目录(可以设置多个,优先级根据添加顺序来)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./app/config")
		// 设置配置文件
		viper.SetConfigType("yaml")
		viper.SetConfigName("app")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig: %v\n", err)
	}

	// 解析配置信息
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Printf("%+v\n", appConfig)
}

func init() {
	// 初始化配置信息
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./app.yaml | ./config/app.yaml )")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

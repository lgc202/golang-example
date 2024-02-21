package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

var (
	serverCmd = &cobra.Command{
		Use:   "service",
		Short: "启动http服务,使用方法: app service?",
		Run: func(cmd *cobra.Command, args []string) {
			// 使用配置
			if appConfig.App.Port == "" {
				fmt.Println("port不能为空!")
				os.Exit(-1)
			}
			engine := gin.Default()
			_ = engine.Run(":" + appConfig.App.Port)
		},
	}
)

func init() {
	// 添加命令
	rootCmd.AddCommand(serverCmd)
}

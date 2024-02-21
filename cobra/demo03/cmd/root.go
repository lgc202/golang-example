package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "命令行的简要描述",
	Long: `使用cobra开发cli命令,
-app: 指的是编译后的文件名`,
	// 根命令执行方法，需要就添加
	// Run: func(cmd *cobra.Command, args []string) {
	//
	// },
}

func init() {
	rootCmd.PersistentFlags().String("version", "", "版本")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

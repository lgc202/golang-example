package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"unicode/utf8"
)

var (
	name string
	list []string

	// userCmd 父命令
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "用户操作",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("用户列表: ", list)
		},
	}

	// 添加用户子命令
	addUserCmd = &cobra.Command{
		Use:   "add",
		Short: "添加用户: user add --name=?",
		// 非选项参数校验方式1：用内置函数；RangeArgs表示传入位置参数个数 min<= N <= max，否则报错
		Args: cobra.RangeArgs(1, 3),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("非选项参数(args):", args)
		},
	}

	// 删除用户子命令
	delUserCmd = &cobra.Command{
		Use:   "del",
		Short: "删除用户: user del --name=?",
		// 非选项参数校验方式2：自定义参数限制
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("参数数量不对")
			}
			// 判断姓名长度
			count := utf8.RuneCountInString(args[0])
			fmt.Printf("%v %v \n", args[0], count)
			if count > 4 {
				return errors.New("姓名长度过长")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("非选项参数(args):", args)
		},
	}
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(delUserCmd)

	addUserCmd.Flags().StringVar(&name, "name", "", "用户名")
	// 选项参数校验
	err := addUserCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println("--name 不能为空")
		return
	}
}

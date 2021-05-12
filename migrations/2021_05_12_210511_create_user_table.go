package main

import (
	"fmt"
	"github.com/Ravior/go-fast-migrate/schema"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "up" {
			up()
		} else if cmd == "down" {
			down()
		}
	}
}

// UP 迁移操作
func up()  {
	fmt.Println("CreateUserTable2021_05_12_210511 Up")
	// Write your migrate action here
	// 创建用户表
	schema.NewSchema().CreateTable("user", func(builder *schema.Builder) {
		builder.BigIncrements("id")
		builder.String("open_id", 32).Unique().Default("").Comment("平台用户标识")
		builder.String("nickname", 64).Default("").Comment("用户在平台昵称")
		builder.String("avatar").Default("").Comment("用户在平台头像")
		builder.UnsignedInteger("last_login_time").Default(0).Comment("最新登录时间")
		builder.UnsignedInteger("login_times").Default(0).Comment("登陆次数")
		builder.UnsignedInteger("continue_login_days").Default(0).Comment("连续登陆天数")
		builder.UnsignedInteger("date").Default(0).Comment("注册日期")
		builder.String("qq", 20).Unique().Default("").Comment("用户QQ号")
		builder.UnsignedTinyInteger("status").Default(0).Comment("用户状态，0:正常;1:可疑；2:封禁")

		builder.Timestamps()
	})
}

// Down 迁移回滚
func down()  {
	fmt.Println("CreateUserTable2021_05_12_210511 Down")
	// Write your rollback action here
	// 创建用户表
	schema.NewSchema().DropTable("user")
}


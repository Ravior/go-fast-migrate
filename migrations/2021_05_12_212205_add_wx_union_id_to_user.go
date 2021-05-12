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
	fmt.Println("AddWxUnionIdToUser2021_05_12_212205 Up")
	schema.NewSchema().AlterTable("user", func(builder *schema.Builder) {
		builder.String("wx_union_id", 32).Default("").Comment("微信UnionId")
	})
	// Write your migrate action here
}

// Down 迁移回滚
func down()  {
	fmt.Println("AddWxUnionIdToUser2021_05_12_212205 Down")
	// Write your rollback action here
	schema.NewSchema().AlterTable("user", func(builder *schema.Builder) {
		builder.DropColumn("wx_union_id")
	})
}


# go-fast-migrate
Like PHP Lavarel Framework, Use go-fast-migrate to Manage your database

### 简介

所谓迁移就像是数据库的版本控制，这种机制允许团队简单轻松的编辑并共享应用的数据库表结构。

迁移可以很容易地构建应用的数据库表结构。如果你曾经频繁告知团队成员需要手动添加列到本地数据库表结构以维护本地开发环境，那么这正是数据库迁移所致力于解决的问题。

###  编译

```go
go build -o migrate main.go
```

### 新建迁移文件

#### 命令

命令：`./migrate new xxxxx`

#### 示例

- 创建迁移文件
```shell script
./migrate new create_user_table
```

执行以上命令后，将会在`migrations`目录生成迁移文件：

```
devops@localhost migrate % ls -al migrations 
total 16
drwxr-xr-x   4 zhoufei  staff   128  5 10 20:30 .
drwxr-xr-x  14 zhoufei  staff   448  5 10 20:38 ..
-rw-r--r--   1 zhoufei  staff  1338  5 10 20:30 2021_05_10_182112_create_user_table.go
-rw-r--r--   1 zhoufei  staff   113  5  9 22:40 init.go
```

生成的迁移文件如下：

```go
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

```

其中，我们需要在`up()`中补充创建表/修改表等操作，在`down()`方法中实现`up()`操作的回退。

以创建用户表操作为例：

```go
func up() {
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

func down() {
	// 删除用户表
	schema.NewSchema().DropTableIfExists("user")
}
```

我们再创建一个修改表操作，`./migrate new add_wx_union_id_to_user`, 修改迁移文件如下：

```go
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

```

### 执行迁移

```shell script
./migrate start
```

执行迁移后，可以看到数据库中增加了`user`表：

```shell script
mysql> desc user;
+---------------------+------------------+------+-----+---------+----------------+
| Field               | Type             | Null | Key | Default | Extra          |
+---------------------+------------------+------+-----+---------+----------------+
| id                  | bigint unsigned  | NO   | PRI | NULL    | auto_increment |
| open_id             | varchar(32)      | NO   | UNI |         |                |
| nickname            | varchar(64)      | NO   |     |         |                |
| avatar              | varchar(255)     | NO   |     |         |                |
| last_login_time     | int unsigned     | NO   |     | 0       |                |
| login_times         | int unsigned     | NO   |     | 0       |                |
| continue_login_days | int unsigned     | NO   |     | 0       |                |
| date                | int unsigned     | NO   |     | 0       |                |
| qq                  | varchar(20)      | NO   | UNI |         |                |
| status              | tinyint unsigned | NO   |     | 0       |                |
| created_at          | timestamp        | YES  |     | NULL    |                |
| updated_at          | timestamp        | YES  |     | NULL    |                |
| wx_union_id         | varchar(32)      | NO   |     |         |                |
+---------------------+------------------+------+-----+---------+----------------+
13 rows in set (0.01 sec)

```


### 执行回滚

```shell script
./migrate rollback
```

执行回滚后,`user`表将被删除。

### 重新迁移

```shell script
./migrate refresh
```

## 其他

- 支持多数据库源

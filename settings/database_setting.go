package settings

import (
	"context"
	"fmt"

	// "github.com/redis/go-redis"
	"github.com/go-redis/redis/v8" // 添加缺失的redis库导入
	"github.com/uozi-tech/cosy/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDB() *gorm.DB {
	return db
}

func InitDatabase(dbs DataBase) {
	switch dbs.Type {
	case "mysql":
		InitMysqlDatabase(dbs)
	case "postgresql":
		InitPostgresqlDatabase(dbs)
	}
}
func InitMysqlDatabase(dbs DataBase) {
	context := context.Background()
	fmt.Println("开始连接数据库mysql.........")
	var err error
	MysqlDsn := (fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbs.User, dbs.Password, dbs.Host, dbs.Port, dbs.Name))
	fmt.Println("MysqlDsn:", MysqlDsn)

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       MysqlDsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		logger.Errorf("failed to connect database: %v", err)
		return
	}
	fmt.Println(context, "Mysql 连接成功")
}

func InitPostgresqlDatabase(dbs DataBase) {
	fmt.Println("开始连接数据库postgresql.........")

	PostgresqlDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", dbs.Host, dbs.User, dbs.Password, dbs.Name, dbs.Port)

	fmt.Println("PostgresqlDsn:", PostgresqlDsn)
	var err error
	// https://github.com/go-gorm/postgres
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  PostgresqlDsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		logger.Errorf("failed to connect database: %v", err)
		return
	}
	fmt.Println(db, "Postgresql 连接成功")
}

func InitRedis(dbs Redis) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", dbs.Host, dbs.Port), //连接地址
		Password: dbs.Password,                             //连接密码
		DB:       int(dbs.DB),                              //默认连接库
		PoolSize: 100,                                      //连接池大小
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Errorf("redis连接失败: %v", err)
		return
	}
	fmt.Println("Redis 连接成功")

}

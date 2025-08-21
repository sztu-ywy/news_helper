package main

import (
	"flag"
	"fmt"
	"log"
	"news_helper/model"

	"github.com/uozi-tech/cosy/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	// 指定生成代码的输出目录（默认："./query"）
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	var confPath string
	// 从配置文件（app.ini）中读取数据库连接信息，并生成数据库连接字符串（DSN）

	flag.StringVar(&confPath, "config", "app.ini", "Specify the configuration file")
	flag.Parse()

	settings.Init(confPath)
	dbs := settings.DataBaseSettings

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbs.User, dbs.Password, dbs.Host, dbs.Port, dbs.Name)

	log.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalln(err)
	}

	g.UseDB(db)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	// 为指定模型生成基本的CRUD方法，model.GenerateAllModel()：生成所有数据库表对应的模型。
	g.ApplyBasic(model.GenerateAllModel()...)

	// apply diy interfaces on structs or table models
	g.ApplyInterface(func(method model.Method) {}, model.GenerateAllModel()...)

	// execute the action of code generation，执行代码生成操作
	g.Execute()
}

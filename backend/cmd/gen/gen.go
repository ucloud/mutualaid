package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// generate code
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../internal/data/model",
		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		/* FieldCoverable: true,*/
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	db, _ := gorm.Open(mysql.Open("mutualaiduser:DemoPassWord@tcp(127.0.0.1:3306)/mutualaid?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db)
	g.GenerateAllTable(
		gen.FieldGORMTag("id", "primaryKey;autoIncrement:false"),
		gen.FieldType("id", "uint64"),
		gen.FieldTypeReg("_id$", "uint64"),
		gen.FieldTypeReg("_time$", "int64"),
		gen.FieldGORMTag("create_time", "autoCreateTime"),
		gen.FieldGORMTag("update_time", "autoUpdateTime"))

	// execute the action of code generation
	g.Execute()
}

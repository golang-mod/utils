package dbgenerate

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"gorm.io/gorm"
)

type Config struct {
	Path        string
	Tables      []string
	TablePrefix string // 去除前缀
}

type generate struct {
	DB           *gorm.DB
	DatabaseName string
	Config
}

var mapping = map[string]string{
	"int":       "int",
	"integer":   "int",
	"mediumint": "int",
	"bit":       "int",
	"year":      "int",
	"smallint":  "int",
	"tinyint":   "int8",
	"bigint":    "int64",
	"decimal":   "float64",
	"double":    "float64",
	"float":     "float64",
	"real":      "float32",
	"numeric":   "float64",
	"timestamp": "time.Time",
	"date":      "*string",
	"datetime":  "time.Time",
	"time":      "time.Time",
}

// Generate 生成Model
// 1.如果非空
// 1.1 如默认值有意义，则设默认值，如字符串类型设为空字符串，数字类型为0
// 1.2 如默认值无意义，不设默认值，除created_at/updated_at外忽略默认值omitempty
// 2.如果可空，设默认值
// 2.1 字符串类型默认空字符串，不加指针*
// 2.2 其他类型默认NULL，添加指针*
func Generate(gdb *gorm.DB, config Config) {
	var gen = new(generate)
	gen.DB = gdb
	gen.Config = config
	gen.DatabaseName = gdb.Migrator().CurrentDatabase()
	tableNamesStr := ""
	for _, name := range config.Tables {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	tables := gen.getTables(tableNamesStr) //生成所有表信息
	log.Printf("tables=%s\n", tables)
	for _, table := range tables {
		fields := gen.getFields(table.Name)
		gen.generateModel(table, fields)
	}
}

// 获取表信息
func (g *generate) getTables(tableNames string) []Table {
	var tables []Table
	if tableNames == "" {
		g.DB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + g.DatabaseName + "';").Find(&tables)
	} else {
		g.DB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + g.DatabaseName + "';").Find(&tables)
	}
	return tables
}

// 获取所有字段信息
func (g *generate) getFields(tableName string) []Field {
	var fields []Field
	g.DB.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

// 生成Model
func (g *generate) generateModel(table Table, fields []Field) {
	content := "package models\n\n"
	tbName := strings.Replace(table.Name, g.TablePrefix, "", 1)
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + Case2Camel(tbName) + " " + table.Comment + "\n"
	}
	content += "type " + Case2Camel(tbName) + " struct {\n"
	//生成字段
	for _, field := range fields {
		// log.Printf(">>>>>>>>>>>>>>>>>>field=%+v\n", field)
		fieldName := Case2Camel(strings.ToLower(field.Field))
		// log.Printf("fieldName=%v ,field.Field= %s", fieldName, field.Field)
		//fieldName := Case2Camel(field.Field)
		//根据字段 获取对应的字段的类型
		fieldType := getFiledType(field)
		//根据字段 获取对应的字段的Orm 表字段信息
		fieldOrm := getFieldOrm(field)
		//根据字段 获取对应的JSON tag 信息
		fieldJson := getFieldJson(field)
		//根据字段  获取对应的备注信息
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldOrm + " " + fieldJson + "` " + fieldComment + "\n"
	}
	content += "}\n"
	content += `func (` + Case2Camel(tbName) + `) TableName() string {
	return "` + table.Name + `"
}`

	filename := "./" + g.Path + "/" + tbName + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		//if !conf.ModelReplace {
		fmt.Println(tbName + " 已存在，需删除才能重新生成...")
		return
		//}
		//f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //打开文件
		//if err != nil {
		//	panic(err)
		//}
	}
	//创建目录

	if !checkFileIsExist(g.Path) {
		err = os.Mkdir(g.Path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //打开文件
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	fmt.Println(content)
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		//导入 需要的系统 包
		_, err = exec.Command("go", "fmt", filename).Output()
		if err != nil {
			log.Println(err)
		}
		_, err = exec.Command("goimports", "-l", "-w", filename).Output()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(Case2Camel(tbName) + " 已生成...")
	}
}

// 获取字段类型
func getFiledType(field Field) string {
	if field.Field == "deleted_at" {
		return "gorm.DeletedAt"
	}
	if field.Field == "deleted_at" {
		return "gorm.DeletedAt"
	}
	typeArr := strings.Split(field.Type, "(")
	value, ok := mapping[typeArr[0]]
	if ok {
		if field.Default == nil && value != "string" && value != "*string" && field.Null != "NO" {
			return "*" + value
		}
		return value
	} else {
		return "string"
	}
}

// 获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + strings.ToLower(field.Field) + `"`
}

// 获取字段gorm描述
func getFieldOrm(field Field) string {
	var fieldField string
	if field.Field != "" {
		fieldField = "column:" + field.Field
	}
	var primaryKeyField string
	if field.Key == "PRI" {
		primaryKeyField = ";primaryKey"
	}
	var typeField string
	if field.Type != "" {
		typeField = ";type:" + field.Type
	}
	var defaultField string
	if field.Default != nil {
		defaultField = ";default:" + *field.Default + ""
	} else {
		// 没设默认值并且非空
		if field.Null == "NO" && field.Field != "created_at" && field.Field != "updated_at" {
			defaultField = ";default:;omitempty"
		}
	}
	var nullField string
	if field.Null == "NO" {
		nullField = ";not null"
		//if field.Default == nil && field.Field != "created_at" && field.Field != "updated_at" {
		//	nullField += ";default:;omitempty"
		//}
	}
	return `gorm:"` + fieldField + primaryKeyField + typeField + nullField + defaultField + `"`

}

// 获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

// 检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// Case2Camel 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = cases.Title(language.English).String(name)
	return strings.Replace(name, " ", "", -1)
}

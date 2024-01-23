package gormx

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	Host             string
	Port             string
	Database         string
	Username         string
	Password         string
	Charset          string
	Timeout          string
	TablePrefix      string
	LogMode          int
	MaxOpenConns     int
	MaxIdleConns     int
	MaxLifetimeConns int64
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Connection(cfg Config) (gdb *gorm.DB, err error) {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local&timeout=%s&parseTime=True",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.Timeout,
	)
	var mysqlConfig = mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}
	var gormConfig = gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀
			SingularTable: true, // 使用单数表名，启用该选项后，`UserInfo` 表将是`user`
			//NameReplacer: strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		},
	}
	// 连接数据库
	gdb, err = gorm.Open(mysql.New(mysqlConfig), &gormConfig)
	if err != nil {
		return nil, err
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	db, _ := gdb.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	//db.SetMaxIdleConns(10)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	//db.SetMaxOpenConns(100)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	//db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeConns) * time.Second)

	return gdb, nil
}

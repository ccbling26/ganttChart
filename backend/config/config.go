package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

type ServiceConfig struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type DBConfig struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	User                string `mapstructure:"user" json:"user" yaml:"user"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"" json:"charset" yaml:"charset"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdleConn         int    `mapstructure:"max_idle_conn" json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn         int    `mapstructure:"max_open_conn" json:"max_open_conn" yaml:"max_open_conn"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

type Config struct {
	Service ServiceConfig `mapstructure:"service" json:"service" yaml:"service"`
	DB      DBConfig      `mapstructure:"db" json:"db" yaml:"db"`
	Log     Log           `mapstructure:"log" json:"log" yaml:"log"`
}

type GlobalConfig struct {
	ProjectPath string
	Config      Config
	ConfigViper *viper.Viper
	DB          *gorm.DB
	Logger      *zap.Logger
}

var Global = new(GlobalConfig)

func initialConfig() {
	// 配置文件路径
	configFile := fmt.Sprintf("%s/config/config.yaml", Global.ProjectPath)
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		configFile = configEnv
	}

	// 初始化 Viper
	vip := viper.New()
	vip.SetConfigFile(configFile)
	vip.SetConfigType("yaml")
	if err := vip.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Read config file failed! Details: %s\n", err))
	}

	// 监听配置文件
	vip.WatchConfig()
	vip.OnConfigChange(func(in fsnotify.Event) {
		Global.Logger.Info("Config File changed: " + in.Name)
		// 重新配置
		if err := vip.Unmarshal(&Global.Config); err != nil {
			Global.Logger.Error(err.Error())
		}
	})

	// 设置全局变量
	if err := vip.Unmarshal(&Global.Config); err != nil {
		panic(err)
	}
	Global.ConfigViper = vip
}

func initialDB() {
	switch Global.Config.DB.Driver {
	case "mysql":
		initialMySQL()
	default:
		panic("只支持 MySQL 目前")
	}
}

func initialMySQL() {
	// parameter details: https://github.com/go-sql-driver/mysql#parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Global.Config.DB.User,
		Global.Config.DB.Password,
		Global.Config.DB.Host,
		Global.Config.DB.Port,
		Global.Config.DB.Database,
	)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁止自动创建外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 使用自定义 Logger
		Logger: getGormLogger(),
	}); err != nil {
		Global.Logger.Error("MySQL init failed! Details: ", zap.Any("err", err))
		panic(fmt.Errorf("MySQL init failed! Details: %s\n", err))
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(Global.Config.DB.MaxIdleConn)
		sqlDB.SetMaxOpenConns(Global.Config.DB.MaxOpenConn)

		Global.DB = db
	}
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch Global.Config.DB.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		// 慢 SQL 阈值
		SlowThreshold: 200 * time.Microsecond,
		// 日志级别
		LogLevel: logMode,
		// 是否忽略 ErrRecordNotFound（记录未找到错误）
		IgnoreRecordNotFoundError: false,
		// 是否允许彩色打印
		Colorful: !Global.Config.DB.EnableFileLogWriter,
	})
}

func getGormLogWriter() logger.Writer {
	var writer io.Writer
	if Global.Config.DB.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   Global.Config.Log.RootDir + "/" + Global.Config.DB.LogFilename,
			MaxSize:    Global.Config.Log.MaxSize,
			MaxBackups: Global.Config.Log.MaxBackups,
			MaxAge:     Global.Config.Log.MaxAge,
			Compress:   Global.Config.Log.Compress,
		}
	} else {
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func initialLog() {
	createRootDir()
	setLogLevel()
	if Global.Config.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}
	Global.Logger = zap.New(getZapCore(), options...)
}

func InitialGlobal() {
	// 项目路径
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Global.ProjectPath = projectPath

	// 加载配置
	initialConfig()

	// 初始化 logger
	initialLog()

	// 初始化数据库
	initialDB()
}

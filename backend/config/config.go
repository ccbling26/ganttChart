package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type ServiceConfig struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type DBConfig struct {
	IP          string `mapstructure:"ip" json:"ip" yaml:"ip"`
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`
	User        string `mapstructure:"user" json:"user" yaml:"user"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	Database    string `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdleConn int    `mapstructure:"maxIdleConn" json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"maxOpenConn" json:"max_open_conn" yaml:"max_open_conn"`
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

func initialMySQL() {
	// parameter details: https://github.com/go-sql-driver/mysql#parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Global.Config.DB.User,
		Global.Config.DB.Password,
		Global.Config.DB.IP,
		Global.Config.DB.Port,
		Global.Config.DB.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("MySQL init failed! Details: %s\n", err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(Global.Config.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(Global.Config.DB.MaxOpenConn)

	Global.DB = db
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

	// 加载 MySQL
	initialMySQL()

	// 初始化 logger
	initialLog()
}

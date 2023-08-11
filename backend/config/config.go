package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type ServiceConfig struct {
	Addr string `mapstructure:"addr" yaml:"addr"`
}

type DBConfig struct {
	IP          string `mapstructure:"ip" yaml:"ip"`
	Port        int    `mapstructure:"port" yaml:"port"`
	User        string `mapstructure:"user" yaml:"user"`
	Password    string `mapstructure:"password" yaml:"password"`
	Database    string `mapstructure:"database" yaml:"database"`
	MaxIdleConn int    `mapstructure:"maxIdleConn" yaml:"maxIdleConn"`
	MaxOpenConn int    `mapstructure:"maxOpenConn" yaml:"maxOpenConn"`
}

type Config struct {
	Service ServiceConfig `mapstructure:"service" yaml:"service"`
	DB      DBConfig      `mapstructure:"db" yaml:"db"`
}

type GlobalConfig struct {
	ProjectPath string
	Config      Config
	ConfigViper *viper.Viper
	DB          *gorm.DB
}

var Global = new(GlobalConfig)

func InitialConfig() {
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
		panic(fmt.Errorf("Read config file failed! Details: %s \n", err))
	}

	// 监听配置文件
	vip.WatchConfig()
	vip.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config File changed: ", in.Name)
		// 重新配置
		if err := vip.Unmarshal(&Global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 设置全局变量
	if err := vip.Unmarshal(&Global.Config); err != nil {
		fmt.Println(err)
	}
	Global.ConfigViper = vip
}

func InitialMySQL() {
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
		panic(err.Error())
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(Global.Config.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(Global.Config.DB.MaxOpenConn)

	Global.DB = db
}

func InitialGlobal() {
	// 项目路径
	projectPath, _ := os.Getwd()
	Global.ProjectPath = projectPath

	// 加载配置
	InitialConfig()

	// 加载 MySQL
	InitialMySQL()
}

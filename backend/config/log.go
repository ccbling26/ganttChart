package config

import (
	"backend/src/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Log struct {
	// 日志等级
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 日志根目录
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	// 日志文件名称
	Filename string `mapstructure:"filename" json:"filename" yaml:"filename"`
	// 写入格式，可选 json
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// 是否显示调用行
	ShowLine bool `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	// 旧文件的最大个数
	MaxBackups int `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	// 日志文件最大大小（MB）
	MaxSize int `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	// 旧文件的最大保留天数
	MaxAge int `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	// 是否压缩
	Compress bool `mapstructure:"compress" json:"compress" yaml:"compress"`
}

var (
	// zap 日志等级
	level zapcore.Level
	// zap 配置项
	options []zap.Option
)

func createRootDir() {
	if ok, _ := utils.PathExists(Global.Config.Log.RootDir); !ok {
		_ = os.Mkdir(Global.Config.Log.RootDir, os.ModePerm)
	}
}

func setLogLevel() {
	switch Global.Config.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "d_panic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

func getZapCore() zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(Global.Config.Service.Env + "." + level.String())
	}

	// 设置编码器
	if Global.Config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   Global.Config.Log.RootDir + "/" + Global.Config.Log.Filename,
		MaxSize:    Global.Config.Log.MaxSize,
		MaxBackups: Global.Config.Log.MaxBackups,
		MaxAge:     Global.Config.Log.MaxAge,
		Compress:   Global.Config.Log.Compress,
	}
	return zapcore.AddSync(file)
}

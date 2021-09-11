package hy_zap

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapConfig struct {
	LogPath      string // 日志文件路径
	LogLevel     string //日志级别 debug/info/warn/error
	MaxSize      int    //单个文件大小,MB
	MaxBackups   int    // 保存的文件个数
	MaxAge       int    // 保存的天数， 没有的话不删除
	Compress     bool   // 压缩
	JsonFormat   bool   // 是否输出为json格式
	ShowLine     bool   // 显示代码行
	LogInConsole bool   // 是否同时输出到控制台
	ServerName   string // 服务名称
}

func DefaultConfig() *ZapConfig {
	return &ZapConfig{
		LogPath:      "./logs/server.log",
		LogLevel:     "debug",
		MaxSize:      500,
		MaxBackups:   10,
		MaxAge:       7,
		Compress:     false,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
		ServerName:   "server",
	}
}

var ZapLog *zap.SugaredLogger

func InitLogger(conf *ZapConfig) {
	if conf == nil {
		conf = DefaultConfig()
	}

	hook := lumberjack.Logger{
		Filename:   conf.LogPath,    // 日志文件路径
		MaxSize:    conf.MaxSize,    // megabytes
		MaxBackups: conf.MaxBackups, // 最多保留300个备份
		Compress:   conf.Compress,   // 是否压缩 disabled by default
	}
	if conf.MaxAge > 0 {
		hook.MaxAge = conf.MaxAge // days
	}

	var syncer zapcore.WriteSyncer
	if conf.LogInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if conf.JsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch conf.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	log := zap.New(core)
	if len(conf.ServerName) > 0 {
		log = log.WithOptions(zap.Fields(zap.String("serviceName", conf.ServerName)))
	}
	if conf.ShowLine {
		log = log.WithOptions(zap.Development(), zap.AddCaller())
	}
	ZapLog = log.Sugar()
}

//Debug 打印"Debug"级别日志信息
func Debug(v ...interface{}) {
	ZapLog.Debug(v...)
}

//Debugf 打印"Debug"级别日志信息
func Debugf(format string, v ...interface{}) {
	ZapLog.Debugf(format, v...)
}

//Info 打印"Info"级别日志信息
func Info(v ...interface{}) {
	ZapLog.Info(v...)
}

//Infof 打印"Info"级别日志信息
func Infof(format string, v ...interface{}) {
	ZapLog.Infof(format, v...)
}

//Warn 打印"Warn"级别日志信息
func Warn(v ...interface{}) {
	ZapLog.Warn(v...)
}

//Warnf 打印"Warn"级别日志信息
func Warnf(format string, v ...interface{}) {
	ZapLog.Warnf(format, v...)
}

//Error 打印"Error"级别日志信息
func Error(v ...interface{}) {
	ZapLog.Error(v...)
}

//Errorf 打印"Error"级别日志信息
func Errorf(format string, v ...interface{}) {
	ZapLog.Errorf(format, v...)
}

// Print 打印info级别日志
func Print(v ...interface{}) {
	ZapLog.Info(v...)
}

// Print 打印info级别日志
func Printf(format string, v ...interface{}) {
	ZapLog.Infof(format, v...)
}

// Print 打印info级别日志
func Println(v ...interface{}) {
	ZapLog.Info(v...)
}

// Panic 打印"Panic"级别日志信息
func Panic(v ...interface{}) {
	ZapLog.Panic(v...)
}

//Panicln 打印"Panic"级别日志信息
func Panicln(v ...interface{}) {
	ZapLog.Panic(v...)
}

//Panicf 打印"Panic"级别日志信息
func Panicf(format string, v ...interface{}) {
	ZapLog.Panicf(format, v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	ZapLog.Fatal(v...)
}

//Fatalln 打印"Fatal"级别日志信息
func Fatalln(v ...interface{}) {
	ZapLog.Fatal(v...)
}

//Fatalf 打印"Fatal"级别日志信息
func Fatalf(format string, v ...interface{}) {
	ZapLog.Fatalf(format, v...)
}

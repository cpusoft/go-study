package zaplog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/osutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

// var sugarLogger := logger.Sugar()
// var sugarLogger
func init() {

	var logLevelStr string
	logLevel := conf.String("logs::level")

	//DEBUG<INFO<WARN<ERROR<FATAL
	switch logLevel {
	case "LevelEmergency":
		fallthrough
	case "LevelAlert":
		fallthrough
	case "LevelCritical":
		fallthrough
	case "LevelError":
		logLevelStr = "error"
	case "LevelWarning":
		fallthrough
	case "LevelNotice":
		fallthrough
	case "LevelInformational":
		logLevelStr = "info"
	case "LevelDebug":
		logLevelStr = "debug"
	}
	if logLevel == "" {
		logLevel = "info" // default level
	}
	// get level
	// get process file name as log name
	logName := filepath.Base(os.Args[0])
	if logName != "" {
		logName = strings.Split(logName, ".")[0] + ".log"
	} else {
		logName = conf.String("logs::name")
	}
	path, err := osutil.GetCurrentOrParentAbsolutePath("log")
	if err != nil {
		fmt.Println("found " + path + " failed, " + err.Error())
		return
	}
	filePath := path + string(os.PathSeparator) + logName
	fmt.Println(filePath)
	lc := LogConfig{
		Level:    logLevelStr, // DEBUG<INFO<WARN<ERROR<FATAL
		FileName: filePath,
	}
	err = initLogger(lc)
	if err != nil {
		fmt.Println(err)
	}
	// L()：获取全局logger
	logger = zap.L()
	sugaredLogger = logger.Sugar()
}

type LogConfig struct {
	Level      string `json:"level"`       // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string `json:"file_name"`   // FileName 日志文件位置
	MaxSize    int    `json:"max_size"`    // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int    `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int    `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
}

// 负责设置 encoding 的日志格式
func getEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConfig := zap.NewProductionEncoderConfig()

	// 设置每个日志条目使用的键。如果有任何键为空，则省略该条目的部分。

	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "time"
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 以 package/file:行 的格式 序列化调用程序，从完整路径中删除除最后一个目录外的所有目录。
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

// 负责日志写入的位置
func getLogWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer。
	// 它试图变得智能：如果 io.Writer 的具体类型实现了 WriteSyncer，我们将使用现有的 Sync 方法。
	// 如果没有，我们将添加一个无操作同步。

	return zapcore.AddSync(lumberJackLogger)
}

// initLogger 初始化Logger
func initLogger(lCfg LogConfig) (err error) {
	// 获取日志写入位置
	writeSyncer := getLogWriter(lCfg.FileName, lCfg.MaxSize, lCfg.MaxBackups, lCfg.MaxAge)
	// 获取日志编码格式
	encoder := getEncoder()

	// 获取日志最低等级，即>=该等级，才会被写入。
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(lCfg.Level))
	if err != nil {
		fmt.Println("initLogger(): UnmarshalText Level fail:", lCfg.Level)
		return
	}

	// 创建一个将日志写入 WriteSyncer 的核心。
	core := zapcore.NewCore(encoder, writeSyncer, l)
	logger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel))

	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(logger)
	return
}

type Field = zap.Field
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserName string `json:"userName"`
	UserId   uint64 `json:"userId"`
	TraceId  string `json:"traceId"`
	LogId    uint64 `json:"LogId"`
}

// Debug in Json:DebugJ("msg", zap.String("aa","bb"), zap.Int("id",33)) -> ["aa","bb","id",33]
func DebugJ(cx context.Context, msg string, fields ...Field) {
	cc := cx.Value("CustomClaims")
	customClaims := cc.(CustomClaims)
	fields = append(fields, zap.String("userName", customClaims.UserName))
	fields = append(fields, zap.Uint64("userId", customClaims.UserId))
	fields = append(fields, zap.String("traceId", customClaims.TraceId))
	fields = append(fields, zap.Uint64("logId", customClaims.LogId))
	logger.Debug(msg, fields...)
}

// Debug with Json: DebugJw("msg", "aaa","bbb", "id",33) -> ["aa","bb","id",33]
func DebugJw(cx context.Context, msg string, args ...interface{}) {
	cc := cx.Value("CustomClaims")
	customClaims := cc.(CustomClaims)
	args = append(args, "userName", customClaims.UserName)
	args = append(args, "userId", customClaims.UserId)
	args = append(args, "traceId", customClaims.TraceId)
	args = append(args, "logId", customClaims.LogId)
	sugaredLogger.Debugw(msg, args...)
}

// Debug in Line: DebugL("msg","aaa","bbb", "id",33) -> msg aaa bbb id  33
func DebugL(cx context.Context, msg string, args ...interface{}) {
	sugaredLogger.Debugw(msg + " " + convert.Interfaces2String(args))
}

// Info in Json:DebugJ("msg", zap.String("aa","bb"), zap.Int("id",33)) -> ["aa","bb","id",33]
func InfoJ(cx context.Context, msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

// Info wit Json: DebugJw("msg", "aaa","bbb", "id",33) -> ["aa","bb","id",33]
func InfoJw(cx context.Context, msg string, args ...interface{}) {
	sugaredLogger.Infow(msg, args...)
}

// Info in Line: DebugL("msg","aaa","bbb", "id",33) -> msg aaa bbb id 33
func InfoL(cx context.Context, msg string, args ...interface{}) {
	sugaredLogger.Infow(msg + " " + convert.Interfaces2String(args))
}

func ErrorJ(cx context.Context, msg string, fields ...Field) {
	logger.Error(msg, fields...)
}
func ErrorJw(cx context.Context, msg string, args ...interface{}) {
	sugaredLogger.Errorw(msg, args...)
}
func ErrorL(cx context.Context, msg string, args ...interface{}) {
	sugaredLogger.Errorw(msg + " " + convert.Interfaces2String(args))
}

func DeferSync() {
	// 调用内核的Sync方法，刷新所有缓冲的日志条目。
	// 应用程序应该注意在退出之前调用Sync。
	logger.Sync()
	sugaredLogger.Sync()
}

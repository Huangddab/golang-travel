package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
)

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	}
	return "UNKNOWN"
}

// WithLevel 设置日志级别
type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// 设置日志公共字段
func (l *Logger) WithFileds(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

// 设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

// 设置当前某一层调用栈的信息
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()                            // 克隆Logger实例
	pc, file, line, ok := runtime.Caller(skip) // 获取调用者信息
	if ok {                                    // 如果获取成功
		f := runtime.FuncForPC(pc)                                            // 获取函数信息
		ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())} // 格式化调用者信息
	}
	return ll
}

// 设置完整调用栈的信息
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25                                                  // 最大调用层级
	minCallerDepth := 1                                                   // 最小调用层级
	callers := []string{}                                                 // 存储调用栈信息
	pcs := make([]uintptr, maxCallerDepth)                                // program counters
	depth := runtime.Callers(minCallerDepth, pcs)                         // 获取调用栈
	frames := runtime.CallersFrames(pcs[:depth])                          // 获取调用栈帧
	for frame, more := frames.Next(); more; frame, more = frames.Next() { // 遍历每一帧
		callers = append(callers, fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)) // 格式化调用栈信息
	}
	ll := l.clone()      // 克隆Logger实例
	ll.callers = callers // 设置调用栈信息
	return ll
}

// 格式化日志内容
func (l *Logger) JSONFormat(lv Level, msg string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)                    // 初始化日志数据容器
	data["level"] = lv.String()                              // 设置日志级别
	data["time"] = fmt.Sprintf("%d", runtime.NumGoroutine()) // 设置日志时间
	data["message"] = msg                                    // 设置日志消息
	data["callers"] = l.callers                              // 设置调用栈信息
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}
func (l *Logger) Output(lev Level, msg string) {
	body, _ := json.Marshal(l.JSONFormat(lev, msg))
	content := string(body)
	switch lev {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

// 日志分级输出
func (l *Logger) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Logger) Debug(v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}
func (l *Logger) Warn(v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}
func (l *Logger) Error(v ...interface{}) {
	l.Output(LevelError, fmt.Sprint(v...))
}
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}
func (l *Logger) Panic(v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprint(v...))
}

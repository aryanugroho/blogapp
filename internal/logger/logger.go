package logger

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/aryanugroho/blogapp/internal/contextprop"
	"go.uber.org/zap"
)

var zapLogger *zap.Logger

func Init() {
	zapLogger, _ = zap.NewProduction()
	zapLogger = zapLogger.WithOptions(zap.WithCaller(false))
}

func getLogger() *zap.Logger {
	if zapLogger == nil {
		Init()
	}
	return zapLogger
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(ctx, "cid")),
	}
	fields = append(fields, fileField, cidField)
	getLogger().Debug(msg, convertFields(fields)...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(ctx, "cid")),
	}
	fields = append(fields, fileField, cidField)
	getLogger().Info(msg, convertFields(fields)...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(ctx, "cid")),
	}
	fields = append(fields, fileField, cidField)
	getLogger().Warn(msg, convertFields(fields)...)
}

func Error(ctx context.Context, msg string, err error, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	errField := Field{
		zap.String("err", err.Error()),
	}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(ctx, "cid")),
	}
	fields = append(fields, fileField, errField, cidField)
	getLogger().Error(msg, convertFields(fields)...)
}

func Fatal(ctx context.Context, msg string, err error, fields ...Field) {
	file, line := getFileAndLine()
	fileNLine := fmt.Sprintf("%s:%d", file, line)
	fileField := Field{
		zap.String("file", fileNLine),
	}
	errField := Field{
		zap.String("err", err.Error()),
	}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(ctx, "cid")),
	}
	fields = append(fields, fileField, errField, cidField)
	getLogger().Fatal(msg, convertFields(fields)...)
}

func Http(r *http.Request) {
	fields := []Field{}
	cidField := Field{
		zap.String("cid", contextprop.GetContextValue(r.Context(), "cid")),
	}
	clientID := Field{
		zap.String("client_id", contextprop.GetContextValue(r.Context(), "client_id")),
	}
	pathField := Field{
		zap.String("path", r.URL.Path),
	}
	methodField := Field{
		zap.String("method", r.Method),
	}
	ipField := Field{
		zap.String("ip", r.RemoteAddr),
	}

	fields = append(fields, cidField, clientID, pathField, methodField, ipField)
	Debug(r.Context(), "http-request", fields...)
}

const (
	skip = 2
)

func getFileAndLine() (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		i := strings.Index(file, "blogapp")
		file = file[i:]
	}

	return file, line
}

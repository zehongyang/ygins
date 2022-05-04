package logger

import (
	"go.uber.org/zap"
	"log"
	"ygin/config"
)

var (
	gLogger *zap.Logger
	lc LogConfig
	err error
)

type LogConfig struct {
	Log zap.Config
}

func init()  {
	config.Load(&lc)
	gLogger, err = lc.Log.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
}

func D(fields ...zap.Field)  {
	gLogger.Debug("",fields...)
}


func Debug(msg string, fields ...zap.Field)  {
	gLogger.Debug(msg,fields...)
}

func I(fields ...zap.Field)  {
	gLogger.Info("",fields...)
}

func Info(msg string, fields ...zap.Field)  {
	gLogger.Info(msg,fields...)
}

func W(fields ...zap.Field)  {
	gLogger.Warn("",fields...)
}

func Warn(msg string, fields ...zap.Field)  {
	gLogger.Warn(msg,fields...)
}

func E(fields ...zap.Field)  {
	gLogger.Error("",fields...)
}

func Error(msg string, fields ...zap.Field)  {
	gLogger.Error(msg,fields...)
}

func F(fields ...zap.Field)  {
	gLogger.Fatal("",fields...)
}

func Fatal(msg string, fields ...zap.Field)  {
	gLogger.Fatal(msg,fields...)
}
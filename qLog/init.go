package qLog

import (
	"os"

	"github.com/quincy0/qpro/qConfig"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(cfgLog qConfig.Log) {
	rotete := &lumberjack.Logger{
		Filename:   qConfig.Settings.Log.Path,
		MaxSize:    qConfig.Settings.Log.MaxSize,
		MaxAge:     qConfig.Settings.Log.MaxAge,
		MaxBackups: qConfig.Settings.Log.MaxBackups,
		LocalTime:  false,
		Compress:   cfgLog.Compress,
	}

	level, err := zapcore.ParseLevel(cfgLog.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	ws := []zapcore.WriteSyncer{
		zapcore.AddSync(rotete),
	}
	if cfgLog.ConsoleStdout {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	SetDefaultLogger(
		MustNew(
			OptLevel(level),
			OptWriteSyncers(ws...),
		),
	)
}

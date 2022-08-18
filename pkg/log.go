package pkg

import (
	"time"

	"github.com/pyihe/plogs"
)

type Logger interface {
	Panic(...interface{})
	Panicf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Close()
}

func InitLogger(app string) Logger {
	opts := []plogs.Option{
		plogs.WithName(app),
		plogs.WithLogPath("./logs"),
		plogs.WithFileOption(plogs.WriteByLevelMerged),
		plogs.WithStdout(true),
		plogs.WithLogLevel(plogs.LevelInfo | plogs.LevelDebug | plogs.LevelWarn | plogs.LevelError | plogs.LevelFatal | plogs.LevelPanic),
		plogs.WithMaxAge(24 * time.Hour),
		plogs.WithMaxSize(60 * 1024 * 1024),
	}
	return plogs.NewLogger(opts...)
}

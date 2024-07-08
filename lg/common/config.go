package common

type LogConfig struct {
	DisableToFile bool   `desc:"disable log output to file (default false)"`
	FileName      string `desc:"output filename (default runlog.log)"`
	MaxSize       int    `desc:"file max szie (default 3)"`
	MaxBackup     int    `desc:"max backup count (default 3)"`
	MaxAge        int    `desc:"max backup age (default 30)"`
	Compress      bool   `desc:"whether to use compress (default false)"`
}

func (l *LogConfig) SetDefault() {
	if l.FileName == "" {
		l.FileName = "runlog.log"
	}

	if l.MaxAge == 0 {
		l.MaxAge = 30
	}

	if l.MaxBackup == 0 {
		l.MaxBackup = 3
	}

	if l.MaxSize == 0 {
		l.MaxSize = 3
	}
}

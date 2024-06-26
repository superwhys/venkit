package common

type LogConfig struct {
	FileName  string `desc:"output filename (default runlog.log)"`
	MaxSize   int    `desc:"file max szie (default 3)"`
	MaxBackup int    `desc:"max backup count (default 3)"`
	MaxAge    int    `desc:"max backup age (default 30)"`
	Compress  bool   `desc:"whether to use compress (default false)"`
}

func (l *LogConfig) SetDefault() {
	l.FileName = "runlog.log"
	l.MaxSize = 3
	l.MaxBackup = 3
	l.MaxAge = 30
	l.Compress = false
}

package log

type Config struct {
	Filename   string   `yaml:"filename"`
	Level      string   `yaml:"level"`
	Targets    []string `yaml:"targets"`
	MaxSize    int      `yaml:"max_size"`
	MaxBackups int      `yaml:"max_backups"`
	Compress   bool     `yaml:"compress"`
}

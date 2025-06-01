package mailer

type Config struct {
	Host      string            `yaml:"host"`
	Port      int               `yaml:"port"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	Sender    string            `yaml:"sender"`
	Templates map[string]string `yaml:"templates"`
}

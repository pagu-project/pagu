package whatsapp

type Config struct {
	WebHookToken   string  `yaml:"webhook_token"`
	GraphToken     string  `yaml:"graph_token"`
	WebHookAddress string  `yaml:"webhook_address"`
	WebHookPath    string  `yaml:"webhook_path"`
	Session        Session `yaml:"session"`
}

type Session struct {
	SessionTTL    int `yaml:"session_ttl"`
	CheckInterval int `yaml:"check_interval"`
}

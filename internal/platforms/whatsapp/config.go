package whatsapp

type Config struct {
	WebHookToken   string  `yaml:"webhook_token"`
	AccessToken    string  `yaml:"access_token"`
	WebHookAddress string  `yaml:"webhook_address"`
	WebHookPath    string  `yaml:"webhook_path"`
	Session        Session `yaml:"session"`
}

type Session struct {
	SessionTTL    int `yaml:"session_ttl"`
	CheckInterval int `yaml:"check_interval"`
}

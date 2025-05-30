package nowpayments

type Config struct {
	Webhook       string `yaml:"webhook"`
	ListenPort    string `yaml:"listen_port"`
	APIToken      string `yaml:"api_token"`
	IPNSecret     string `yaml:"ipn_secret"`
	APIURL        string `yaml:"api_url"` //nolint:tagliatelle // unable to detect api_url
	PaymentURL    string `yaml:"payment_url"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	FeePaidByUser bool   `yaml:"fee_paid_by_user"`
	FixedRate     bool   `yaml:"fixed_rate"`
}

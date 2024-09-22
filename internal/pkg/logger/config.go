package logger

type Config struct {
	TgBotToken    string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	TgBotBaseURL  string `envconfig:"TELEGRAM_BOT_API_BASE_URL" required:"true"`
	SelfIPAddress string `envconfig:"SELF_IP_ADDRESS" default:""`
	IsCronJob     bool   `envconfig:"IS_CRON_JOB" required:"true"`
}

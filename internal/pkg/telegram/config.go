package telegram

type Config struct {
	TgBotToken   string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	TgBotBaseURL string `envconfig:"TELEGRAM_BOT_API_BASE_URL" required:"true"`
}

func (c *Config) GetTgBotToken() string   { return c.TgBotToken }
func (c *Config) GetTgBotBaseURL() string { return c.TgBotBaseURL }

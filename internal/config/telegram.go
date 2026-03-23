package config

type (
	TelegramConfig struct {
		ApiId      int    `mapstructure:"api_id"`
		ApiHash    string `mapstructure:"api_hash"`
		SessionDir string `mapstructure:"session_dir"`
	}
)

func NewTelegramConfig() *TelegramConfig {
	return new(TelegramConfig)
}

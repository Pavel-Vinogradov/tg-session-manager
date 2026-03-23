package config

type (
	TelegramConfig struct {
		ApiId   int    `mapstructure:"api_id"`
		ApiHash string `mapstructure:"api_hash"`
	}
)

func NewTelegramConfig() *TelegramConfig {
	return new(TelegramConfig)
}

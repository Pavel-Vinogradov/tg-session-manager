package config

import "github.com/gotd/td/telegram"

type (
	TelegramService struct {
		client *telegram.Client
		config AppConfig
	}
	TelegramConfig struct {
		ApiId   int    `mapstructure:"api_id"`
		ApiHash string `mapstructure:"api_hash"`
	}
)

func NewTelegramService() *ServiceChainMember {
	return &ServiceChainMember{
		&TelegramService{},
	}
}

func (t *TelegramService) Health() error {
	return nil
}

func (t *TelegramService) setupMember(services *Services) error {
	if services.App != nil && services.App.Config.TelegramServer != nil {
		t.config = services.App.Config
		t.client = telegram.NewClient(
			t.config.TelegramServer.ApiId,
			t.config.TelegramServer.ApiHash,
			telegram.Options{},
		)
	}

	services.Telegram = t
	return nil
}

func (t *TelegramService) appendMemberToServices(services *Services) {
	services.Telegram = t
}

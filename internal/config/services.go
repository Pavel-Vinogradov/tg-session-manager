package config

type (
	Services struct {
		App        *AppService
		GrpcServer *GrpcServer
		Telegram   *TelegramService
		Warnings   []string
	}

	union struct {
		Internal *Services
	}

	configurator struct {
		member serviceConfigurator
		union  *union
	}

	serviceConfigurator interface {
		Health() error
		setupMember(s *Services) error
		appendMemberToServices(*Services)
	}
	ServiceChainMember struct {
		serviceConfigurator
	}
	memberHandler func()
)

func NewConfigurations() (s *Services) {
	unions := newUnion()
	handlers := []memberHandler{
		func() {
			runServiceChain([]configurator{
				{member: NewAppService(), union: unions},
				{member: NewGrpcServer(), union: unions},
				{member: NewTelegramService(), union: unions},
			}...)
		},
	}
	return unions.withServicesChainHandler(handlers...)
}

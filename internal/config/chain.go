package config

import (
	"errors"

	"github.com/sirupsen/logrus"
)

func newUnion() *union {
	return &union{
		Internal: &Services{},
	}
}

func (u union) withServicesChainHandler(handlers ...memberHandler) *Services {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("error while execute one of handlers in servicesChain: %v", r)
		}
	}()

	for _, handler := range handlers {
		handler()
	}

	return u.Internal
}

func runServiceChain(members ...configurator) {
	for _, m := range members {
		err := m.member.setupMember(m.union.Internal)

		if err != nil {
			var appServiceError AppServiceError
			if errors.As(err, &appServiceError) {
				panic(appServiceError)
			}

			m.union.Internal.Warnings = append(m.union.Internal.Warnings, err.Error())
			logrus.Warn(err)
		}

		m.member.appendMemberToServices(m.union.Internal)
	}
}

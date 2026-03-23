package telegram

import "errors"

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionNotActive = errors.New("session is not active")
)

type Repository interface {
	DeleteSession(sessionId string) error
	SendMessage(sessionId, peer, message string) error
}

type Service interface {
	CreateSession() (string, string, error)
	DeleteSession(sessionId string) error
	SendMessage(sessionId, peer, message string) error
}

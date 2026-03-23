package telegram

type Repository interface {
	CreateSession() (string, string, error)
	DeleteSession(sessionId string) error
	SendMessage(sessionId, message string) error
}

type Service interface {
	CreateSession() (string, string, error)
	DeleteSession(sessionId string) error
	SendMessage(sessionId, message string) error
}

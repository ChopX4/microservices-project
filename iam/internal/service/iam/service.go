package iam

import (
	"time"

	"github.com/ChopX4/raketka/iam/internal/repository"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
)

type service struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	txManager         pgxtx.TxManager
	sessionTTL        time.Duration
}

func NewIamService(sessionRepository repository.SessionRepository, userRepository repository.UserRepository, txManager pgxtx.TxManager, sessionTTL time.Duration) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		txManager:         txManager,
		sessionTTL:        sessionTTL,
	}
}

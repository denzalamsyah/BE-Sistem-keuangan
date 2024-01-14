package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type SessionService interface {
	GetSessionByEmail(email string) (models.Session, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) *sessionService {
	return &sessionService{sessionRepo}
}

func (c *sessionService) GetSessionByEmail(email string) (models.Session, error) {
	
	session, err := c.sessionRepo.SessionAvailEmail(email)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

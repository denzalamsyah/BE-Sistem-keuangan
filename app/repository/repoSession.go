package repository

import (
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSessions(session models.Session) error
	DeleteSession(token string) error
	UpdateSessions(session models.Session) error
	SessionAvailEmail(email string) (models.Session, error)
	SessionAvailToken(token string) (models.Session, error)
	TokenExpired(session models.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) *sessionsRepo {
	return &sessionsRepo{db}
}

func (u *sessionsRepo) AddSessions(session models.Session) error {
	err := u.db.Create(&session).Error
	if err != nil{
		return err
	}
	return nil // TODO: replace this
}

func (u *sessionsRepo) DeleteSession(token string) error {

	result := u.db.Where("token = ?", token).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *sessionsRepo) UpdateSessions(session models.Session) error {
	result := u.db.Model(&session).Where("email = ?", session.Email).Updates(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *sessionsRepo) SessionAvailEmail(email string) (models.Session, error) {
	var session models.Session
	result := u.db.Where("email = ?", email).First(&session)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.Session{}, result.Error
		}
		return models.Session{}, result.Error
	}
	return session, nil // TODO: replace this
}

func (u *sessionsRepo) SessionAvailToken(token string) (models.Session, error) {
	var session models.Session
	result := u.db.Where("token = ?", token).First(&session)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.Session{}, result.Error
		}
		return models.Session{}, result.Error
	}
	return session, nil// TODO: replace this
}

func (u *sessionsRepo) TokenValidity(token string) (models.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return models.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSession(token)
		if err != nil {
			return models.Session{}, err
		}
		return models.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session models.Session) bool {
	return session.Expiry.Before(time.Now())
}

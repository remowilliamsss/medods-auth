package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	dao *SessionDAO
	jwt *JWT
}

func NewService(dao *SessionDAO, jwt *JWT) *Service {
	return &Service{dao: dao, jwt: jwt}
}

func (s *Service) Auth(ctx context.Context, uuid *uuid.UUID) (*TokenPair, error) {
	_, err := s.dao.FindById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.jwt.GenerateAccessToken(uuid)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRefreshToken(accessToken)
	if err != nil {
		return nil, err
	}
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 14)
	if err != nil {
		return nil, err
	}
	session := &Session{
		UUID:         uuid,
		RefreshToken: hashedRefreshToken,
		ExpiresAt:    time.Now().Add(14 * 24 * time.Hour)}
	err = s.dao.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return &TokenPair{accessToken, refreshToken}, nil
}

func (s *Service) Refresh(ctx context.Context, accessToken string, refreshToken []byte) (*TokenPair, error) {
	userId, err := s.jwt.ParseUUID(accessToken)
	if err != nil {
		return nil, err
	}
	session, err := s.dao.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(session.RefreshToken, refreshToken)
	if err != nil {
		return nil, ErrWrongToken
	}
	return s.Auth(ctx, userId)
}

func generateRefreshToken(accessToken string) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b) + accessToken[len(accessToken)-7:], nil
}

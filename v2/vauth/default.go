package vauth

import "github.com/google/uuid"

type AuthToken struct {
	Uid   string
	Value any
}

func (t *AuthToken) GetKey() string {
	return t.Uid
}

func NewAuthToken(data any) Token {
	return &AuthToken{
		Uid:   uuid.NewString(),
		Value: data,
	}
}

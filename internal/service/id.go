package service

import "github.com/google/uuid"

type ID struct{}

func NewID() *ID {
	return &ID{}
}

func (i *ID) GenerateID() string {
	return uuid.NewString()
}

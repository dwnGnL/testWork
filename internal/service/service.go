package service

import (
	"context"

	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/internal/models"
	"github.com/dwnGnL/testWork/internal/service/auth"
	"github.com/dwnGnL/testWork/internal/service/book"
)

type ServiceImpl struct {
	conf *config.Config
	auth Auther
	book BookService
}

type Auther interface {
	Login(ctx context.Context, username string, password string) (string, error)
	CheckToken(tokenStr string) (int64, error)
}

type BookService interface {
	AllBook(ctx context.Context) []*models.Book
}

type Option func(*ServiceImpl)

func New(conf *config.Config, opts ...Option) *ServiceImpl {

	s := ServiceImpl{
		conf: conf,
		auth: auth.New(conf),
		book: book.New(conf),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s ServiceImpl) GetAuth() Auther {
	return s.auth
}

func (s ServiceImpl) GetBook() BookService {
	return s.book
}

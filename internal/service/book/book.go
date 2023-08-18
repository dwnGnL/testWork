package book

import (
	"context"
	"time"

	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/lib/bookClient"
	"github.com/dwnGnL/testWork/lib/goerrors"

	"github.com/dwnGnL/testWork/internal/models"
)

type book struct {
	client *bookClient.Client
}

func New(conf *config.Config) *book {
	client, err := bookClient.New(&bookClient.Config{Host: conf.BookClient.Host, ConnectionTimeout: 10 * time.Second})
	if err != nil {
		goerrors.Log().WithError(err).Warn("book client connect err")
	}
	return &book{client: client}
}

func (book) AllBook(ctx context.Context) []*models.Book {
	return []*models.Book{
		{
			Name:   "Book1",
			Author: "Author1",
		},
		{
			Name:   "Book2",
			Author: "Author2",
		},
		{
			Name:   "Book3",
			Author: "Author3",
		},
	}
}

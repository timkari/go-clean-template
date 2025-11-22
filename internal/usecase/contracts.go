// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) (entity.TranslationHistory, error)
	}

	// Comment -.
	Comment interface {
		Create(context.Context, entity.Comment) (entity.Comment, error)
		GetByID(context.Context, int64) (entity.Comment, error)
		GetList(context.Context, entity.CommentFilter) (entity.CommentList, error)
		Update(context.Context, entity.Comment) (entity.Comment, error)
	}
)

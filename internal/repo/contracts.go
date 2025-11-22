// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}
)

type (
	// CommentRepo -.
	CommentRepo interface {
		Create(ctx context.Context, comment entity.Comment) (entity.Comment, error)
		GetByID(ctx context.Context, id int64) (entity.Comment, error)
		GetList(ctx context.Context, filter entity.CommentFilter) (entity.CommentList, error)
		Update(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	}

	// CommentWebAPI - placeholder for external services if needed.
	CommentWebAPI interface {
		// Methods for external comment services can be added here
	}
)

package comment

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo   repo.CommentRepo
	logger logger.Interface
}

// New -.
func New(r repo.CommentRepo, l logger.Interface) *UseCase {
	return &UseCase{
		repo:   r,
		logger: l,
	}
}

// Create - создает новый комментарий.
func (uc *UseCase) Create(ctx context.Context, comment entity.Comment) (entity.Comment, error) {

	if comment.EntityType == "" {
		return entity.Comment{}, fmt.Errorf("entity type is required")
	}
	if comment.EntityID == 0 {
		return entity.Comment{}, fmt.Errorf("entity ID is required")
	}
	if comment.UserID == 0 {
		return entity.Comment{}, fmt.Errorf("user ID is required")
	}
	if comment.UserName == "" {
		return entity.Comment{}, fmt.Errorf("user name is required")
	}
	if comment.Content == "" {
		return entity.Comment{}, fmt.Errorf("content is required")
	}

	createdComment, err := uc.repo.Create(ctx, comment)
	if err != nil {
		uc.logger.Error(err, "CommentUseCase - Create - uc.repo.Create")
		return entity.Comment{}, fmt.Errorf("failed to create comment: %w", err)
	}

	uc.logger.Info("Comment created",
		"id", createdComment.ID,
		"entity_type", createdComment.EntityType,
		"entity_id", createdComment.EntityID,
		"user_id", createdComment.UserID,
	)

	return createdComment, nil
}

// GetByID - получает комментарий по ID.
func (uc *UseCase) GetByID(ctx context.Context, id int64) (entity.Comment, error) {
	if id == 0 {
		return entity.Comment{}, fmt.Errorf("comment ID is required")
	}

	comment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error(err, "CommentUseCase - GetByID - uc.repo.GetByID", "id", id)
		return entity.Comment{}, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

// GetList - получает список комментариев с пагинацией.
func (uc *UseCase) GetList(ctx context.Context, filter entity.CommentFilter) (entity.CommentList, error) {

	if filter.EntityType == "" {
		return entity.CommentList{}, fmt.Errorf("entity type is required")
	}
	if filter.EntityID == 0 {
		return entity.CommentList{}, fmt.Errorf("entity ID is required")
	}

	comments, err := uc.repo.GetList(ctx, filter)
	if err != nil {
		uc.logger.Error(err, "CommentUseCase - GetList - uc.repo.GetList",
			"entity_type", filter.EntityType,
			"entity_id", filter.EntityID,
		)
		return entity.CommentList{}, fmt.Errorf("failed to get comments: %w", err)
	}

	uc.logger.Info("Comments retrieved",
		"entity_type", filter.EntityType,
		"entity_id", filter.EntityID,
		"total", comments.Total,
		"page", comments.Page,
	)

	return comments, nil
}

// Update - обновляет комментарий.
func (uc *UseCase) Update(ctx context.Context, comment entity.Comment) (entity.Comment, error) {

	if comment.ID == 0 {
		return entity.Comment{}, fmt.Errorf("comment ID is required")
	}
	if comment.Content == "" {
		return entity.Comment{}, fmt.Errorf("content is required")
	}

	existingComment, err := uc.repo.GetByID(ctx, comment.ID)
	if err != nil {
		uc.logger.Error(err, "CommentUseCase - Update - uc.repo.GetByID", "id", comment.ID)
		return entity.Comment{}, fmt.Errorf("comment not found")
	}

	existingComment.Content = comment.Content

	updatedComment, err := uc.repo.Update(ctx, existingComment)
	if err != nil {
		uc.logger.Error(err, "CommentUseCase - Update - uc.repo.Update", "id", comment.ID)
		return entity.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}

	uc.logger.Info("Comment updated",
		"id", updatedComment.ID,
		"user_id", updatedComment.UserID,
	)

	return updatedComment, nil
}

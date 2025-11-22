package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

const _defaultCommentCap = 64

// CommentRepo - структура репозитория.
type CommentRepo struct {
	*postgres.Postgres
}

// NewCommentRepo - конструктор репозитория.
func NewCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

// Create - создает новый комментарий в БД.
func (r *CommentRepo) Create(ctx context.Context, comment entity.Comment) (entity.Comment, error) {

	sql, args, err := r.Builder.
		Insert("comments").
		Columns("entity_type, entity_id, user_id, user_name, content").
		Values(comment.EntityType, comment.EntityID, comment.UserID, comment.UserName, comment.Content).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - Create - r.Builder: %w", err)
	}

	var createdComment entity.Comment
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&createdComment.ID,
		&createdComment.CreatedAt,
		&createdComment.UpdatedAt,
	)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - Create - r.Pool.QueryRow: %w", err)
	}

	createdComment.EntityType = comment.EntityType
	createdComment.EntityID = comment.EntityID
	createdComment.UserID = comment.UserID
	createdComment.UserName = comment.UserName
	createdComment.Content = comment.Content

	return createdComment, nil
}

// GetByID - получает комментарий по ID.
func (r *CommentRepo) GetByID(ctx context.Context, id int64) (entity.Comment, error) {

	sql, _, err := r.Builder.
		Select("id, entity_type, entity_id, user_id, user_name, content, created_at, updated_at").
		From("comments").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - GetByID - r.Builder: %w", err)
	}

	var comment entity.Comment
	err = r.Pool.QueryRow(ctx, sql, id).Scan(
		&comment.ID,
		&comment.EntityType,
		&comment.EntityID,
		&comment.UserID,
		&comment.UserName,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return comment, nil
}

// GetList - получает список комментариев с пагинацией и сортировкой.
func (r *CommentRepo) GetList(ctx context.Context, filter entity.CommentFilter) (entity.CommentList, error) {

	baseQuery := r.Builder.
		Select("id, entity_type, entity_id, user_id, user_name, content, created_at, updated_at").
		From("comments").
		Where("entity_type = ? AND entity_id = ?", filter.EntityType, filter.EntityID)

	countSQL, countArgs, err := baseQuery.
		RemoveLimit().
		RemoveOffset().
		Prefix("SELECT COUNT(*) FROM (").
		Suffix(") AS count_query").
		ToSql()
	if err != nil {
		return entity.CommentList{}, fmt.Errorf("CommentRepo - GetList - count query: %w", err)
	}

	var total int64
	err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return entity.CommentList{}, fmt.Errorf("CommentRepo - GetList - count scan: %w", err)
	}

	// Применяем сортировку
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == entity.SortAsc {
		sortOrder = "ASC"
	}
	baseQuery = baseQuery.OrderBy(fmt.Sprintf("%s %s", sortBy, sortOrder))

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	baseQuery = baseQuery.Limit(uint64(pageSize)).Offset(uint64(offset))

	sql, args, err := baseQuery.ToSql()
	if err != nil {
		return entity.CommentList{}, fmt.Errorf("CommentRepo - GetList - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.CommentList{}, fmt.Errorf("CommentRepo - GetList - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	comments := make([]entity.Comment, 0, _defaultCommentCap)

	for rows.Next() {
		var comment entity.Comment
		err = rows.Scan(
			&comment.ID,
			&comment.EntityType,
			&comment.EntityID,
			&comment.UserID,
			&comment.UserName,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return entity.CommentList{}, fmt.Errorf("CommentRepo - GetList - rows.Scan: %w", err)
		}
		comments = append(comments, comment)
	}

	return entity.CommentList{
		Comments: comments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Update - обновляет существующий комментарий.
func (r *CommentRepo) Update(ctx context.Context, comment entity.Comment) (entity.Comment, error) {

	sql, args, err := r.Builder.
		Update("comments").
		Set("content", comment.Content).
		Set("updated_at", "NOW()").
		Where("id = ?", comment.ID).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - Update - r.Builder: %w", err)
	}

	var updatedAt string
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&updatedAt)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - Update - r.Pool.QueryRow: %w", err)
	}

	return r.GetByID(ctx, comment.ID)
}

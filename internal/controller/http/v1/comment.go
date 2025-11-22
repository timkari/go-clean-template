package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// CreateCommentRequest -.
type CreateCommentRequest struct {
	Content    string `json:"content" validate:"required"`
	UserName   string `json:"user_name" validate:"required"`
	EntityType string `json:"entity_type" validate:"required"`
	EntityID   int64  `json:"entity_id" validate:"required"`
	UserID     int64  `json:"user_id" validate:"required"`
}

// UpdateCommentRequest -.
type UpdateCommentRequest struct {
	Content  string `json:"content" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
}

// createComment создает новый комментарий
func (r *V1) createComment(ctx *fiber.Ctx) error {
	var req CreateCommentRequest

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createComment")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Валидация
	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createComment - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	comment := entity.Comment{
		Content:    req.Content,
		UserName:   req.UserName,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		UserID:     req.UserID,
	}

	createdComment, err := r.c.Create(ctx.Context(), comment)
	if err != nil {
		r.l.Error(err, "http - v1 - createComment")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusCreated).JSON(createdComment)
}

// getCommentsList возвращает список комментариев
func (r *V1) getCommentsList(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	entityID, err := strconv.ParseInt(ctx.Query("entity_id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid entity_id")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))

	filter := entity.CommentFilter{
		EntityType: entityType,
		EntityID:   entityID,
		Page:       page,
		PageSize:   pageSize,
		SortBy:     ctx.Query("sort_by", "created_at"),
		SortOrder:  entity.SortDirection(ctx.Query("sort_order", "desc")),
	}

	comments, err := r.c.GetList(ctx.Context(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - getCommentsList")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(comments)
}

// getCommentByID возвращает комментарий по ID
func (r *V1) getCommentByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid comment ID")
	}

	comment, err := r.c.GetByID(ctx.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getCommentByID", "id", id)
		return errorResponse(ctx, http.StatusNotFound, "comment not found")
	}

	return ctx.Status(http.StatusOK).JSON(comment)
}

// updateComment обновляет комментарий
func (r *V1) updateComment(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid comment ID")
	}

	var req UpdateCommentRequest
	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updateComment")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Валидация
	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updateComment - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	// Получаем существующий комментарий
	existingComment, err := r.c.GetByID(ctx.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - updateComment - get existing")
		return errorResponse(ctx, http.StatusNotFound, "comment not found")
	}

	// Обновляем поля
	existingComment.Content = req.Content
	existingComment.UserName = req.UserName

	updatedComment, err := r.c.Update(ctx.Context(), existingComment)
	if err != nil {
		r.l.Error(err, "http - v1 - updateComment")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(updatedComment)
}

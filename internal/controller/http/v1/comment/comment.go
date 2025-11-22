package comment

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type commentRoutes struct {
	t usecase.Comment
	l logger.Interface
}

func newCommentRoutes(handler fiber.Router, t usecase.Comment, l logger.Interface) {
	r := &commentRoutes{t, l}

	handler.Post("/comments", r.create)
	handler.Get("/comments", r.getList)
	handler.Get("/comments/:id", r.getByID)
	handler.Put("/comments/:id", r.update)
}

// @Summary     Create comment
// @Description Create a new comment
// @ID          create-comment
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       request body entity.Comment true "Comment data"
// @Success     201 {object} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments [post]
func (r *commentRoutes) create(ctx *fiber.Ctx) error {
	var comment entity.Comment

	if err := ctx.BodyParser(&comment); err != nil {
		r.l.Error(err, "http - v1 - comment - create")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	createdComment, err := r.t.Create(ctx.UserContext(), comment)
	if err != nil {
		r.l.Error(err, "http - v1 - comment - create")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusCreated).JSON(createdComment)
}

// @Summary     Get comments list
// @Description Get paginated list of comments for entity
// @ID          get-comments-list
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       entity_type query string true "Entity type"
// @Param       entity_id   query int    true "Entity ID"
// @Param       page        query int    false "Page number" default(1)
// @Param       page_size   query int    false "Page size" default(10)
// @Param       sort_by     query string false "Sort field" default(created_at)
// @Param       sort_order  query string false "Sort order" Enums(asc, desc) default(desc)
// @Success     200 {object} entity.CommentList
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments [get]
func (r *commentRoutes) getList(ctx *fiber.Ctx) error {
	filter := entity.CommentFilter{
		EntityType: ctx.Query("entity_type"),
		Page:       getIntQuery(ctx, "page", 1),
		PageSize:   getIntQuery(ctx, "page_size", 10),
		SortBy:     ctx.Query("sort_by", "created_at"),
		SortOrder:  entity.SortDirection(ctx.Query("sort_order", "desc")),
	}

	entityID, err := strconv.ParseInt(ctx.Query("entity_id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid entity_id")
	}
	filter.EntityID = entityID

	comments, err := r.t.GetList(ctx.UserContext(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - comment - getList")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(comments)
}

// @Summary     Get comment by ID
// @Description Get comment by ID
// @ID          get-comment-by-id
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       id path int true "Comment ID"
// @Success     200 {object} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments/{id} [get]
func (r *commentRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid comment ID")
	}

	comment, err := r.t.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - comment - getByID", "id", id)
		return errorResponse(ctx, http.StatusNotFound, "comment not found")
	}

	return ctx.Status(http.StatusOK).JSON(comment)
}

// @Summary     Update comment
// @Description Update comment content
// @ID          update-comment
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       id      path int             true "Comment ID"
// @Param       request body entity.Comment true "Comment data"
// @Success     200 {object} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments/{id} [put]
func (r *commentRoutes) update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid comment ID")
	}

	var comment entity.Comment
	if err := ctx.BodyParser(&comment); err != nil {
		r.l.Error(err, "http - v1 - comment - update")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	comment.ID = id

	updatedComment, err := r.t.Update(ctx.UserContext(), comment)
	if err != nil {
		r.l.Error(err, "http - v1 - comment - update", "id", id)
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(updatedComment)
}

// Вспомогательные функции
func getIntQuery(ctx *fiber.Ctx, key string, defaultValue int) int {
	if value := ctx.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func errorResponse(ctx *fiber.Ctx, code int, message string) error {
	return ctx.Status(code).JSON(map[string]string{
		"error": message,
	})
}

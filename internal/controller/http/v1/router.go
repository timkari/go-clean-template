package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group fiber.Router, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}
}

// NewCommentRoutes -.
func NewCommentRoutes(apiV1Group fiber.Router, c usecase.Comment, l logger.Interface) {
	r := &V1{c: c, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	commentGroup := apiV1Group.Group("/comments")

	{
		commentGroup.Post("/", r.createComment)
		commentGroup.Get("/", r.getCommentsList)
		commentGroup.Get("/:id", r.getCommentByID)
		commentGroup.Put("/:id", r.updateComment)
	}
}

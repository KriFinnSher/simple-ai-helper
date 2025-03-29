package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"support/internal/usecase"
)

type Knowledge struct {
	UseCase *usecase.KnowledgeUseCase
}

func NewKnowledgeInstance(uu *usecase.KnowledgeUseCase) *Knowledge {
	return &Knowledge{UseCase: uu}
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Answer string `json:"answer"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func (h *Knowledge) Answer(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Errors: "bad request body"})
	}

	intent, err := h.UseCase.ExtractIntent(ctx.Request().Context(), req.Message)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Errors: "failed to define intent"})
	}

	answer, err := h.UseCase.GetAnswer(ctx.Request().Context(), intent)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Errors: "failed to fetch answer"})
	}

	return ctx.JSON(http.StatusOK, Response{Answer: answer})
}

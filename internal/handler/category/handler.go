package category

import (
	"errors"
	"strconv"

	categorydto "sys-admin-serve/internal/dto/category"
	"sys-admin-serve/internal/response"
	servicecategory "sys-admin-serve/internal/service/category"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	categoryService *servicecategory.Service
}

func NewHandler(categoryService *servicecategory.Service) *Handler {
	return &Handler{categoryService: categoryService}
}

func (h *Handler) List(c *gin.Context) {
	var req categorydto.ListCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid category query params")
		return
	}

	result, err := h.categoryService.ListCategories(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, servicecategory.ErrCategoryInvalidRequest):
			response.BadRequest(c, "invalid category query params")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Page(c, result.List, result.Total, result.Page, result.PageSize)
}

func (h *Handler) Tree(c *gin.Context) {
	var req categorydto.CategoryTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid category tree query params")
		return
	}

	result, err := h.categoryService.GetCategoryTree(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, servicecategory.ErrCategoryParentNotFound):
			response.BadRequest(c, "parent category not found")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req categorydto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid create category request")
		return
	}

	result, err := h.categoryService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, servicecategory.ErrCategoryInvalidRequest):
			response.BadRequest(c, "invalid create category request")
		case errors.Is(err, servicecategory.ErrCategoryParentNotFound):
			response.BadRequest(c, "parent category not found")
		case errors.Is(err, servicecategory.ErrCategoryInvalidParent):
			response.BadRequest(c, "invalid parent category")
		case errors.Is(err, servicecategory.ErrCategoryCodeExists):
			response.BadRequest(c, "category code already exists")
		case errors.Is(err, servicecategory.ErrCategoryNameExists):
			response.BadRequest(c, "category name already exists")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	categoryID, err := parseUint64ID(c.Param("id"))
	if err != nil || categoryID == 0 {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req categorydto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid update category request")
		return
	}

	result, err := h.categoryService.UpdateCategory(c.Request.Context(), categoryID, req)
	if err != nil {
		switch {
		case errors.Is(err, servicecategory.ErrCategoryNotFound):
			response.NotFound(c, "category not found")
		case errors.Is(err, servicecategory.ErrCategoryInvalidRequest):
			response.BadRequest(c, "invalid update category request")
		case errors.Is(err, servicecategory.ErrCategoryParentNotFound):
			response.BadRequest(c, "parent category not found")
		case errors.Is(err, servicecategory.ErrCategoryInvalidParent):
			response.BadRequest(c, "invalid parent category")
		case errors.Is(err, servicecategory.ErrCategoryCodeExists):
			response.BadRequest(c, "category code already exists")
		case errors.Is(err, servicecategory.ErrCategoryNameExists):
			response.BadRequest(c, "category name already exists")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, result)
}

func parseUint64ID(raw string) (uint64, error) {
	return strconv.ParseUint(raw, 10, 64)
}

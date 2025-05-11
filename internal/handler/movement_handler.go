package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gopkg.in/go-playground/validator.v9"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/util"
)

// MovementHandler handles movement-related requests
type MovementHandler struct {
	movementService service.MovementService
	accountService  service.AccountService
	validator       *validator.Validate
}

// NewMovementHandler creates a new movement handler
func NewMovementHandler(
	movementService service.MovementService,
	accountService service.AccountService,
) *MovementHandler {
	return &MovementHandler{
		movementService: movementService,
		accountService:  accountService,
		validator:       validator.New(),
	}
}

// CreateMovementRequest represents a request to create a new movement
type CreateMovementRequest struct {
	Amount      string `json:"amount" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=credit debit"`
	Description string `json:"description"`
}

// List returns a paginated list of movements for the user's account
// @Summary List account movements
// @Description Get a paginated list of account movements
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} util.PaginatedResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /accounts/movements [get]
func (h *MovementHandler) List(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{
			Error: util.NewUnauthorizedError("unauthorized"),
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{
			Error: util.NewInternalServerError("user context invalid"),
		})
		return
	}

	// Get account
	account, err := h.accountService.GetByUserID(c, userModel.ID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Get movements
	response, err := h.movementService.GetByAccountID(c, account.ID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// Create creates a new movement for the user's account
// @Summary Create account movement
// @Description Create a new movement (credit or debit) in the account
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param movement body CreateMovementRequest true "Movement details"
// @Success 201 {object} model.Movement
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /accounts/movements [post]
func (h *MovementHandler) Create(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse{
			Error: util.NewUnauthorizedError("unauthorized"),
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse{
			Error: util.NewInternalServerError("user context invalid"),
		})
		return
	}

	// Parse and validate request
	var req CreateMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("invalid request body"),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError(err.Error()),
		})
		return
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("invalid amount"),
		})
		return
	}

	// Get account
	account, err := h.accountService.GetByUserID(c, userModel.ID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Create movement
	movement, err := h.movementService.Create(
		c,
		account.ID,
		amount,
		req.Type,
		req.Description,
	)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusCreated, movement)
}

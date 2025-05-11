package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/util"
)

// TransferHandler handles transfer-related requests
type TransferHandler struct {
	transferService service.TransferService
	accountService  service.AccountService
	validator       *validator.Validate
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(
	transferService service.TransferService,
	accountService service.AccountService,
) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
		accountService:  accountService,
		validator:       validator.New(),
	}
}

// TransferRequest represents a request to create a new transfer
type TransferRequest struct {
	ToAccount   string `json:"to_account" validate:"required,uuid4"`
	Amount      string `json:"amount" validate:"required"`
	Description string `json:"description"`
}

// Transfer performs a transfer from the user's account to another account
// @Summary Create a transfer
// @Description Transfer funds from the authenticated user's account to another account
// @Tags transfers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transfer body TransferRequest true "Transfer details"
// @Success 201 {object} model.Transfer
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /transfers [post]
func (h *TransferHandler) Transfer(c *gin.Context) {
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
	var req TransferRequest
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

	// Parse to account ID
	toAccountID, err := uuid.Parse(req.ToAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("invalid to_account"),
		})
		return
	}

	// Get from account (user's account)
	fromAccount, err := h.accountService.GetByUserID(c, userModel.ID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Perform transfer
	transfer, err := h.transferService.Transfer(
		c,
		fromAccount.ID,
		toAccountID,
		amount,
		req.Description,
	)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusCreated, transfer)
}

// List returns a paginated list of transfers for the user's account
// @Summary List transfers
// @Description Get a paginated list of transfers for the authenticated user's account
// @Tags transfers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} util.PaginatedResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /transfers [get]
func (h *TransferHandler) List(c *gin.Context) {
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

	// Get transfers
	response, err := h.transferService.GetByAccountID(c, account.ID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

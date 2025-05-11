package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/util"
)

// AccountHandler handles account-related requests
type AccountHandler struct {
	accountService service.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// BalanceResponse represents an account balance response
type BalanceResponse struct {
	AccountID uuid.UUID `json:"account_id"`
	Balance   string    `json:"balance"`
	Currency  string    `json:"currency"`
}

// Balance gets the balance of the authenticated user's account
// @Summary Get account balance
// @Description Get the current balance of the authenticated user's account
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} BalanceResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /accounts/balance [get]
func (h *AccountHandler) Balance(c *gin.Context) {
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

	// Get balance (uses cache when available)
	balance, err := h.accountService.GetBalance(c, account.ID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Create response
	response := BalanceResponse{
		AccountID: account.ID,
		Balance:   balance.String(),
		Currency:  account.Currency,
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/saalikmubeen/backend-masterclass-go/authToken"
	generated_db "github.com/saalikmubeen/backend-masterclass-go/db/sqlc"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,len=3,oneof=USD EUR CAD"`
	// Balance  int64  `json:"balance" binding:"required,gt=0"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var json createAccountRequest

	// ShouldBindJSON: To convert or bind the request JSON body
	// to the createAccountRequest struct
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the owner from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*authToken.Payload)

	arg := generated_db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: json.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var params getAccountRequest

	// ShouldBindUri: To get the URI parameter(:id) from the request
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, params.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the logged in user from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*authToken.Payload)

	// Check if the account belongs to the logged in user
	if account.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUserAccounts(ctx *gin.Context) {
	var query listAccountRequest

	// ShouldBindQuery: To get the query parameters from the request
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the logged in user from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*authToken.Payload)

	arg := generated_db.ListUserAccountsParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
		Owner:  authPayload.Username,
	}

	accounts, err := server.store.ListUserAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) listAllAccounts(ctx *gin.Context) {
	var query listAccountRequest

	// ShouldBindQuery: To get the query parameters from the request
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := generated_db.ListAllAccountsParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	}

	accounts, err := server.store.ListAllAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

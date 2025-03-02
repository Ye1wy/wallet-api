package controller

import (
	"log/slog"
	"net/http"
	"wallet-api/internal/config"
	"wallet-api/internal/dto"
	"wallet-api/internal/logger"
	"wallet-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type WalletController struct {
	service service.WalletService
	log     *slog.Logger
}

func NewWalletController(service service.WalletService, logger *slog.Logger) *WalletController {
	return &WalletController{
		service: service,
		log:     logger,
	}
}

// GetWalletById godoc
//
//	@Summary		Get Wallat By Id
//	@Description	Get wallet data by uuid
//	@Tags			wallets
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uuid.UUID true "Wallet ID"
//	@Success		200	{object}	dto.WalletDTO
//	@Failure		400	{object}	dto.ErrorDTO
//	@Failure		404	{object}	dto.ErrorDTO
//	@Failure		500	{object}	dto.ErrorDTO
//	@Router			/api/v1/wallets/{id} [get]
func (ctrl *WalletController) GetWalletById(ctx *gin.Context) {
	op := "controller.wallet-controller.GetWalletById"

	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		ctrl.log.Error("Bad request: invalid id", logger.Err(err), "op", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctrl.log.Info("Start validate wallet id from param", "op", op)

	dto, err := ctrl.service.GetWalletById(ctx, id)

	if err == pgx.ErrNoRows {
		ctrl.log.Warn("Wallet not found with id", "id", id, "op", op)
		ctx.JSON(http.StatusNotFound, gin.H{"massage": "Wallet with that uuid is not found"})
		return
	}

	if err != nil {
		ctrl.log.Error("Error to get wallet by id", logger.Err(err), "op", op)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctrl.log.Info("Wallet data retrived successfully", "op", op)
	ctx.JSON(http.StatusOK, dto)
}

// OperationWithWalletById godoc
//
//	@Summary		Operation with wallet by id
//	@Description	Change wallet data with operation deposit and withdraw
//	@Tags			tasks
//	@Accept			json
//	@Produce		dto.ErrorDTO
//	@Success		200	{object}	dto.ErrorDTO
//	@Failure		400	{object}	dto.ErrorDTO
//	@Failure		422	{object}	dto.ErrorDTO
//	@Failure		500	{object}	dto.ErrorDTO
//	@Router			/api/v1/wallet [post]
func (ctrl *WalletController) OperationWithWalletByID(ctx *gin.Context) {
	op := "controller.wallet-controller.OperationWithWalletByID"

	var dto dto.WalletOperationRequestDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctrl.log.Error("Failed to bind JSON to DTO", logger.Err(err), "op", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.service.OperationWithWalletByID(ctx, dto)

	if err != nil {
		if err.Error() == config.InvalidOperation || err.Error() == config.AmountIsNotValid {
			ctrl.log.Warn("One of the parameters contains unprocessable entity", logger.Err(err), "op", op)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
			return
		}

		ctrl.log.Error("Error to update wallet by id", logger.Err(err), "op", op)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctrl.log.Info("Wallet data retrived successfully", "op", op)
	ctx.JSON(http.StatusOK, gin.H{"massage": "balance is updated"})
}

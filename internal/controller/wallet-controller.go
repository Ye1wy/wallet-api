package controller

import (
	"log/slog"
	"net/http"
	"wallet-api/internal/service"

	"github.com/gin-gonic/gin"
)

type walletController struct {
	service service.WalletService
	log     *slog.Logger
}

func NewWalletController(service service.WalletService, logger *slog.Logger) *walletController {
	return &walletController{
		service: service,
		log:     logger,
	}
}

func (ctrl *walletController) GetWalletById(ctx *gin.Context) {
	op := "controller.wallet-controller.GetWalletById"

	ctrl.log.Info("Wallet data retrived successfully", "op", op)
	ctx.JSON(http.StatusOK, nil)
}

func (ctrl *walletController) DepositToWalletByID(ctx *gin.Context) {
	op := "controller.wallet-controller.DepositToWalletByID"

	ctrl.log.Info("Wallet data retrived successfully", "op", op)
	ctx.JSON(http.StatusOK, nil)
}

package routes

import (
	"wallet-api/internal/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type router struct {
	router *gin.Engine
}

type RouterConfig struct {
	WalletController *controller.WalletController
}

func NewRouter(cfg RouterConfig) router {
	r := router{
		router: gin.Default(),
	}

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.router.GET("/api/v1/wallets/:id", cfg.WalletController.GetWalletById)
	r.router.POST("/api/v1/wallet", cfg.WalletController.OperationWithWalletByID)

	return r
}

func (r *router) Run(addr ...string) error {
	return r.router.Run(addr...)
}

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/silasrm/api-obm/internal/interface/http/handler"
	"github.com/silasrm/api-obm/internal/interface/http/middleware"

	_ "github.com/silasrm/api-obm/docs"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	searchHandler *handler.SearchHandler,
	vmpHandler *handler.VMPHandler,
	ampHandler *handler.AMPHandler,
	supplierHandler *handler.SupplierHandler,
	domainHandler *handler.DomainHandler,
	genericHandler *handler.GenericHandler,
	adminHandler *handler.AdminHandler,
	cmedHandler *handler.CMEDHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.POST("/auth/login", authHandler.Login)
	r.GET("/health", adminHandler.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/search", searchHandler.Search)

		api.GET("/vmp", vmpHandler.List)
		api.GET("/vmp/:id", vmpHandler.GetByID)
		api.GET("/vmp/:id/detail", vmpHandler.GetDetail)

		api.GET("/amp", ampHandler.List)
		api.GET("/amp/:id", ampHandler.GetByID)
		api.GET("/amp/:id/detail", ampHandler.GetDetail)

		api.GET("/vtm", genericHandler.ListVTM)
		api.GET("/vtm/:id", genericHandler.GetVTM)

		api.GET("/vmpp", genericHandler.ListVMPP)
		api.GET("/vmpp/:id", genericHandler.GetVMPP)

		api.GET("/ampp", genericHandler.ListAMPP)
		api.GET("/ampp/:id", genericHandler.GetAMPP)
		api.GET("/ampp/:id/cmed", cmedHandler.GetAMPPWithCMED)

		api.GET("/suppliers", supplierHandler.List)
		api.GET("/suppliers/:id", supplierHandler.GetByID)

		api.GET("/dcb", genericHandler.ListDCB)
		api.GET("/dcb/:id", genericHandler.GetDCB)

		api.GET("/ingredients", genericHandler.ListIngredients)
		api.GET("/ingredients/:id", genericHandler.GetIngredient)

		api.GET("/domains/:domain", domainHandler.List)
		api.GET("/domains/:domain/:id", domainHandler.GetByID)

		api.POST("/admin/reindex", adminHandler.Reindex)

		cmed := api.Group("/cmed")
		{
			cmed.GET("", cmedHandler.List)
			cmed.GET("/registro/:registro", cmedHandler.GetByRegistro)
			cmed.GET("/ean/:ean", cmedHandler.GetByEAN)
			cmed.GET("/:id/historico", cmedHandler.GetHistorico)
			cmed.GET("/:id", cmedHandler.GetByID)
		}
	}

	return r
}

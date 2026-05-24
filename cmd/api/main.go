// @title OBM API
// @version 1.0
// @description API do Observatorio de Medicamentos - Brasil (OBM). Gerencia dados de medicamentos seguindo o padrao dm%d adaptado para o Brasil.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Tipo: Bearer <token>. Insira o token JWT obtido via /auth/login
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	meilisearchrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/meilisearch"
	"github.com/silasrm/api-obm/internal/infrastructure/persistence/postgres"
	"github.com/silasrm/api-obm/internal/interface/http/handler"
	"github.com/silasrm/api-obm/internal/interface/http/router"
	"github.com/silasrm/api-obm/internal/usecase"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	var pgPool *pgxpool.Pool = pool

	meiliClient := meilisearch.New(cfg.Meilisearch.URL, meilisearch.WithAPIKey(cfg.Meilisearch.APIKey))

	meiliRepo := meilisearchrepo.NewMeilisearchRepo(cfg.Meilisearch)

	var userRepo repository.UserRepository = postgres.NewUserRepo(pool)
	var vmpRepo repository.VMPRepository = postgres.NewVMPRepo(pool)
	var ampRepo repository.AMPRepository = postgres.NewAMPRepo(pool)
	var supplierRepo repository.SupplierRepository = postgres.NewSupplierRepo(pool)
	var domainRepo repository.DomainRepository = postgres.NewDomainRepo(pool)
	var vtmRepo repository.VTMRepository = postgres.NewVTMRepo(pool)
	var vmppRepo repository.VMPPRepository = postgres.NewVMPPRepo(pool)
	var amppRepo repository.AMPPRepository = postgres.NewAMPPRepo(pool)
	var dcbRepo repository.DCBRepository = postgres.NewDCBRepo(pool)
	var ingredientRepo repository.IngredientSubstanceRepository = postgres.NewIngredientRepo(pool)
	var syncRepo repository.SyncRepository = postgres.NewSyncRepo(pool)

	authUseCase := usecase.NewAuthUsecase(userRepo, cfg.JWT)
	searchUseCase := usecase.NewSearchUsecase(meiliRepo)
	vmpUseCase := usecase.NewVMPUsecase(vmpRepo)
	ampUseCase := usecase.NewAMPUsecase(ampRepo)
	supplierUseCase := usecase.NewSupplierUsecase(supplierRepo)
	domainUseCase := usecase.NewDomainUsecase(domainRepo)
	genericUseCase := usecase.NewGenericUsecase(vtmRepo, vmppRepo, amppRepo, dcbRepo, ingredientRepo)
	reindexUseCase := usecase.NewReindexUsecase(syncRepo, meiliRepo)

	authHandler := handler.NewAuthHandler(authUseCase)
	searchHandler := handler.NewSearchHandler(searchUseCase)
	vmpHandler := handler.NewVMPHandler(vmpUseCase)
	ampHandler := handler.NewAMPHandler(ampUseCase)
	supplierHandler := handler.NewSupplierHandler(supplierUseCase)
	domainHandler := handler.NewDomainHandler(domainUseCase)
	genericHandler := handler.NewGenericHandler(genericUseCase)
	adminHandler := handler.NewAdminHandler(reindexUseCase, pgPool, meiliClient)

	if cfg.Sync.OnStartup {
		log.Println("Sync on startup enabled, running reindex...")
		if _, err := reindexUseCase.Reindex(context.Background()); err != nil {
			log.Printf("Warning: reindex on startup failed: %v", err)
		}
	}

	r := router.SetupRouter(
		authHandler,
		searchHandler,
		vmpHandler,
		ampHandler,
		supplierHandler,
		domainHandler,
		genericHandler,
		adminHandler,
		cfg.JWT.Secret,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

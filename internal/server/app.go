package server

import (
	"context"
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/middleware"
	"infotecstechtask/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	dhttp "infotecstechtask/internal/delivery/http"
	prepo "infotecstechtask/internal/payment/repository"
	trepo "infotecstechtask/internal/transaction/repository"
	tservice "infotecstechtask/internal/transaction/service"
	wrepo "infotecstechtask/internal/wallet/repository"
	wservice "infotecstechtask/internal/wallet/service"
)

// Структура, хранящая в себе указатель на http сервер и экземпляр фасада
type App struct {
	httpServer *http.Server

	facade facade.TransactionFacade
}

func NewApp() *App {
	config := database.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbClient, err := database.NewClient(ctx, config)
	if err != nil {
		log.Fatal("Failed to create DB client: %w", err)
	}

	walletRepository := wrepo.NewWalletRepository(dbClient)
	transactionRepository := trepo.NewTransactionRepository(dbClient)
	paymentRepository := prepo.NewPaymentRepository(dbClient)

	walletService := wservice.NewWalletService(walletRepository)
	transactionService := tservice.NewTransactionService(transactionRepository)

	return &App{
		facade: *facade.NewFacade(walletService, transactionService, paymentRepository),
	}
}

// Функция для запуска http сервера
// Реализован базовый механизм graceful shutdown
func (a App) Run(port string) error {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	limiterInstance := limiter.New(store, rate)

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.RateLimiter(limiterInstance),
	)

	validate := validator.New()

	dhttp.RegisterHTTPEndpoints(&router.RouterGroup, a.facade, validate)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

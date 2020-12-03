package bootstrap

import (
	"nyala-backend/pkg/logruslogger"
	api "nyala-backend/server/handler"
	"nyala-backend/server/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RegisterRoutes ...
func (boot *Bootup) RegisterRoutes() {
	handlerType := api.Handler{
		DB:         boot.DB,
		EnvConfig:  boot.EnvConfig,
		Validate:   boot.Validator,
		Translator: boot.Translator,
		ContractUC: &boot.ContractUC,
		Jwe:        boot.Jwe,
		Jwt:        boot.Jwt,
	}
	mJwt := middleware.VerifyMiddlewareInit{
		ContractUC: &boot.ContractUC,
	}

	boot.R.Route("/v1", func(r chi.Router) {
		// Define a limit rate to 1000 requests per IP per request.
		rate, _ := limiter.NewRateFromFormatted("1000-S")
		store, _ := sredis.NewStoreWithOptions(boot.ContractUC.Redis, limiter.StoreOptions{
			Prefix:   "limiter_global",
			MaxRetry: 3,
		})
		rateMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
		r.Use(rateMiddleware.Handler)

		// Logging setup
		r.Use(chimiddleware.RequestID)
		r.Use(logruslogger.NewStructuredLogger(boot.EnvConfig["LOG_FILE_PATH"], boot.EnvConfig["LOG_DEFAULT"], boot.ContractUC.ReqID))
		r.Use(chimiddleware.Recoverer)

		// API
		r.Route("/api", func(r chi.Router) {
			logicHandler := api.LogicHandler{Handler: handlerType}
			r.Route("/logic", func(r chi.Router) {
				r.Get("/fibonacci", logicHandler.GetFibonacciHandler)
			})

			customerHandler := api.CustomerHandler{Handler: handlerType}
			r.Route("/customer", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Post("/register", customerHandler.RegisterHandler)
					r.Post("/login", customerHandler.LoginHandler)
				})

				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyRefreshTokenCredential)
					r.Get("/refresh_token", customerHandler.RefreshTokenHandler)
				})
			})

			orderHandler := api.OrderHandler{Handler: handlerType}
			r.Route("/order", func(r chi.Router) {
				r.Use(mJwt.VerifyCustomerTokenCredential)
				r.Post("/", orderHandler.CreateHandler)
			})
		})
	})
}

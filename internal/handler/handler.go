package handler

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/service"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	productH := newProductHandler(h.service.Product)
	customerH := newCustomerHandler(h.service.Customer)
	orderH := newOrderHandler(h.service.Order)

	r.Use(middleware.RedirectSlashes)

	r.Post("/v1/staff/register", userH.Register)
	r.Post("/v1/staff/login", userH.Login)
	r.Get("/v1/product/customer", productH.GetProductShop)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Patch("/v1/user", userH.UpdateAccount)

		r.Post("/v1/product", productH.AddProduct)
		r.Get("/v1/product", productH.GetProduct)
		r.Put("/v1/product/{id}", productH.UpdateProduct)
		r.Delete("/v1/product/{id}", productH.DeleteProduct)

		r.Post("/v1/customer/register", customerH.RegisterCustomer)
		r.Get("/v1/customer", customerH.GetCustomer)

		r.Post("/v1/product/checkout", orderH.Checkout)
	})
}

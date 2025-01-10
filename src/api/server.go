package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shyam0507/pd-payment/src/storage"
)

type Server struct {
	port     string
	storage  storage.Storage
	producer storage.Producer
	r        *echo.Echo
}

func NewServer(port string, storage storage.Storage, producer storage.Producer) *Server {
	return &Server{port: port, storage: storage, producer: producer, r: echo.New()}
}

func (s *Server) Start() {
	e := s.r
	e.Use(middleware.RequestID())

	g := e.Group("/api/payment/v1.0")
	g.POST("/orderId/:orderId", s.updatePayment)

	e.Start(":" + s.port)
}

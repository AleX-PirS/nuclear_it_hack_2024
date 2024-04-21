package http

import (
	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/http/dto"
	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/http/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Config struct {
	host string
	port string
}

var config *Config 

func init() {
	config = &Config{
		host: "127.0.0.1",
		port: "8888",
	}
} 

type Server struct {
	app *fiber.App
	h *handlers.Handler
	config *Config
}

func (s *Server) GetChans() (chan dto.Request, chan dto.Response) {
	return s.h.GetChans()
}

func NewServer(f *fiber.App, h *handlers.Handler) *Server{
	return &Server{
		app: f,
		config: config,
		h: h,
	}
}

func New() *fiber.App{
	return fiber.New()
}

func (s *Server) Register() {
	s.app.Post("/upload_go", s.h.HandleJsons)
	s.app.Get("/", func(c *fiber.Ctx) error {
		log.Info("New request")
		return c.SendStatus(200)
	})
}

func (s *Server) ConfigurateAndRun() {
	s.Register()
	log.Fatal(s.app.Listen(s.config.host+":"+s.config.port))
}


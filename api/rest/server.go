package rest

import (
	"github.com/janminder/content-delivery-s3-backend/api/services"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Services struct {
	fileService services.FileService
}

type Server struct {
	echo     *echo.Echo
	services Services
	conf     *viper.Viper
}

func NewServer(fileService services.FileService, c *viper.Viper) *Server {

	// Setup all Services
	services := Services{
		fileService: fileService,
	}

	// Setup Server
	server := Server{
		services: services,
		conf:     c,
		echo:     nil,
	}

	return &server
}

func (api *Server) InitializeHandler() *echo.Echo {

	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(api.loggingMiddleware)
	e.Use(api.headerMiddleware)

	// Initialize all sub handlers
	api.NewFileHandler(e.Group(""))

	return e
}

// logging middleware for echo
func (api *Server) loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.WithFields(log.Fields{
			"request_uri": c.Request().RequestURI,
			"protocol":    c.Request().Proto,
			"method":      c.Request().Method,
		}).Debug("Calling Endpoint")

		return next(c)
	}
}

// modify HTTP Header
func (api *Server) headerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Add("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		return next(c)
	}
}

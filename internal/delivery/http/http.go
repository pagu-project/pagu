package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

type HTTPServer struct {
	handler HTTPHandler
	eServer *echo.Echo
	cfg     *config.HTTP
}

type HTTPHandler struct {
	engine *engine.BotEngine
}

func NewHTTPServer(be *engine.BotEngine, cfg *config.HTTP) HTTPServer {
	return HTTPServer{
		handler: HTTPHandler{
			engine: be,
		},
		eServer: echo.New(),
		cfg:     cfg,
	}
}

func (hs *HTTPServer) Start() error {
	log.Info("Starting HTTP Server", "listen", hs.cfg.Listen)
	hs.eServer.POST("/run", hs.handler.Run)

	return hs.eServer.Start(hs.cfg.Listen)
}

type RunRequest struct {
	Command string `json:"command"`
}

type RunResponse struct {
	Result string `json:"result"`
}

func (hh *HTTPHandler) Run(ctx echo.Context) error {
	r := new(RunRequest)
	if err := ctx.Bind(r); err != nil {
		return err
	}

	cmdResult := hh.engine.ParseAndExecute(entity.PlatformIDReserved, ctx.RealIP(), r.Command)

	return ctx.JSON(http.StatusOK, RunResponse{
		Result: cmdResult.Message,
	})
}

func (hs *HTTPServer) Stop() error {
	log.Info("Stopping HTTP Server")

	return hs.eServer.Close()
}

package grpc

import (
	"context"

	"github.com/pagu-project/pagu/internal/entity"
	pagu "github.com/pagu-project/pagu/pkg/proto/gen/go"
)

type paguServer struct {
	*Server
}

func newPaguServer(server *Server) *paguServer {
	return &paguServer{
		Server: server,
	}
}

func (ps *paguServer) Execute(_ context.Context, req *pagu.ExecuteRequest) (*pagu.ExecuteResponse, error) {
	res := ps.engine.ParseAndExecute(entity.PlatformIDWeb, req.Id, req.Command)

	return &pagu.ExecuteResponse{
		Response: res.Message,
	}, nil
}

package grpc

import (
	"context"
	"strings"

	pagu "github.com/pagu-project/Pagu/internal/delivery/grpc/gen/go"
	"github.com/pagu-project/Pagu/internal/entity"
)

type paguServer struct {
	*Server
}

func newPaguServer(server *Server) *paguServer {
	return &paguServer{
		Server: server,
	}
}

func (ps *paguServer) Run(_ context.Context, req *pagu.RunRequest) (*pagu.RunResponse, error) {
	beInput := make(map[string]string)

	tokens := strings.Split(req.Command, " ")
	for _, t := range tokens {
		beInput[t] = t
	}

	res := ps.engine.Run(entity.AppIDgRPC, req.Id, nil, beInput)

	return &pagu.RunResponse{
		Response: res.Message,
	}, nil
}

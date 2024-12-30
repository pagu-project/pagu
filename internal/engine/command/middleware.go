package command

import (
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type MiddlewareFunc func(caller *entity.User, cmd *Command, args map[string]string) error

type MiddlewareHandler struct {
	db     repository.IDatabase
	wallet wallet.IWallet
}

func NewMiddlewareHandler(d repository.IDatabase, w wallet.IWallet) *MiddlewareHandler {
	return &MiddlewareHandler{
		db:     d,
		wallet: w,
	}
}

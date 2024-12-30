package command

import (
	"errors"

	"github.com/pagu-project/pagu/internal/entity"
)

func (*MiddlewareHandler) OnlyAdmin(caller *entity.User, _ *Command, _ map[string]string) error {
	if caller.Role != entity.Admin {
		return errors.New("this command is Only Admin")
	}

	return nil
}

func (*MiddlewareHandler) OnlyModerator(caller *entity.User, _ *Command, _ map[string]string) error {
	if caller.Role != entity.Moderator {
		return errors.New("this command is Only Moderator")
	}

	return nil
}

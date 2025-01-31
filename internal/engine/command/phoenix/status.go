package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (*PhoenixCmd) statusHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	return cmd.RenderFailedTemplate("TODO!")
}

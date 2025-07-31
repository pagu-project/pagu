package market

import (
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/job"
	"github.com/pagu-project/pagu/pkg/log"
)

func (c *MarketCmd) priceHandler(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	priceData, ok := c.priceCache.Get(job.PriceCacheKey)
	if !ok {
		return cmd.RenderFailedTemplate("failed to get price from markets. please try again later")
	}

	tradeOgre, err := strconv.ParseFloat(priceData.TradeOgrePacToUSDT.Price, 64)
	if err != nil {
		log.Error("unable to parse float", "error", err)
	}

	return cmd.RenderResultTemplate(
		"tradeOgre", tradeOgre,
		"azbitPrice", priceData.AzbitPacToUSDT.Price,
	)
}

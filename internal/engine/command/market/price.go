package market

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (m *MarketCmd) handlerPrice(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	priceData, ok := m.priceCache.Get(config.PriceCacheKey)
	if !ok {
		return cmd.ErrorResult(fmt.Errorf("failed to get price from markets. please try again later"))
	}

	bldr := strings.Builder{}
	xeggexPrice, err := strconv.ParseFloat(priceData.XeggexPacToUSDT.LastPrice, 64)
	if err == nil {
		bldr.WriteString(fmt.Sprintf("Xeggex Price: %f	USDT\n https://xeggex.com/market/PACTUS_USDT \n\n",
			xeggexPrice))
	}

	if priceData.AzbitPacToUSDT.Price > 0 {
		bldr.WriteString(fmt.Sprintf("Azbit Price: %f	USDT\n https://azbit.com/exchange/PAC_USDT \n\n",
			priceData.AzbitPacToUSDT.Price))
	}

	return cmd.SuccessfulResult(bldr.String())
}

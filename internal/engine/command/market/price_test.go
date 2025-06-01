package market

import (
	"context"
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/job"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func setup() *MarketCmd {
	priceCache := cache.NewBasic[string, entity.Price](1 * time.Second)
	priceJob := job.NewPrice(context.Background(), priceCache)
	priceJobScheduler := job.NewScheduler()
	priceJobScheduler.Submit(priceJob)
	go priceJobScheduler.Run()
	m := NewMarketCmd(nil, priceCache)

	return m
}

func TestGetPrice(t *testing.T) {
	market := setup()
	time.Sleep(10 * time.Second)

	cmd := &command.Command{}
	result := market.priceHandler(nil, cmd, nil)
	assert.Equal(t, result.Successful, true)
}

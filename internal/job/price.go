package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/log"
)

const (
	PriceCacheKey = "PriceCacheKey"

	_defaultTradeOgreEndpoint  = "https://tradeogre.com/api/v1/ticker/PAC-USDT"
	_defaultAzbitPriceEndpoint = "https://data.azbit.com/api/tickers?currencyPairCode=PAC_USDT"
)

type PriceChecker struct {
	ctx    context.Context
	cache  cache.Cache[string, entity.Price]
	ticker *time.Ticker
}

func NewPrice(ctx context.Context, cch cache.Cache[string, entity.Price]) *PriceChecker {
	return &PriceChecker{
		cache:  cch,
		ticker: time.NewTicker(128 * time.Second),
		ctx:    ctx,
	}
}

func (p *PriceChecker) Start() {
	p.start()
	go p.runTicker()
}

func (p *PriceChecker) start() {
	var (
		wg        sync.WaitGroup
		price     entity.Price
		tradeOgre entity.TradeOgrePriceResponse
		azbit     []entity.AzbitPriceResponse
	)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultTradeOgreEndpoint, &tradeOgre); err != nil {
			log.Error("unable to get TradeOgre price", "error", err)

			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultAzbitPriceEndpoint, &azbit); err != nil {
			log.Error("unable to get Azbit price", "error", err)

			return
		}
	}()

	wg.Wait()

	price.TradeOgrePacToUSDT = tradeOgre
	if len(azbit) > 0 {
		price.AzbitPacToUSDT = azbit[0]
	}

	ok := p.cache.Exists(PriceCacheKey)
	if ok {
		p.cache.Update(PriceCacheKey, price, 0)
	} else {
		p.cache.Add(PriceCacheKey, price, 0)
	}
}

func (p *PriceChecker) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.start()
		}
	}
}

func (*PriceChecker) getPrice(ctx context.Context, endpoint string, priceResponse any) error {
	cli := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return err
	}

	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is %v", res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)

	return dec.Decode(priceResponse)
}

func (p *PriceChecker) Stop() {
	p.ticker.Stop()
}

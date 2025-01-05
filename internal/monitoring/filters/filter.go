package filters

import (
	"pump_fun/internal/models"
)

type FilterPipeline struct {
	filters []Filter
}

func (p *FilterPipeline) AddFilter(filter Filter) {
	p.filters = append(p.filters, filter)
}

func (p *FilterPipeline) ApplyFilters(coin *models.Coin) *models.Coin {
	var data *models.Coin
	for _, filter := range p.filters {
		data = filter(coin)
		if data == nil {
			return nil
		}
	}
	return data
}

type Filter func(*models.Coin) *models.Coin

func HasWebsite() Filter {
	return func(coin *models.Coin) *models.Coin {
		if coin.IPFSData.WebsiteURL == "" {
			return nil
		}
		return coin
	}
}

func HasTwitter() Filter {
	return func(coin *models.Coin) *models.Coin {
		if coin.IPFSData.TwitterURL == "" {
			return nil
		}
		return coin
	}
}

func HasTelegram() Filter {
	return func(coin *models.Coin) *models.Coin {
		if coin.IPFSData.TelegramURL == "" {
			return nil
		}
		return coin
	}
}

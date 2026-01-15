package filters

import (
	"personal_bot/internal/monitoring/models"
)

type FilterInfo struct {
	Name  string
	Fn    Filter
	Value string
}

type Filter func(*models.Coin) *models.Coin

type FilterPipeline struct {
	filters []FilterInfo
}

func (p *FilterPipeline) AddFilter(filter FilterInfo) {
	p.filters = append(p.filters, filter)
}

func (p *FilterPipeline) ApplyFilters(coin *models.Coin) *models.Coin {
	var data *models.Coin

	if len(p.filters) == 0 {
		return coin
	}

	for _, filter := range p.filters {
		data = filter.Fn(coin)
		if data == nil {
			return nil
		}
	}
	return data
}

func HasWebsite() FilterInfo {
	return FilterInfo{
		Name: HasWebsiteFilter,
		Fn: func(coin *models.Coin) *models.Coin {
			if coin.IPFSData.WebsiteURL == "" {
				return nil
			}
			return coin
		},
	}

}

func HasTwitter() FilterInfo {
	return FilterInfo{
		Name: HasTwitterFilter,
		Fn: func(coin *models.Coin) *models.Coin {
			if coin.IPFSData.TwitterURL == "" {
				return nil
			}
			return coin
		},
	}

}

func HasTelegram() FilterInfo {
	return FilterInfo{
		Name: HasTelegramFilter,
		Fn: func(coin *models.Coin) *models.Coin {
			if coin.IPFSData.TelegramURL == "" {
				return nil
			}
			return coin
		},
	}
}

func DevWallet(devWallet string) FilterInfo {
	return FilterInfo{
		Name: DevWalletFilter,
		Fn: func(coin *models.Coin) *models.Coin {
			if coin.CoinData.CreatorAddr != devWallet {
				return nil
			}
			return coin
		},
		Value: devWallet,
	}
}

package transactiondecoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/monitoring/stream/response"
	eventModels "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/models"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pda"

	jsonparse "personal_bot/pkg/json_parse"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
)

func ParseTransactionNotification(ctx context.Context, data []byte, httpClient *rpc.Client, filters filters.FilterPipeline) (*models.Coin, error) {
	tx, err := jsonparse.Decode[response.TransactionNotification](data)
	if err != nil {
		return nil, err
	}

	txInstructions := tx.Params.Result.Transaction.TransactionDetails.Message.Instructions

	for _, inst := range txInstructions {
		if len(inst.Data) < 8 {
			continue
		}

		instructionData, err := base58.Decode(inst.Data)
		if err != nil {
			continue
		}

		migrationV1 := []byte{155, 234, 231, 146, 236, 158, 162, 30}
		migrationV2 := []byte{187, 203, 18, 31, 206, 237, 254, 41}

		if !bytes.Equal(instructionData[:8], migrationV1[:]) && !bytes.Equal(instructionData[:8], migrationV2[:]) {
			continue
		}

		isV1 := bytes.Equal(instructionData[:8], migrationV1[:])
		var tokenIndex int
		if isV1 {
			if len(inst.Accounts) < 3 {
				continue
			}
			tokenIndex = 2
		} else {
			if len(inst.Accounts) < 4 {
				continue
			}
			if inst.Accounts[2] == constants.WSOLTokenAddress {
				tokenIndex = 3
			} else {
				tokenIndex = 2
			}
		}

		token, err := solana.PublicKeyFromBase58(inst.Accounts[tokenIndex])
		if err != nil {
			logger.Error("failed to parse token address from migration instruction: ", err)
			continue
		}

		coin, err := parseCoin(ctx, token, httpClient, filters)
		if err != nil {
			return nil, err
		}

		if coin != nil {
			coin = filters.ApplyFilters(coin)
			if coin != nil {
				logger.Information("found new coin: ", coin.CoinData.Name, " sig: ", tx.Params.Result.Signature)
				return coin, nil
			}
		}
	}

	return nil, nil
}

func parseCoin(ctx context.Context, mint solana.PK, httpClient *rpc.Client, filters filters.FilterPipeline) (*models.Coin, error) {
	accountInfo, err := httpClient.GetAccountInfoWithOpts(ctx, mint, &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingJSONParsed,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if accountInfo.Value.Owner == solana.Token2022ProgramID {
		return parseV2Migration(accountInfo, mint, filters)
	} else {
		return parseV1Migration(ctx, mint, httpClient, filters)
	}
}

func parseV2Migration(accountInfo *rpc.GetAccountInfoResult, mint solana.PK, filters filters.FilterPipeline) (*models.Coin, error) {
	var mintAccountData eventModels.MintAccountData
	mintAccountData, err := jsonparse.Decode[eventModels.MintAccountData](accountInfo.Value.Data.GetRawJSON())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	for _, ext := range mintAccountData.Parsed.Info.Extensions {
		state := ext.State
		coin := models.Coin{}
		isMetadata := state.Name != nil && state.Symbol != nil && state.Uri != nil
		if isMetadata {
			coin.CoinData.TokenAddr = mint.String()
			coin.CoinData.Name = *ext.State.Name
			coin.CoinData.Symbol = *ext.State.Symbol
			if filters.ShouldAccessIPFS() {
				ipfs, err := getIPFSData(*state.Uri)
				if err != nil {
					return nil, err
				}
				coin.IPFSData = *ipfs
			}
			return &coin, nil
		} else {
			logger.Error("not a valid extension")
		}
	}

	return nil, fmt.Errorf("unable to find metadata within mint account data")

}

func parseV1Migration(ctx context.Context, mint solana.PK, httpClient *rpc.Client, filters filters.FilterPipeline) (*models.Coin, error) {
	output, err := pda.GetMetadataAccount(mint)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	accountInfo, err := httpClient.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(output))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	coin, err := getCoinDataFromMetadata(accountInfo.Value.Data.GetBinary(), filters)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return coin, nil
}

func getIPFSData(ipfsURL string) (*models.IPFS, error) {
	resp, err := http.Get(ipfsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ipfsData models.IPFS
	if err := json.NewDecoder(resp.Body).Decode(&ipfsData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &ipfsData, nil
}

func getCoinDataFromMetadata(data []byte, filters filters.FilterPipeline) (*models.Coin, error) {
	var meta eventModels.MetaplexMetadata
	err := borsh.Deserialize(&meta, data[1:])

	coin := models.Coin{
		CoinData: models.MintData{
			Name:      meta.Name,
			Symbol:    meta.Symbol,
			IpfsUrl:   meta.Uri,
			TokenAddr: solana.PublicKeyFromBytes(meta.Mint[:]).String(),
		},
	}

	if filters.ShouldAccessIPFS() {
		ipfs, err := getIPFSData(meta.Uri)
		if err != nil {
			return nil, err
		}
		coin.IPFSData = *ipfs
	}

	return &coin, err
}

package transaction

import (
	"context"
	"fmt"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	wallets "personal_bot/internal/solana/wallet"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func BuildTx(ctx context.Context, rpcClient *rpc.Client, wallet solana.PrivateKey, instructions *[]solana.Instruction) (*solana.Transaction, error) {
	latestHash, err := client.GetLatestBlockhash(ctx, rpcClient)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return nil, err
	}

	opts := []solana.TransactionOption{
		solana.TransactionPayer(wallet.PublicKey()),
	}

	//TODO - Add constants into our address lookup table
	accountLookupMap, err := lookuptable.GetAddressLookupTable(rpcClient)
	if err != nil {
		logger.Error("Error getting address lookup table, proceeding without it: ", err)
	} else {
		opts = append(opts, solana.TransactionAddressTables(accountLookupMap))
	}

	if instructions == nil {
		return nil, fmt.Errorf("instructions passed into BuildTx were nil ")
	}

	tx, err := solana.NewTransaction(*instructions, latestHash.Value.Blockhash, opts...)
	if err != nil {
		logger.Error("Error creating transaction", err)
		return nil, err
	}

	tx.Message.SetVersion(solana.MessageVersionV0)

	err = wallets.SignTx(tx, wallet)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func SendTx(ctx context.Context, rpcClient *rpc.Client, tx *solana.Transaction) (solana.Signature, error) {
	// SIMULATE TRANSACTION
	// txResp, err := rpcClient.SimulateTransaction(bt.BuyTask.Ctx(), bt.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return nil
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH NO OPTS
	// fmt.Println(bt.transaction.String())
	// txResp, err := rpcClient.SendTransaction(bt.BuyTask.Ctx(), bt.transaction)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	txResp, err := rpcClient.SendTransactionWithOpts(ctx, tx, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true})
	if err != nil {
		logger.Error(err)
		return solana.Signature{}, err
	}

	return txResp, nil

}

func ConfirmTx(ctx context.Context, rpcClient *rpc.Client, signature solana.Signature, task tasks.Task, publisher subscriptionhub.Publisher) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(ctx, rpcClient, signature, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			publisher.PublishMessage(task, msg.Err)
			return fmt.Errorf("%v", msg.Err)
		}
		publisher.PublishMessage(task, msg.Message)
	}
	return nil
}

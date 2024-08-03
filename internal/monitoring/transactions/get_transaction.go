package transactions

import (
	"context"
	"fmt"
	"time"
	
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

func GetTransaction(signature string) {
	transaction, err := getTransaction(signature)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	instruction_data, err := getMintInstructionData(transaction)
	fmt.Println("instruction_data: ", instruction_data)
	// ipfs := getIPFS(coin.Ipfs)
	// combine the coin and ipfs struct into TransactionData struct and return
}

// TODO: Refactor the return type to make it reusable in the future, its currently coupled to the rpc package
func getTransaction(signature string) (*rpc.GetTransactionResult, error) {
	cluster := rpc.MainNetBeta

	rpcClient := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		cluster.RPC,
		rate.Every(time.Second), 
		5,                       
	))

	version := uint64(0)
	transaction, err := rpcClient.GetTransaction(
		context.Background(),
		solana.MustSignatureFromBase58(signature),
		&rpc.GetTransactionOpts{
			MaxSupportedTransactionVersion: &version,
			Encoding:                       solana.EncodingBase64,
		},
	)

	if err != nil {
        return nil, err
    }

	return transaction, nil
}

// TODO: Refactor the return type to make it reusable in the future, its currently coupled to the solana package
func getMintInstructionData(transaction *rpc.GetTransactionResult) ([]solana.CompiledInstruction, error) {
	parsed, err := transaction.Transaction.GetTransaction()
	
	if err != nil {
        return nil, err
    }

    return parsed.Message.Instructions, nil
}

func parseInstructionData(instruction_data string) {
	// logic to parse the instruction data into a coin struct
}

func getIPFS(ipfs_url string) {
	// logic to parse the instruction data into a coin struct
}

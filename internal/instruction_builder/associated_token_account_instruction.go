package instructionbuilder

import (
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
)

func GetIdEmponentInstruction(wallet solana.PublicKey, mintAddress solana.PublicKey) *associatedtokenaccount.Instruction {
	idEmponentInstruction := associatedtokenaccount.NewCreateInstruction(wallet, wallet, mintAddress).Build()
	return idEmponentInstruction

}

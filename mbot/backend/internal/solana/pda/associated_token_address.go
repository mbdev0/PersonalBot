package pda

import "github.com/gagliardetto/solana-go"

func FindToken2022AssociatedTokenAddress(
	walletAddress solana.PublicKey,
	mintAddress solana.PublicKey,
) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			walletAddress[:],
			solana.Token2022ProgramID[:],
			mintAddress[:],
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)
}

func FindTokenAssociatedTokenAddress(
	walletAddress solana.PublicKey,
	mintAddress solana.PublicKey,
) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			walletAddress[:],
			solana.TokenProgramID[:],
			mintAddress[:],
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)
}

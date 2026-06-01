package pda

import (
	"encoding/binary"
	"personal_bot/internal/core/constants"

	"github.com/gagliardetto/solana-go"
)

var (
	pumpfun_amm_program = solana.MustPublicKeyFromBase58(constants.PumpFunAMMProgram)
	pumpfun_program     = solana.MustPublicKeyFromBase58(constants.PumpFunProgram)
)

func getPoolAuthority(baseMint string) (solana.PK, error) {
	poolAuth := []byte("pool-authority")

	baseMintAddress, err := solana.PublicKeyFromBase58(baseMint)
	if err != nil {
		return solana.PK{}, err
	}

	seeds := [][]byte{poolAuth, baseMintAddress.Bytes()}

	address, _, err := solana.FindProgramAddress(seeds, pumpfun_program)
	return address, err
}

func GetPoolPda(baseMint string, quoteMint string) (string, error) {
	poolAuth, err := getPoolAuthority(baseMint)
	if err != nil {
		return "", err
	}

	baseMintPubkey, err := solana.PublicKeyFromBase58(baseMint)
	if err != nil {
		return "", err
	}

	quoteMintPubkey, err := solana.PublicKeyFromBase58(quoteMint)
	if err != nil {
		return "", err
	}

	indexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(indexBytes, 0)

	seeds := [][]byte{
		[]byte("pool"),
		indexBytes,
		poolAuth.Bytes(),
		baseMintPubkey.Bytes(),
		quoteMintPubkey.Bytes(),
	}

	pubkey, _, err := solana.FindProgramAddress(seeds, pumpfun_amm_program)
	if err != nil {
		return "", err
	}

	return pubkey.String(), nil
}

func GetPoolTokenAccount(poolAddress, mintAddress string, isNewTokenAccount bool) (string, error) {
	pool, err := solana.PublicKeyFromBase58(poolAddress)
	if err != nil {
		return "", err
	}

	var tokenAccount solana.PK
	if isNewTokenAccount {
		tokenAccount = solana.MustPublicKeyFromBase58(constants.Token2022Program)
	} else {
		tokenAccount = solana.MustPublicKeyFromBase58(constants.TokenProgram)
	}

	mint, err := solana.PublicKeyFromBase58(mintAddress)
	if err != nil {
		return "", err
	}

	poolTokenAccount, _, err := solana.FindAssociatedTokenAddressWithProgram(pool, mint, tokenAccount)
	if err != nil {
		return "", err
	}

	return poolTokenAccount.String(), nil
}

func GetProtocalFeeRecipientTokenAccount(protocolFeeRecipient, tokenProgram, quoteMint string) (string, error) {
	pfa, err := solana.PublicKeyFromBase58(protocolFeeRecipient)
	if err != nil {
		return "", err
	}

	tokenprogram, err := solana.PublicKeyFromBase58(tokenProgram)
	if err != nil {
		return "", err
	}

	quotemint, err := solana.PublicKeyFromBase58(quoteMint)
	if err != nil {
		return "", err
	}

	address, _, err := solana.FindAssociatedTokenAddressWithProgram(pfa, quotemint, tokenprogram)

	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func GetCoinCreatorVaultAuthority(creator string) (string, error) {
	vault := []byte("creator_vault")
	creatorBytes, err := solana.PublicKeyFromBase58(creator)
	if err != nil {
		return "", err
	}
	seeds := [][]byte{vault, creatorBytes.Bytes()}

	address, _, err := solana.FindProgramAddress(seeds, pumpfun_amm_program)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func GetCoinCreatorVaultAta(coinCreatorVaultAuthority, tokenProgram, wsolMint string) (string, error) {
	coinCreatorVaultAuthorityBytes, err := solana.PublicKeyFromBase58(coinCreatorVaultAuthority)
	if err != nil {
		return "", err
	}

	tokenProgramBytes, err := solana.PublicKeyFromBase58(tokenProgram)
	if err != nil {
		return "", err
	}

	wsolMintBytes, err := solana.PublicKeyFromBase58(wsolMint)
	if err != nil {
		return "", err
	}

	address, _, err := solana.FindAssociatedTokenAddressWithProgram(coinCreatorVaultAuthorityBytes, wsolMintBytes, tokenProgramBytes)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func GetUserVolumeAccumulatorAddress(wallet string) (string, error) {
	walletBytes, err := solana.PublicKeyFromBase58(wallet)
	if err != nil {
		return "", err
	}

	userVolumeAccumulator := []byte("user_volume_accumulator")

	seeds := [][]byte{userVolumeAccumulator, walletBytes.Bytes()}
	address, _, err := solana.FindProgramAddress(seeds, pumpfun_amm_program)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func GetUserVolumeAccumulatorWsolTokenAccount(userVolumeAccumulator, quoteTokenProgram, quoteMint string) (string, error) {
	userVolumeAccumulatorBytes, err := solana.PublicKeyFromBase58(userVolumeAccumulator)
	if err != nil {
		return "", err
	}

	quoteTokenProgramBytes, err := solana.PublicKeyFromBase58(quoteTokenProgram)
	if err != nil {
		return "", err
	}

	quoteMintBytes, err := solana.PublicKeyFromBase58(quoteMint)
	if err != nil {
		return "", err
	}

	address, _, err := solana.FindAssociatedTokenAddressWithProgram(userVolumeAccumulatorBytes, quoteMintBytes, quoteTokenProgramBytes)

	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func GetPoolV2Account(baseMint string) (string, error) {
	poolV2 := []byte("pool-v2")
	baseMintBytes, err := solana.PublicKeyFromBase58(baseMint)
	if err != nil {
		return "", err
	}
	seeds := [][]byte{poolV2, baseMintBytes.Bytes()}

	address, _, err := solana.FindProgramAddress(seeds, pumpfun_amm_program)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func GetBuybackFeeRecipientTokenAccount(quoteMint, buybackFeeRecipient, quoteTokenProgram string) (string, error) {

	quoteMintAddy, err := solana.PublicKeyFromBase58(quoteMint)
	if err != nil {
		return "", err
	}
	buybackFeeRecipientAddy, err := solana.PublicKeyFromBase58(buybackFeeRecipient)
	if err != nil {
		return "", err
	}

	quoteTokenProgramAddy, err := solana.PublicKeyFromBase58(quoteTokenProgram)
	if err != nil {
		return "", err
	}

	address, _, err := solana.FindAssociatedTokenAddressWithProgram(quoteMintAddy, buybackFeeRecipientAddy, quoteTokenProgramAddy)
	return address.String(), err
}

func GetMetadataAccount(mint solana.PK) (string, error) {
	metaplexSeed := solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	metadataSeed := []byte("metadata")

	seeds := [][]byte{metadataSeed, metaplexSeed[:], mint[:]}

	out, _, err := solana.FindProgramAddress(seeds, pumpfun_program)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

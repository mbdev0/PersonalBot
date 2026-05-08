package models

import "github.com/gagliardetto/solana-go"

type BuyEvent struct {
	Timestamp                        int64            `json:"timestamp"`
	BaseAmountOut                    uint64           `json:"base_amount_out"`
	MaxQuoteAmountIn                 uint64           `json:"max_quote_amount_in"`
	UserBaseTokenReserves            uint64           `json:"user_base_token_reserves"`
	UserQuoteTokenReserves           uint64           `json:"user_quote_token_reserves"`
	PoolBaseTokenReserves            uint64           `json:"pool_base_token_reserves"`
	PoolQuoteTokenReserves           uint64           `json:"pool_quote_token_reserves"`
	QuoteAmountIn                    uint64           `json:"quote_amount_in"`
	LpFeeBasisPoints                 uint64           `json:"lp_fee_basis_points"`
	LpFee                            uint64           `json:"lp_fee"`
	ProtocolFeeBasisPoints           uint64           `json:"protocol_fee_basis_points"`
	ProtocolFee                      uint64           `json:"protocol_fee"`
	QuoteAmountInWithLpFee           uint64           `json:"quote_amount_in_with_lp_fee"`
	UserQuoteAmountIn                uint64           `json:"user_quote_amount_in"`
	Pool                             solana.PublicKey `json:"pool"`
	User                             solana.PublicKey `json:"user"`
	UserBaseTokenAccount             solana.PublicKey `json:"user_base_token_account"`
	UserQuoteTokenAccount            solana.PublicKey `json:"user_quote_token_account"`
	ProtocolFeeRecipient             solana.PublicKey `json:"protocol_fee_recipient"`
	ProtocolFeeRecipientTokenAccount solana.PublicKey `json:"protocol_fee_recipient_token_account"`
	CoinCreator                      solana.PublicKey `json:"coin_creator"`
	CoinCreatorFeeBasisPoints        uint64           `json:"coin_creator_fee_basis_points"`
	CoinCreatorFee                   uint64           `json:"coin_creator_fee"`
	TrackVolume                      bool             `json:"track_volume"`
	TotalUnclaimedTokens             uint64           `json:"total_unclaimed_tokens"`
	TotalClaimedTokens               uint64           `json:"total_claimed_tokens"`
	CurrentSolVolume                 uint64           `json:"current_sol_volume"`
	LastUpdateTimestamp              int64            `json:"last_update_timestamp"`
	MinBaseAmountOut                 uint64           `json:"min_base_amount_out"`
	IxName                           string           `json:"ix_name"`
	CashbackFeeBasisPoints           uint64           `json:"cashback_fee_basis_points"`
	Cashback                         uint64           `json:"cashback"`
	BuybackFeeBasisPoints            uint64           `json:"buyback_fee_basis_points"`
	BuybackFee                       uint64           `json:"buyback_fee"`
}

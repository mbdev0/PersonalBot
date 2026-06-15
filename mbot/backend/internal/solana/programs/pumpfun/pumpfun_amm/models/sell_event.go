package models

import "github.com/gagliardetto/solana-go"

type SellEvent struct {
	Timestamp                        int64            `json:"timestamp"`
	BaseAmountIn                     uint64           `json:"base_amount_in"`
	MinQuoteAmountOut                uint64           `json:"min_quote_amount_out"`
	UserBaseTokenReserves            uint64           `json:"user_base_token_reserves"`
	UserQuoteTokenReserves           uint64           `json:"user_quote_token_reserves"`
	PoolBaseTokenReserves            uint64           `json:"pool_base_token_reserves"`
	PoolQuoteTokenReserves           uint64           `json:"pool_quote_token_reserves"`
	QuoteAmountOut                   uint64           `json:"quote_amount_out"`
	LpFeeBasisPoints                 uint64           `json:"lp_fee_basis_points"`
	LpFee                            uint64           `json:"lp_fee"`
	ProtocolFeeBasisPoints           uint64           `json:"protocol_fee_basis_points"`
	ProtocolFee                      uint64           `json:"protocol_fee"`
	QuoteAmountOutWithoutLpFee       uint64           `json:"quote_amount_out_without_lp_fee"`
	UserQuoteAmountOut               uint64           `json:"user_quote_amount_out"`
	Pool                             solana.PublicKey `json:"pool"`
	User                             solana.PublicKey `json:"user"`
	UserBaseTokenAccount             solana.PublicKey `json:"user_base_token_account"`
	UserQuoteTokenAccount            solana.PublicKey `json:"user_quote_token_account"`
	ProtocolFeeRecipient             solana.PublicKey `json:"protocol_fee_recipient"`
	ProtocolFeeRecipientTokenAccount solana.PublicKey `json:"protocol_fee_recipient_token_account"`
	CoinCreator                      solana.PublicKey `json:"coin_creator"`
	CoinCreatorFeeBasisPoints        uint64           `json:"coin_creator_fee_basis_points"`
	CoinCreatorFee                   uint64           `json:"coin_creator_fee"`
	CashbackFeeBasisPoints           uint64           `json:"cashback_fee_basis_points"`
	Cashback                         uint64           `json:"cashback"`
	BuybackFeeBasisPoints            uint64           `json:"buyback_fee_basis_points"`
	BuybackFee                       uint64           `json:"buyback_fee"`
}

package models

type WalletRepository struct {
	Id         string `db:"id"`
	WalletName string `db:"wallet_name"`
	Chain      string `db:"chain"`
	PrivateKey string `db:"private_key"`
}

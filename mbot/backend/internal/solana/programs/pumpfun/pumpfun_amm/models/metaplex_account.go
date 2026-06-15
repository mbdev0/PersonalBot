package models

type MetaplexMetadata struct {
	UpdateAuthority      [32]byte `borsh:"update_authority"`
	Mint                 [32]byte `borsh:"mint"`
	Name                 string   `borsh:"name"`
	Symbol               string   `borsh:"symbol"`
	Uri                  string   `borsh:"uri"`
	SellerFeeBasisPoints uint16   `borsh:"seller_fee_basis_points"`
}

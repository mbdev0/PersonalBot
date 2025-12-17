import { type Wallet, type WalletDto, type WalletPut, type WalletPutDTO } from '../types/wallet';

export function walletToDto(wallet: Wallet): WalletDto {
  return {
    id: wallet.id,
    wallet_name: wallet.wallet_name,
    chain: wallet.chain,
    public_key: wallet.public_key,
  };
}

export function dtoToWallet(dto: WalletDto): Wallet {
  return {
    id: dto.id,
    wallet_name: dto.wallet_name,
    chain: dto.chain,
    public_key: dto.public_key,
  };
}

export function walletPutToDto(wallet: WalletPut): WalletPutDTO {
  return {
    wallet_name: wallet.wallet_name,
    chain: wallet.chain,
    private_key: wallet.private_key,
  };
}

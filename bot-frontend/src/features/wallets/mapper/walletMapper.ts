import {
  type Wallet,
  type WalletDto,
  type WalletPost,
  type WalletPostDTO,
  type WalletPut,
  type WalletPutDTO,
} from '../types/wallet';

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

export function walletPostToDto(wallet: WalletPost): WalletPostDTO {
  return {
    wallet_name: wallet.wallet_name,
    chain: wallet.chain,
    private_key: wallet.private_key,
  };
}

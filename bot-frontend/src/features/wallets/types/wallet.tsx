// What we recieve from the API
export interface WalletDto {
  id: string;
  wallet_name: string;
  chain: Chain;
  public_key: string;
}

// how it will look in the app
export interface Wallet {
  id: string;
  wallet_name: string;
  chain: Chain;
  public_key: string;
}

// what we send for any updates/post requests
export interface WalletPostDTO {
  wallet_name: string;
  chain: string;
  private_key: string;
}

enum Chain {
  Solana = 'Solana',
}

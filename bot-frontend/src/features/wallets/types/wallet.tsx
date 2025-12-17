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

// what we send for any post requests
export interface WalletPostDTO {
  wallet_name: string;
  chain: string;
  private_key: string;
}

// what we have inside the app on the update screen
export interface WalletPut {
  id: string;
  wallet_name: string;
  chain: string;
  public_key: string; // this will be what's displayed when you click on the edit
  private_key: string; // will always be empty when first displayed
}

// what we send to the API
export interface WalletPutDTO {
  wallet_name: string;
  chain: string;
  private_key: string;
}

export enum Chain {
  Solana = 'Solana',
  BSC = 'BSC',
}

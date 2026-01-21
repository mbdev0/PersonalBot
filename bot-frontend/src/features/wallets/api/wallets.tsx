// this will hold the API calls.. or in otherwords the query functions
// we will use fetch

import { API_BASE } from '../../../config/urls';
import type { Wallet, WalletDto, WalletPost, WalletPut } from '../types/wallet';
import { dtoToWallet, walletPostToDto, walletPutToDto } from '../mapper/walletMapper';

//get wallets
//TODO: map each wallet into Wallet instead of Dto
export async function getWallets(): Promise<Wallet[]> {
  let resp = await fetch(API_BASE + '/wallet/wallets');
  if (!resp.ok) {
    throw new Error('Failed to fetch wallets: ' + resp.json());
  }

  const walletDtos: WalletDto[] = await resp.json();
  const wallets = walletDtos.map(dtoToWallet);

  return wallets;
}

//get wallet by id -> maybe have the id as the query key? -> Do we need this though?
export async function getWalletsById(id: string): Promise<WalletDto> {
  let resp = await fetch(API_BASE + `/wallet/wallets/${id}}`);
  if (!resp.ok) {
    throw new Error(`Failed to fetch wallet with ID: ${id} - Error: ${resp.json()}`);
  }

  return resp.json();
}

//put for update
export async function updateWallet(wallet: WalletPut): Promise<WalletDto> {
  const walletDto = walletPutToDto(wallet);
  const url = API_BASE + `/wallet/wallets/${wallet.id}`;

  const response = await fetch(url, {
    method: 'PUT',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(walletDto),
  });

  if (!response.ok) {
    throw new Error(`Failed to update wallet: ${response.json()}`);
  }

  return response.json();
}

//post
export async function postWallet(wallet: WalletPost) {
  const walletDto = walletPostToDto(wallet);
  const url = API_BASE + `/wallet/wallets`;

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(walletDto),
  });

  if (!response.ok) {
    throw new Error(`Failed to add new wallet: ${response.json()}`);
  }

  return response.status;
}

//delete
export async function deleteWallet(id: string) {
  const url = API_BASE + `/wallet/wallets/${id}`;

  const response = await fetch(url, {
    method: 'DELETE',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to add new wallet: ${response.json()}`);
  }

  return response.status;
}

import { API_BASE } from '../../../config/urls';
import type { Wallet, WalletDto, WalletPost, WalletPut } from '../types/wallet';
import { dtoToWallet, walletPostToDto, walletPutToDto } from '../mapper/walletMapper';

//get wallets
export async function getWallets(): Promise<Wallet[]> {
  const resp = await fetch(API_BASE + '/wallet/wallets');
  if (!resp.ok) {
    throw new Error('Failed to fetch wallets: ' + resp.json());
  }

  const walletDtos: WalletDto[] = await resp.json();
  const wallets = walletDtos.map(dtoToWallet);

  return wallets;
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
    const text = await response.text();
    throw new Error(text || `Request failed (${response.status})`);
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
    const text = await response.text();
    throw new Error(text || `Request failed (${response.status})`);
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
    const text = await response.text();
    throw new Error(text || `Request failed (${response.status})`);
  }

  return response.status;
}

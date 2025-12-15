// this will hold the API calls.. or in otherwords the query functions
// we will use fetch

import { API_BASE } from '../../../config/urls';
import type { WalletDto } from '../types/wallet';

//get wallets
export async function getWallets(): Promise<WalletDto[]> {
  let resp = await fetch(API_BASE + '/wallet/wallets');
  if (!resp.ok) {
    throw new Error('Failed to fetch wallets: ' + resp.json());
  }

  return resp.json();
}

//get wallet by id -> maybe have the id as the query key?
export async function getWalletsById(id: string): Promise<WalletDto> {
  let resp = await fetch(API_BASE + `/wallet/wallets/${id}}`);
  if (!resp.ok) {
    throw new Error(`Failed to fetch wallet with ID: ${id} - Error: ${resp.json()}`);
  }

  return resp.json();
}

//put for update

//delete

//post

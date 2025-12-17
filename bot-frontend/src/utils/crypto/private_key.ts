import nacl from "tweetnacl";
import bs58 from 'bs58'

export function getPubkeyFromPrivateKey(priv_key: string): string | null {
  try {
    const key_pair = nacl.sign.keyPair.fromSecretKey(bs58.decode(priv_key));
    return bs58.encode(key_pair.publicKey);
  } catch (e) {
    return null;  
  }
}

export function isValidSolanaPrivateKey(priv_key: string): boolean {
  if (!priv_key || priv_key.trim() === '') {
    return false;
  }

  try {
    const decoded = bs58.decode(priv_key);

    if (decoded.length !== 64) {
      return false;
    }

    nacl.sign.keyPair.fromSecretKey(decoded);

    return true;
  } catch (e) {
    return false;
  }
}
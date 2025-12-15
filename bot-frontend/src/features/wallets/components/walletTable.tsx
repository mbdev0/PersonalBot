import { useState } from 'react';
import { useWallets } from '../hooks/useWallets';
import './walletTable.css';
import Modal from '../../../components/modal';
import WalletEntry from './walletEntry';

function WalletTable() {
  const { isPending, isError, data, error } = useWallets();
  const [isModalShowing, setModalShow] = useState(false);

  if (isPending) {
    return <div className="loading">Loading...</div>;
  }

  if (isError) {
    return <div className="error"> Error: {error.message}</div>;
  }

  return (
    <>
      <Modal isOpen={isModalShowing} onClose={() => setModalShow(false)}>
        <WalletEntry></WalletEntry>
      </Modal>
      <table>
        <thead>
          <tr>
            <th scope="col">Wallet Name</th>
            <th scope="col">Chain</th>
            <th scope="col">Private Key</th>
            <th scope="col">Edit</th>
          </tr>
        </thead>
        <tbody>
          {data?.map((wallet) => (
            <tr key={wallet.wallet_name}>
              <td>{wallet.wallet_name}</td>
              <td>{wallet.chain}</td>
              <td>{`${wallet.public_key.slice(0, 4)}...${wallet.public_key.slice(-4)}`}</td>
              <td>
                <button onClick={() => setModalShow(true)}>✏️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}

export default WalletTable;

package datastream

import (
	"context"
	"personal_bot/internal/solana/monitoring/stream"
	"personal_bot/pkg/logger"
	"sync"
)

type SolanaStream struct {
	activeConnections map[string]*connections
	ctx               context.Context
	mu                sync.Mutex
}

type connections struct {
	addressConnections map[string]*addressConnection
}

type addressConnection struct {
	stream     chan []byte
	references int32
	cancel     context.CancelFunc
}

func NewSolanaStream(ctx context.Context) *SolanaStream {
	return &SolanaStream{
		activeConnections: map[string]*connections{},
		ctx:               ctx,
	}
}

func (ss *SolanaStream) SubscribeToTransaction(wsurl, program string) (out <-chan []byte, err error) {
	return ss.subscribeToStream(wsurl, program, stream.NewStartGeyserTransactionStream)

}

func (ss *SolanaStream) SubscribeToAccountStream(wsurl, account string) (out <-chan []byte, err error) {
	return ss.subscribeToStream(wsurl, account, stream.NewGeyserStreamAccountInfo)
}

func (ss *SolanaStream) subscribeToStream(wsurl, account string, subFunc func(ctx context.Context, address, wsurl string, output chan<- []byte) error) (out <-chan []byte, err error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, ok := ss.activeConnections[wsurl]; !ok {
		ss.activeConnections[wsurl] = &connections{
			addressConnections: map[string]*addressConnection{},
		}
	}

	ws := ss.activeConnections[wsurl]

	if _, ok := ws.addressConnections[account]; !ok {

		streamContext, cancel := context.WithCancel(ss.ctx)
		messageChan := make(chan []byte, 100)
		ws.addressConnections[account] = &addressConnection{
			stream:     messageChan,
			references: 0,
			cancel:     cancel,
		}

		go func() {
			err := subFunc(streamContext, account, wsurl, messageChan)
			if err != nil {
				logger.Error(err)
				close(messageChan)
				return
			}
		}()
	}

	ws.addressConnections[account].references++
	return ws.addressConnections[account].stream, nil
}

func (ss *SolanaStream) Unsubscribe(wsurl, address string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ws, ok := ss.activeConnections[wsurl]
	if !ok {
		return
	}

	connection, ok := ws.addressConnections[address]
	if !ok {
		return
	}

	connection.references--
	if connection.references >= 0 {
		connection.cancel()
		close(connection.stream)
	}
}

package datastream

type DataStream interface {
	SubscribeToTransaction(wsurl, program string) (out <-chan []byte, err error)
	SubscribeToAccountStream(wsurl, account string) (out <-chan []byte, err error)
}

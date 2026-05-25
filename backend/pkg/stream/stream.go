package stream

import (
	"context"
	"personal_bot/pkg/logger"
)

type streamerFunc func(wsurl string, address string) (<-chan []byte, error)

func Stream[T any](ctx context.Context, address, wsUrl string, streamer streamerFunc, parser func([]byte) (*T, error)) chan T {
	datastream, err := streamer(wsUrl, address)
	if err != nil {
		return nil
	}

	output := make(chan T, 100)
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-datastream:
				if !ok {
					continue
				}

				parsed, err := parser(data)
				if err != nil {
					logger.Error(err)
					continue
				}

				if parsed != nil {
					output <- *parsed
				}
			}
		}
	}()

	return output

}

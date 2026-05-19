package stream

import (
	"context"
	"personal_bot/pkg/logger"
)

type streamerFunc func(context.Context, string, string, chan<- []byte) error

func Stream[T any](ctx context.Context, address, wsUrl string, streamer streamerFunc, parser func([]byte) (T, error)) chan T {
	dataStream := make(chan []byte, 100)
	go func(ctx context.Context, s1, s2 string, c chan<- []byte) {
		defer close(dataStream)
		if err := streamer(ctx, s1, s2, c); err != nil {
			logger.Error(err)
			return
		}
	}(ctx, address, wsUrl, dataStream)

	output := make(chan T, 100)
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-dataStream:
				if !ok {
					continue
				}

				parsed, err := parser(data)
				if err != nil {
					logger.Error(err)
					continue
				}

				output <- parsed

			}
		}
	}()

	return output

}

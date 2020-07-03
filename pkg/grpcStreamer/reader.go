package grpcStreamer

import "io"

func NewStreamReader(config StreamReaderConfig) *streamReader {
	streamer := &streamReader{
		recvByte: config.RecvByte,
	}

	return streamer
}

type streamReader struct {
	recvByte RecvByte
}

func (s *streamReader) Read(p []byte) (n int, err error) {
	n = 0
	for {
		if n >= len(p) {
			return n, nil
		}

		data, err := s.recvByte()
		if err != nil {
			if err == io.EOF {
				return n, nil
			} else {
				return 0, err
			}
		} else {
			p[n] = data
			n += 1
		}
	}
}

type RecvByte func() (byte, error)

type StreamReaderConfig struct {
	RecvByte RecvByte
}

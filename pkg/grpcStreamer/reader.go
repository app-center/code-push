package grpcStreamer

import "io"

func NewStreamReader(config StreamReaderConfig) *streamReader {
	streamer := &streamReader{
		recvByte: config.RecvByte,
	}

	return streamer
}

type streamReader struct {
	recvByte  RecvByte
	streamEOF bool
}

func (s *streamReader) Read(p []byte) (n int, err error) {
	if s.streamEOF {
		return 0, io.EOF
	}

	n = 0
	for {
		if n >= len(p) {
			return n, nil
		}

		data, recvErr := s.recvByte()
		if recvErr != nil {
			if recvErr == io.EOF {
				s.streamEOF = true
				return n, nil
			} else {
				return n, recvErr
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

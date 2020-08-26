package grpc_streamer

func NewSender(sendByte SendByte) *streamSender {
	return &streamSender{sendByte: sendByte}
}

type streamSender struct {
	sendByte SendByte
}

func (s *streamSender) Write(p []byte) (n int, err error) {
	n = 0
	for {
		if n >= len(p) {
			return
		}

		err = s.sendByte(p[n])
		if err != nil {
			return
		} else {
			n += 1
		}
	}
}

type SendByte func(p byte) error

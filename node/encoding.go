package node

import (
	"encoding/gob"
	"io"
)

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

type Decoder interface {
	Decode(r io.Reader, msg *MSG) error
}

type GOBDecoder struct{}

func (gb GOBDecoder) Decode(r io.Reader, msg *MSG) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dg DefaultDecoder) Decode(r io.Reader, msg *MSG) error {
	// if we read 1 byte it is just a message, and we don't need to decode it
	optBuff := make([]byte, 1)
	_, err := r.Read(optBuff)
	if err != nil {
		return nil
	}

	//well decode it in next logic, now just set is it a stream
	stream := optBuff[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}

	buff := make([]byte, 1028)
	if _, err = r.Read(buff); err != nil {
		return err
	}

	msg.Payload = buff

	return nil
}

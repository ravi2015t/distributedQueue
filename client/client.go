package client

import (
	"bytes"
	"errors"
)

const defaultScratchBufferSize = 64 * 1024

//Basic represnts an instance of queue
type Basic struct {
	addrs   []string
	buf     bytes.Buffer
	restBuf bytes.Buffer
}

var errBufTooSmall = errors.New("buffer is too small to fit the message")

//NewBasic creates a new instance of client
func NewBasic(addrs []string) *Basic {
	return &Basic{
		addrs: addrs,
	}
}

func (c *Basic) Send(msgs []byte) error {
	_, err := c.buf.Write(msgs)
	return err

}

func (c *Basic) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchBufferSize)
	}

	off := 0

	if c.restBuf.Len() > 0 {
		if c.restBuf.Len() >= len(scratch) {
			return nil, errBufTooSmall
		}
		n, err := c.restBuf.Read(scratch)
		if err != nil {
			return nil, err
		}
		c.restBuf.Reset()
		off += n
	}
	n, err := c.buf.Read(scratch[off:])
	if err != nil {
		return nil, err
	}

	truncated, rest, err := cutToLastMessage(scratch[0 : n+off])
	if err != nil {
		return nil, err
	}

	c.restBuf.Reset()
	c.restBuf.Write(rest)

	return truncated, nil
}

func cutToLastMessage(res []byte) (truncated []byte, rest []byte, err error) {

	n := len(res)
	if n == 0 {
		return res, nil, nil
	}

	if res[n-1] == '\n' {
		return res, nil, nil
	}

	lastPos := bytes.LastIndexByte(res, '\n')

	if lastPos < 0 {
		return nil, nil, errBufTooSmall
	}
	return res[0 : lastPos+1], res[lastPos+1:], nil
}

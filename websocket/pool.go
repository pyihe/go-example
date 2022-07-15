package websocket

import "sync"

var mPool sync.Pool

type message struct {
	mType int
	conn  *wsConn
	data  []byte
}

func (m *message) write(b []byte) {
	m.data = make([]byte, len(b))
	copy(m.data, b)
}

func getMessage() *message {
	data := mPool.Get()
	if data == nil {
		return &message{}
	}
	return data.(*message)
}

func putMessage(m *message) {
	if m == nil {
		return
	}
	m.mType = 0
	m.conn = nil
	m.data = m.data[:0]
	mPool.Put(m)
}

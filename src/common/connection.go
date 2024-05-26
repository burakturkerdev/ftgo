package common

import (
	"encoding/json"
	"net"
)

type Connection struct {
	conn    net.Conn
	content []byte
	readed  uint64
}

func CreateConnection(c net.Conn) *Connection {

	con := &Connection{
		conn:    c,
		content: nil,
		readed:  0,
	}
	con.readed = 0
	return con
}

func (c *Connection) Read() *Connection {
	c.content = make([]byte, 1024)
	c.conn.Read(c.content)
	c.content = trim(c.content)
	c.readed = 0
	// if content is empty, set message to default
	if len(c.content) < 4 {
		c.content = make([]byte, 4)
	}
	return c
}

func (c *Connection) ReadFile(buffer []byte) {

}

func (c *Connection) GetMessage(m *Message) *Connection {
	*m = Message(uint32(c.content[3]))
	c.readed += 4
	return c
}

func (c *Connection) GetString(s *string) {
	*s = string(c.content[c.readed:])
}

func (c *Connection) SendMessage(m Message) {
	c.conn.Write(MessageToBytes(m))
}

func (c *Connection) SendMessageWithData(m Message, s string) {
	c.conn.Write(append(MessageToBytes(m), []byte(s)...))
}

func (c *Connection) GetJson(t any) {
	json.Unmarshal(c.content[c.readed:], t)
}

func (c *Connection) SendJson(t any) {
	j, e := json.Marshal(t)
	if e != nil {
		println("Log => Json marshal error -> " + e.Error())
	}
	bytes := []byte{}

	bytes = append(bytes, MessageToBytes(Success)...)

	bytes = append(bytes, j...)

	c.conn.Write(bytes)
}

func trim(b []byte) []byte {
	for len(b) > 0 && b[len(b)-1] == 0x00 {
		b = b[:len(b)-1]
	}
	return b
}

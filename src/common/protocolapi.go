package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Connection struct {
	conn          net.Conn
	content       []byte
	EOF           bool
	readed        uint64
	messageReaded bool
}

func CreateConnection(c net.Conn) *Connection {

	con := &Connection{
		conn:          c,
		content:       nil,
		readed:        0,
		EOF:           false,
		messageReaded: false,
	}
	con.readed = 0
	return con
}

func (c *Connection) Read() *Connection {
	c.content = make([]byte, ExchangeBufferSize+4)
	n, err := c.conn.Read(c.content)
	if err != nil {
		c.EOF = true
	}
	c.content = c.content[0:n]
	c.readed = 0
	c.messageReaded = false

	if len(c.content) < 4 {
		c.content = make([]byte, 4)
	}
	return c
}

// Sends data with success message.
func (c *Connection) SendData(buffer []byte) {
	bytes := messageToBytes(Success)
	bytes = append(bytes, buffer...)
	c.conn.Write(bytes)
}

func (c *Connection) GetData(buffer *[]byte) {
	if !c.messageReaded {
		log.Fatal("MESSAGE NOT EXTRACTED. FIRSTLY EXTRACT MESSAGE WITH GetMessage() OR IgnoreMessage()")
	}
	*buffer = c.content[c.readed:]

}

func (c *Connection) GetMessage(m *Message) *Connection {
	if c.messageReaded {
		log.Fatal("MESSAGE ALREADY EXTRACTED.")
	}
	*m = Message(uint32(c.content[3]))
	c.readed += 4
	c.messageReaded = true
	return c
}

func (c *Connection) IgnoreMessage() *Connection {
	if c.messageReaded {
		log.Fatal("MESSAGE ALREADY EXTRACTED.")
	}
	c.readed += 4
	c.messageReaded = true
	return c
}

func (c *Connection) GetString(s *string) {
	if !c.messageReaded {
		log.Fatal("MESSAGE NOT EXTRACTED. FIRSTLY EXTRACT MESSAGE WITH GetMessage() OR IgnoreMessage()")
	}
	*s = string(Trim(c.content[c.readed:]))
}

// Sends string with success message.
func (c *Connection) SendString(s string) {
	bytes := messageToBytes(Success)
	bytes = append(bytes, []byte(s)...)
	c.conn.Write(bytes)
}

func (c *Connection) SendMessage(m Message) {
	c.conn.Write(messageToBytes(m))
}

func (c *Connection) SendMessageWithString(m Message, s string) {
	c.conn.Write(append(messageToBytes(m), []byte(s)...))
}

func (c *Connection) SendMessageWithData(m Message, data []byte) {
	c.conn.Write(append(messageToBytes(m), data...))
}

func (c *Connection) GetJson(t any) {
	if !c.messageReaded {
		log.Fatal("MESSAGE NOT EXTRACTED. FIRSTLY EXTRACT MESSAGE WITH GetMessage() OR IgnoreMessage()")
	}
	json.Unmarshal(c.content[c.readed:], t)
}

// Sends json with success message.
func (c *Connection) SendJson(t any) {
	j, e := json.Marshal(t)
	if e != nil {
		fmt.Println("Log => Json marshal error -> " + e.Error())
	}
	bytes := []byte{}

	bytes = append(bytes, messageToBytes(Success)...)

	bytes = append(bytes, j...)

	c.conn.Write(bytes)
}

func Trim(b []byte) []byte {
	for len(b) > 0 && b[len(b)-1] == 0x00 {
		b = b[:len(b)-1]
	}
	return b
}

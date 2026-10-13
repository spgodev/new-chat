package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) WriteJSON(data any) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteJSON(data)
}

func (c *Client) ReadMessage() ([]byte, error) {

	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ConfigureHeartbeat(timeout time.Duration) error {
	err := c.conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(timeout))
	})

	return nil
}

func (c *Client) StartPing(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := c.conn.WriteControl(
			websocket.PingMessage,
			nil,
			time.Now().Add(5*time.Second),
		)
		if err != nil {
			return
		}
	}
}

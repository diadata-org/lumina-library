package scrapers

import ws "github.com/gorilla/websocket"

type wsConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

var _ wsConn = (*ws.Conn)(nil)

package scrapers

type wsConn interface {
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	Close() error
}

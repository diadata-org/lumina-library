package scrapers

import (
	"encoding/json"
	"errors"
)

type mockWsConn struct {
	writeJSONCalls []interface{}
	writeJSONErrs  []error

	readJSONQueue []interface{}
	readJSONErrs  []error

	closeCalled bool
	closeErr    error
}

func (m *mockWsConn) WriteJSON(v interface{}) error {
	m.writeJSONCalls = append(m.writeJSONCalls, v)
	return nil
}
func (m *mockWsConn) ReadJSON(v interface{}) error {
	if len(m.readJSONErrs) > 0 {
		err := m.readJSONErrs[0]
		m.readJSONErrs = m.readJSONErrs[1:]
		return err
	}
	if len(m.readJSONQueue) == 0 {
		return errors.New("nothing to read")
	}
	b, _ := json.Marshal(m.readJSONQueue[0])
	json.Unmarshal(b, v)
	m.readJSONQueue = m.readJSONQueue[1:]
	return nil
}
func (m *mockWsConn) Close() error {
	m.closeCalled = true
	return nil
}

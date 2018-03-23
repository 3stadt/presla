package Handlers

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

// EditorSync provides websocket broadcasting
// When a websocket sends data, the data is broadcasted to all connected websockets, including the sender
func (conf *Conf) EditorSync(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	index := len(conf.SyncedEditorSub) + 1

	go conf.reader(index, ws)

	// Each websocket gets it's own writer, but it is saved into SyncedEditorSub
	// This way it can be called in batch. See main.go
	conf.SyncedEditorSub[index] = &SyncedEditorWriter{
		Ws:     ws,
		Writer: conf.writer,
	}
	return nil
}

/**
Writes a message to the specified websocket
*/
func (conf *Conf) writer(e SyncedEditor, ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(e.Update))
}

/**
Used to read messages from websockets.
When a message is received from a browser, it is pushed to SyncedEditorPub channel.
*/
func (conf *Conf) reader(index int, ws *websocket.Conn) error {
	defer delete(conf.SyncedEditorSub, index)
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ws.Close()
			return err
		}
		conf.SyncedEditorPub <- SyncedEditor{
			Update: string(msg),
		}
	}
}

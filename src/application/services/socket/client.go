package socket

import (
	"bytes"
	"log"
	"net/http"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	userId uuid.UUID
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
}

func NewClient(hub *Hub) *Client {
	return &Client{
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (this *Client) Send() chan []byte {
	return this.send
}

func (this *Client) read() {
	defer func() {
		this.hub.Unsubscribe(this.userId, this)
	}()

	this.conn.SetReadLimit(maxMessageSize)
	this.conn.SetReadDeadline(time.Now().Add(pongWait))
	this.conn.SetPongHandler(func(string) error {
		this.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := this.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Socket closed: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		this.hub.Broadcast(message, nil)
	}
}

func (this *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		this.conn.Close()
	}()

	for {
		select {
		case message, ok := <-this.send:
			this.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				this.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := this.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(this.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-this.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			this.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := this.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (this *Client) Run(ctx *gin.Context) {
	reqctx := reqcontext.From(ctx)
	auth := reqctx.Auth()

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	this.conn = conn
	this.userId = auth.UserId
	this.hub.Subscribe(this)

	go this.read()
	go this.write()
}

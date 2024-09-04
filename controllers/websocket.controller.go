package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type WebsocketController struct {
	DB *gorm.DB
}

func NewWebsocketController(DB *gorm.DB) WebsocketController {
	return WebsocketController{DB}
}

var (
	MaxMessageSize int64 = 512
	writeWait            = 10 * time.Second
	pongWait             = 60 * time.Second
	pingPeriod           = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewAuctionHub() *models.AuctionHub {
	return &models.AuctionHub{
		Clients:    map[*models.UserClient]bool{},
		Broadcast:  make(chan *models.Bid),
		Register:   make(chan *models.UserClient),
		Unregister: make(chan *models.UserClient),
	}
}

func GetMessageTemplate(message *models.Bid) []byte {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Println("template parsing:", err)
	} else {
		log.Println("template parsing: success")
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, message)
	if err != nil {
		log.Println("template execution: ", err)
	}
	return renderedMessage.Bytes()
}

func Run(h *models.AuctionHub) {
	basiController := NewBasicController(initializers.DB)
	messages := basiController.GetBidsForAuction(h.Auction.ID.String())
	for i := len(messages) - 1; i >= 0; i-- {
		h.Messages = append(h.Messages, &messages[i])
	}

	for {
		select {
		case client := <-h.Register:
			h.Lock()
			h.Clients[client] = true
			h.Unlock()
			log.Printf("client %v connected", client.User.ID)
			for _, message := range h.Messages {
				client.Send <- GetMessageTemplate(message)
			}
		case client := <-h.Unregister:
			h.Lock()
			if _, ok := h.Clients[client]; ok {
				close(client.Send)
				log.Printf("client %v disconnected", client.User.ID)
				delete(h.Clients, client)
			}
			h.Unlock()
		case message := <-h.Broadcast:
			h.RLock()
			h.Messages = append(h.Messages, message)
			for client := range h.Clients {
				select {
				case client.Send <- GetMessageTemplate(message):
					log.Printf("message %v", message)
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.RUnlock()
		}
	}
}

func (wc WebsocketController) GetUserBos(userId string) models.User {
	var user models.User
	result := wc.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return models.User{}
	}
	return user
}

func (wc WebsocketController) ServeWS(hub *models.AuctionHub, ctx *gin.Context) {
	fmt.Println("SATU")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	fmt.Println("DUA")
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}
	fmt.Println("TIGA")
	tempUser := wc.GetUserBos(user_id)
	fmt.Println("EMPAT")

	client := &models.UserClient{
		Hub:  hub,
		Conn: conn,
		User: &tempUser,
		Send: make(chan []byte, 256),
	}
	fmt.Println("LIMA")

	client.Hub.Register <- client
	fmt.Println("ENAM")

	go wc.WritePump(client)
	fmt.Println("TUJUH")

	go wc.ReadPump(client)
	fmt.Println("DELAPAN")
}

func (wc WebsocketController) WritePump(c *models.UserClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(msg)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (wc WebsocketController) ReadPump(c *models.UserClient) {
	defer func() {
		fmt.Println("READPUMP DISCONNECT")
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	fmt.Println("READPUMP CONNECT")

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	fmt.Println("READPUMP 2")
	for {
		_, price, err := c.Conn.ReadMessage()
		fmt.Println("READPUMP 3")
		if err != nil {
			fmt.Println("READPUMP 4")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				fmt.Println("READPUMP 5")
			}
			log.Printf("?? error: %v", err)
			fmt.Println("READPUMP 6")
			break
		}
		fmt.Println("MASUK USERNAME READPUMP")

		user := wc.GetUserBos(c.User.ID.String())
		fmt.Println("READPUMP 7")
		now := time.Now()

		bid := &models.WSBid{}
		fmt.Println("READPUMP 8")
		reader := bytes.NewReader(price)
		fmt.Println("READPUMP 9")
		decoder := json.NewDecoder(reader)
		fmt.Println("READPUMP 10")
		err = decoder.Decode(bid)
		float_price, _ := strconv.ParseFloat(bid.Price, 64)
		fmt.Println("READPUMP 11")
		if err != nil {
			fmt.Println("READPUMP 12")
			log.Printf("error: %v", err)
		}
		dummy_bid := &models.Bid{AuctionID: c.Hub.Auction.ID, Auction: *c.Hub.Auction, UserID: user.ID, User: user, BidAmount: float_price, BidTime: now}
		result := wc.DB.Create(dummy_bid)
		fmt.Println("READPUMP 13")
		c.Hub.Broadcast <- dummy_bid
		wc.DB.Model(&c.Hub.Auction).Update("CurrentPrice", float_price)
		fmt.Println("READPUMP 14")
		if result.Error != nil {
			fmt.Println("READPUMP 15")
			log.Println(result.Error)
		}
	}
}

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
	wc 				 	 = NewWebsocketController(initializers.DB)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewAuctionHub() *models.AuctionHub {
	return &models.AuctionHub{
		Clients:    make(map[*models.UserClient]bool),
		Broadcast:  make(chan *models.Bid),
		Register:   make(chan *models.UserClient),
		Unregister: make(chan *models.UserClient),
	}
}

func GetMessageTemplate(message *models.Bid) []byte {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Println("template parsing:", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, message)
	if err != nil {
		log.Println("template execution: ", err)
	}
	return renderedMessage.Bytes()
}

func Run(h *models.AuctionHub) {
	for {
		select {
		case client := <-h.Register:
			h.Lock()
			h.Clients[client] = true
			log.Printf("client %v connected", client.User.ID)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("client %v disconnected", client.User.ID)
			}
		case message := <-h.Broadcast:
			h.Messages = append(h.Messages, message)
			for client := range h.Clients {
				select {
				case client.Send <- GetMessageTemplate(message):
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
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
	if wc.DB == nil {
		fmt.Println("DB IS NIL")
	} else {
		fmt.Println("DB IS NOT NIL")
	}
	tempUser := wc.GetUserBos(user_id)
	fmt.Println("SERVE USERNAME:", user_id)
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

	go WritePump(client)
	fmt.Println("TUJUH")

	go ReadPump(client)
	fmt.Println("DELAPAN")
}

func WritePump(c *models.UserClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("WRITEPUMP 1")
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("WRITEPUMP 2")
				return
			}

			w.Write(msg)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				fmt.Println("WRITEPUMP 3")
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			fmt.Println("WRITEPUMP 4")
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("WRITEPUMP 5")
				return
			}
		}
	}
}

func ReadPump(c *models.UserClient) {
	defer func() {
		c.Conn.Close()
		c.Hub.Unregister <- c
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		fmt.Println("READPUMP 1")
		return nil
	})

	for {
		fmt.Println("READPUMP 2")
		_, price, err := c.Conn.ReadMessage()
		fmt.Println("READPUMP 3")
		if err != nil {
			fmt.Println("READPUMP 4")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("READPUMP 5")
				log.Printf("error: %v", err)
			}
			fmt.Println("READPUMP 6")
			break
		}
		log.Println("Value: %v", string(price))
		user := wc.GetUserBos(c.User.ID.String())
		fmt.Println("READPUMP USERNAME", c.User.ID.String())
		fmt.Println("READPUMP 7")
		now := time.Now()

		bid := &models.WSBid{}
		fmt.Println("READPUMP 8")
		reader := bytes.NewReader(price)
		fmt.Println("READPUMP 9")
		decoder := json.NewDecoder(reader)
		fmt.Println("READPUMP 10")
		err = decoder.Decode(bid)
		fmt.Println("READPUMP 11")
		if err != nil {
			fmt.Println("READPUMP 12")
			log.Printf("error: %v", err)
		}
		c.Hub.Broadcast <- &models.Bid{AuctionID: c.Hub.Auction.ID, Auction: *c.Hub.Auction, UserID: user.ID, User: user, BidAmount: bid.Price, BidTime: now}
		fmt.Println("READPUMP 13")
		result := wc.DB.Create(bid)
		fmt.Println("READPUMP 14")
		if result.Error != nil {
			fmt.Println("READPUMP 15")
			log.Println(result.Error)
		}
	}
}

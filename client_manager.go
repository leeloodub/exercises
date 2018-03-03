// This file contains structs and functions related to managing client connections.
package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

// Client is a struct to hold connection to a client, the associated writer, and message channel.
type Client struct {
	notifications chan *Message
	writer        *bufio.Writer
	conn          net.Conn
}

// NewClient creates a new client.
func NewClient(conn net.Conn) *Client {
	return &Client{
		notifications: make(chan *Message),
		conn:          conn,
		writer:        bufio.NewWriter(conn),
	}
}

// Forward calls Client.write in a new goroutine.
func (c *Client) Forward() {
	go c.write()
}

func (c *Client) write() {
	for msg := range c.notifications {
		var err error
		_, err = c.writer.WriteString(msg.payload)
		if err != nil {
			log.Printf("failed to buffer message %v, error: %v", msg, err)
		}
		err = c.writer.Flush()
		if err != nil {
			log.Printf("failed to write buffered message %v into io.Writer, error: %v", msg, err)
		}
	}
}

// ClientManager is a struct to hold a map of client connecitons and followers, as
// well as logics on sending messages to the correct client channels.
type ClientManager struct {
	clients   map[int]*Client
	followers map[int][]int
}

// NewClientManager creates a new ClientManager.
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:   make(map[int]*Client),
		followers: make(map[int][]int),
	}
}

func (p *ClientManager) processMessage(m *Message) {
	switch m.mtype {
	case "F":
		p.addFollower(m)
		p.notifyUser(m)
	case "P":
		p.notifyUser(m)
	case "B":
		p.notifyAllUsers(m)
	case "U":
		p.removeFollower(m)
	case "S":
		p.notifyFollowers(m)
	}
}

func (p *ClientManager) addFollower(m *Message) {
	if cl, ok := p.followers[m.toUser]; ok {
		p.followers[m.toUser] = append(cl, m.fromUser)
	} else {
		p.followers[m.toUser] = append([]int{}, m.fromUser)
	}
}

func (p *ClientManager) removeFollower(m *Message) {
	if _, ok := p.followers[m.toUser]; ok {
		for i, v := range p.followers[m.toUser] {
			if v == m.fromUser {
				p.followers[m.toUser] = append(p.followers[m.toUser][:i], p.followers[m.toUser][i+1:]...)
			}
		}
	}
}

func (p *ClientManager) notifyFollowers(m *Message) {
	if fls, ok := p.followers[m.fromUser]; ok {
		for _, f := range fls {
			if recipient, ok := p.clients[f]; ok {
				recipient.notifications <- m
			}
		}
	}
}

func (p *ClientManager) notifyUser(m *Message) {
	if recipient, ok := p.clients[m.toUser]; ok {
		recipient.notifications <- m
	}
}

func (p *ClientManager) notifyAllUsers(m *Message) {
	for _, cl := range p.clients {
		cl.notifications <- m
	}
}

// ManageClients accepts connections form clients, creates a Client, and
// starts a writer goroutine for each client.
func (p *ClientManager) ManageClients(l net.Listener) {
	for {
		// Wait for a connections from all clients.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		rw := bufio.NewReader(bufio.NewReader(conn))
		msg, _ := rw.ReadString('\n')
		id, err := strconv.Atoi(strings.Trim(msg, terminationChar))
		if err != nil {
			log.Printf("client id should be an integer, got: %s, err: %v", msg, err)
		}

		cl := NewClient(conn)
		cl.Forward()
		p.clients[id] = cl
		p.followers[id] = []int{}
		defer conn.Close()
	}
}

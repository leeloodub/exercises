// This file contains structs and functions for managing the Event Target endpoint.
package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

const (
	startSeq        = 1
	terminationChar = "\n"
)

// queue is a struct ho hold messages that have not been sent yet.
// Used within the Endpoint struct.
type queue struct {
	next  int
	elems map[int]*Message
}

// newQueue creates a new queue.
func newQueue(startSeq int) *queue {
	return &queue{
		next:  startSeq,
		elems: make(map[int]*Message),
	}
}

// addToQueue adds an element into a dictionary of elements in the queue.
func (q *queue) addToQueue(m *Message) {
	q.elems[m.seq] = m
}

// sendCurrentAndSubsequent checks if current message has the sequence number
// which is the next number that can be forwarded to clients. Checks subsequent
// message numbers to see if they are in the queue and sends them out as well.
func (q *queue) sendCurrentAndSubsequent(processor func(m *Message)) {
	if _, ok := q.elems[q.next]; !ok {
		return
	}
	processor(q.elems[q.next])
	delete(q.elems, q.next)
	q.next++
	q.sendCurrentAndSubsequent(processor)
}

// Endpoint is an endpoint to accept messages from the event source.
type Endpoint struct {
	listener  net.Listener
	conn      net.Conn
	reader    *bufio.Reader
	queue     *queue
	processor func(m *Message)
}

// NewEndpoint creates a new Endpoint.
func NewEndpoint(l net.Listener, queueStart int) *Endpoint {
	return &Endpoint{listener: l, queue: newQueue(queueStart)}
}

// ListenEvents listens to connection from the event source and reads its messages.
func (e *Endpoint) ListenEvents() {
	// Ready to only accept one connection, as we expect exactly one event source.
	var err error
	e.conn, err = e.listener.Accept()
	log.Println("accepted event target connection")
	if err != nil {
		log.Fatal(err)
	}
	defer e.conn.Close()
	e.reader = bufio.NewReader(bufio.NewReader(e.conn))

	for {
		msg, err := e.reader.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			go e.ListenEvents()
			return
		case err != nil:
			log.Println(err)
		}
		m, err := CreateMessageFromStr(msg)
		if err != nil {
			log.Print(err)
		}
		if m == nil {
			continue
		}
		e.queue.addToQueue(m)
		e.queue.sendCurrentAndSubsequent(e.processor)
	}
}

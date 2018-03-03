package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	terminator = "CRLF"
	separator  = "|"
)

// Message is a struct to hold data associated with the parsed event message.
type Message struct {
	seq      int
	mtype    string
	fromUser int
	toUser   int
	payload  string
}

// CreateMessageFromStr parses the input string into a Message, extracting
// message sequence number, type, from and to user id, and sets the raw string
// content as messgae payload. Message parsing simplifies the message forwarding logics,
// whereis the payload is used to transfer the message to the client, unmodified.
func CreateMessageFromStr(in string) (*Message, error) {
	// if in == "" {
	// 	return nil, nil
	// }

	parts := strings.Split(in, separator)
	if len(parts) < 2 {
		return nil, fmt.Errorf("string '%s' format violates the protocol", in)
	}
	msg := &Message{}

	var err error
	msg.seq, err = strconv.Atoi(strings.Trim(parts[0], " "))
	if err != nil {
		return nil, fmt.Errorf("message sequence should be an integer, got: %s, error: %v", parts[0], err)
	}

	msg.mtype = strings.Trim(parts[1], "\n")

	if len(parts) >= 3 {
		msg.fromUser, err = strconv.Atoi(strings.Trim(parts[2], "\n"))
		if err != nil {
			return nil, fmt.Errorf("user id should be an integer, got: %s, error: %v", parts[2], err)
		}
	}

	if len(parts) == 4 {
		msg.toUser, err = strconv.Atoi(strings.Trim(parts[3], "\n"))
		if err != nil {
			return nil, fmt.Errorf("user id should be an integer, got: %s, error: %v", parts[3], err)
		}
	}

	msg.payload = in
	return msg, nil
}

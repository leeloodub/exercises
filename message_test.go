package main

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestCreateMessageFromString(t *testing.T) {
	for _, test := range []struct {
		in   string
		want *Message
		err  error
	}{
		{
			in: "666|F|60|50",
			want: &Message{
				seq:      666,
				mtype:    "F",
				fromUser: 60,
				toUser:   50,
				payload:  "666|F|60|50",
			},
			err: nil,
		},
		{
			in: "542532|B",
			want: &Message{
				seq:      542532,
				mtype:    "B",
				fromUser: 0,
				toUser:   0,
				payload:  "542532|B",
			},
			err: nil,
		},
		{
			in: "1|U|12|9",
			want: &Message{
				seq:      1,
				mtype:    "U",
				fromUser: 12,
				toUser:   9,
				payload:  "1|U|12|9",
			},
			err: nil,
		},
	} {
		if got, err := CreateMessageFromStr(test.in); !reflect.DeepEqual(got, test.want) {
			log.Println(err)
			t.Errorf("CreateMessageFromStr(%v) = %v, want %v", test.in, got, test.want)
		}
	}
}

func TestCreateMessageFromString_Fails(t *testing.T) {
	for _, test := range []struct {
		in   string
		want *Message
		err  error
	}{
		{
			in:   "foo test| 123",
			want: nil,
			err:  errors.New("should be an integer"),
		},
	} {
		if _, err := CreateMessageFromStr(test.in); !strings.Contains(err.Error(), test.err.Error()) {
			t.Errorf("CreateMessageFromStr(%v) returned error %v want error %v", test.in, err, test.err)
		}
	}
}

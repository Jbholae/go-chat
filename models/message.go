package models

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID `json:"id"`
	UserId  uuid.UUID `json:"user_id"`
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  Room      `json:"target"`
	Sender  Client    `json:"sender"`
}

func (message *Message) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	return nil
}

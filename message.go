// Copyright 2016 Raul Santos Lebrato
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gowebsockets

import "encoding/json"

// Message is the webrtc-message that will be send by and to the server.
// Name is the name of the message, it will be used for the listeners.
// Content is the content that will be sent, use encoding/json to send and receive.
//
// Example:
// 	gowebsockets.NewMessage("message", struct{ Text string; Number int }{"Welcome from the Server!", 42})
type Message struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
	from    string
	to      string
	rooms   []string
}

// NewMessage is a helper function to create a Message.
// content can be a struct of any type.
// content may also be a []byte variable if you haven't unmarshaled it yet (it will do it for you).
func NewMessage(name string, content interface{}) *Message {
	if bytes, ok := content.([]byte); ok {
		var data *map[string]interface{}
		err := json.Unmarshal(bytes, &data)
		if err == nil {
			content = data
		}
	}
	msg := &Message{
		Name:    name,
		Content: content,
	}
	return msg
}

func (m *Message) toBytes() ([]byte, error) {
	bytes, err := json.Marshal(m.Content)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

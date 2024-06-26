/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package web

import (
	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/server/web"
)

type Input = web.Input

type Output = web.Output

type WSStream struct {
	conn *websocket.Conn
}

func NewWSStream(url string, config *tls.Config) (*WSStream, error) {
	logger.Debugf("Connecting to %s", url)
	dialer := &websocket.Dialer{TLSClientConfig: config}
	ws, _, err := dialer.Dial(url, nil)
	logger.Infof("Successfully connected to websocket")
	if err != nil {
		logger.Errorf("Dial failed: %s\n", err.Error())
		return nil, err
	}
	return &WSStream{conn: ws}, nil
}

func (c *WSStream) Send(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *WSStream) Recv(v interface{}) error {
	return c.conn.ReadJSON(v)
}

func (c *WSStream) Close() error {
	return c.conn.Close()
}

func (c *WSStream) Result() ([]byte, error) {
	output := &Output{}
	if err := c.Recv(output); err != nil {
		return nil, err
	}
	return output.Raw, nil
}

func (c *WSStream) SendInput(in []byte) error {
	return c.Send(&Input{Raw: in})
}

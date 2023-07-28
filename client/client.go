package client

import (
	"github.com/coalalib/coalago"
	coalaMsg "github.com/coalalib/coalago/message"
	log "github.com/ndmsystems/golog"
	"net"
)

type Client struct {
	coala     *coalago.Client
	addr      string
	publicKey string
}

func New(addr string, publicKey []byte) *Client {
	c := new(Client)
	c.coala = coalago.NewClient()
	c.addr = addr
	c.publicKey = string(publicKey)

	return c
}

func (c *Client) SendAlive() {
	requestMessage := coalaMsg.NewCoAPMessage(coalaMsg.CON, coalaMsg.GET)
	requestMessage.SetURIPath("/live?key=" + c.publicKey)
	requestMessage.SetURIQuery("key", c.publicKey)

	address, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = c.coala.Send(requestMessage, address.String())
	if err != nil {
		log.Error(err)
		return
	}
}

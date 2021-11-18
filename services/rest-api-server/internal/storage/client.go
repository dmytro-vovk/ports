package storage

import (
	"context"

	"github.com/dmytro-vovk/ports/services/protocol"
	"google.golang.org/grpc"
)

type Client struct {
	address string
	conn    *grpc.ClientConn
	client  protocol.StorageClient
}

func New(address string) *Client {
	c := Client{address: address}

	return &c
}

func (c *Client) connect() error {
	if c.conn != nil {
		return nil
	}

	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = protocol.NewStorageClient(c.conn)

	return nil
}

func (c *Client) Store(ctx context.Context, request *protocol.StorePortRequest) error {
	if err := c.connect(); err != nil {
		return err
	}

	sc, err := c.client.Store(ctx, grpc.EmptyCallOption{})
	if err != nil {
		return err
	}

	return sc.Send(request)
}

func (c *Client) Get(ctx context.Context, request *protocol.GetPortRequest) (*protocol.Data, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	return c.client.Get(ctx, request)
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

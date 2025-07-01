package yellowstone

import (
	"context"
	"golang/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type YellowstoneClient struct {
	conn   *grpc.ClientConn
	client proto.GeyserClient
}

func NewYellowstoneClient(addr string) (*YellowstoneClient, error) {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	client := proto.NewGeyserClient(conn)
	return &YellowstoneClient{conn: conn, client: client}, nil
}

func (c *YellowstoneClient) SubscribeAndListen(ctx context.Context, req *proto.SubscribeRequest, onUpdate func(*proto.SubscribeUpdate)) error {
	stream, err := c.client.Subscribe(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(req); err != nil {
		return err
	}
	for {
		update, err := stream.Recv()
		if err != nil {
			return err
		}
		onUpdate(update)
	}
}

func (c *YellowstoneClient) Close() error {
	return c.conn.Close()
}

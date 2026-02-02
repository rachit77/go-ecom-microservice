package account

import (
	"context"

	"github.com/rachit77/go-ecom-microservice/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	// conn, err := grpc.Dial(url, grpc.WithInsecure())
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	c := pb.NewAccountServiceClient(conn)
	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {

}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {

}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {

}

func (c *Client) GetAccounts(ctx context.Context, skip uint, take uint) ([]Account, error) {

}

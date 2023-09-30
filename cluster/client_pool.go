package cluster

import (
	"context"
	"errors"

	"github.com/stream1080/godis/resp/client"

	pool "github.com/jolestar/go-commons-pool/v2"
)

type ConnFactory struct {
	Peer string
}

func (f *ConnFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {

	c, err := client.MakeClient(f.Peer)
	if err != nil {
		return nil, err
	}

	c.Start()

	return pool.NewPooledObject(c), nil
}

func (f *ConnFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type not match")
	}

	c.Close()
	return nil

}

func (f *ConnFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (f *ConnFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (f *ConnFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

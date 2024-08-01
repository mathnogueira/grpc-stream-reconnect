package grpcstreamreconnect

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type conn struct {
	mutex       sync.Mutex
	grpcConn    *grpc.ClientConn
	openStreams []*stream
}

// Invoke implements grpc.ClientConnInterface.
func (c *conn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return c.grpcConn.Invoke(ctx, method, args, reply, opts...)
}

// NewStream implements grpc.ClientConnInterface.
func (c *conn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	newStream, err := newStream(ctx, c.grpcConn, desc, method, opts...)
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	c.openStreams = append(c.openStreams, newStream)
	c.mutex.Unlock()
	return newStream, nil
}

func NewClient(target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	c, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	connection := &conn{
		grpcConn: c,
	}

	connection.watchReconnections()

	return connection, nil
}

func (c *conn) watchReconnections() {
	ctx := context.Background()
	go func() {
		shouldReconnect := false
		for {
			c.grpcConn.WaitForStateChange(ctx, connectivity.TransientFailure)
			state := c.grpcConn.GetState()

			if state == connectivity.TransientFailure || state == connectivity.Connecting {
				shouldReconnect = true
			}

			if state == connectivity.Ready && shouldReconnect {
				shouldReconnect = false
				c.reconnectStreams()
			}
		}
	}()
}

func (c *conn) reconnectStreams() {
	for _, stream := range c.openStreams {
		stream.Reconnect()
	}
}

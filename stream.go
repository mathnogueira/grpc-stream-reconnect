package grpcstreamreconnect

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type stream struct {
	mutex sync.Mutex
	conn  grpc.ClientConnInterface
	s     grpc.ClientStream

	// Constructor args
	desc   *grpc.StreamDesc
	method string
	opts   []grpc.CallOption
}

func newStream(ctx context.Context, conn grpc.ClientConnInterface, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (*stream, error) {
	stream := &stream{
		mutex:  sync.Mutex{},
		conn:   conn,
		desc:   desc,
		method: method,
		opts:   opts,
	}

	err := stream.connect(ctx)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (s *stream) connect(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stream, err := s.conn.NewStream(ctx, s.desc, s.method, s.opts...)
	if err != nil {
		return err
	}

	s.s = stream
	return nil
}

// CloseSend implements grpc.ClientStream.
func (s *stream) CloseSend() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.CloseSend()
}

// Context implements grpc.ClientStream.
func (s *stream) Context() context.Context {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.Context()
}

// Header implements grpc.ClientStream.
func (s *stream) Header() (metadata.MD, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.Header()
}

// RecvMsg implements grpc.ClientStream.
func (s *stream) RecvMsg(m any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.RecvMsg(m)
}

// SendMsg implements grpc.ClientStream.
func (s *stream) SendMsg(m any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.SendMsg(m)
}

// Trailer implements grpc.ClientStream.
func (s *stream) Trailer() metadata.MD {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.s.Trailer()
}

func (s *stream) Reconnect() error {
	err := s.CloseSend()
	if err != nil {
		return err
	}

	return s.connect(context.Background())
}

var _ grpc.ClientStream = &stream{}

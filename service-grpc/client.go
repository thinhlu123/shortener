package service_grpc

import (
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"shortener/model/message"
	"strconv"
	"sync"
)

var ClientgRPC *APIServiceClient
var grpcServerAddr string

type APIServiceClient struct {
	Client message.APIServiceClient
	Cons   map[string]*ClientConn
	lock   *sync.Mutex
}

type ClientConn struct {
	id     string
	isBusy bool
	Conn   *grpc.ClientConn
	lock   *sync.Mutex
}

func newConnClient(serverAddr string, isSecure bool) *ClientConn{
	var opts []grpc.DialOption
	if isSecure {

	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	id := rand.Intn(999999999) + 1000000000

	clientConn := &ClientConn{
		isBusy: false,
		Conn:   conn,
		id:     strconv.Itoa(id),
	}
	ClientgRPC.lock.Lock()
	ClientgRPC.Cons[clientConn.id] = clientConn
	ClientgRPC.lock.Unlock()

	return clientConn
}

func SetupgRPCClient(serverAddr string, isSecure bool) *APIServiceClient {
	grpcServerAddr = serverAddr
	clientConn := newConnClient(serverAddr, isSecure)

	conns := make(map[string]*ClientConn)
	conns[clientConn.id] = clientConn

	ClientgRPC = &APIServiceClient{
		Client: message.NewAPIServiceClient(clientConn.Conn),
		Cons:   conns,
	}

	return ClientgRPC
}

func PickConn(createNew bool) *ClientConn {
	if !createNew {
		ClientgRPC.lock.Lock()
		for _, conn := range ClientgRPC.Cons {
			conn.lock.Lock()
			if !conn.isBusy {
				conn.isBusy = true
				conn.lock.Unlock()
				ClientgRPC.lock.Unlock()
				return conn
			}
			conn.lock.Unlock()
		}
		ClientgRPC.lock.Unlock()
	}

	//Create new conn
	con := newConnClient(grpcServerAddr, false)
	return con
}

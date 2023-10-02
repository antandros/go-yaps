package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/antandros/go-yaps/helper"
	"github.com/antandros/go-yaps/yaperror"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	client       PluginProtocolClient
	ctx          context.Context
	logger       *zap.Logger
	conn         *grpc.ClientConn
	responseTime time.Duration
	isRemote     bool
	addr         string
}

func (c *Client) WaitConnect() {
	state := c.conn.GetState()
	for {
		if state == connectivity.Ready {
			return
		}
		state = c.conn.GetState()
		<-time.After(time.Millisecond * 100)
	}
}
func (c *Client) ConnectionStatus() connectivity.State {
	state := c.conn.GetState()
	return state

}
func (c *Client) Disconnect() error {

	return c.conn.Close()
}
func (c *Client) Connect() error {
	if c.conn == nil {
		addr := c.addr
		isRemote := c.isRemote
		var err error
		dialer := func(context.Context, string) (net.Conn, error) {
			if isRemote {
				return net.Dial("tcp", addr)
			}
			return net.Dial("unix", addr)
		}
		dd := grpc.WithContextDialer(dialer)
		var kacp = keepalive.ClientParameters{
			Time:                4 * time.Second, // send pings every 4 seconds if there is no activity
			Timeout:             time.Second * 2, // wait 2 second for ping ack before considering the connection dead
			PermitWithoutStream: true,            // send pings even without active streams
		}

		c.conn, err = grpc.Dial(addr, dd, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
		if err != nil {
			fmt.Println(err)
			return err
		}

		c.client = NewPluginProtocolClient(c.conn)
	}

	c.conn.Connect()
	return nil
}
func NewClient(addr string, remote bool, ctx context.Context, logger *zap.Logger) *Client {
	if logger == nil {
		logger = helper.LoggerNamed("logs/client.log")
	}

	baseCli := &Client{
		addr:     addr,
		ctx:      ctx,
		isRemote: remote,
		logger:   logger,
	}
	return baseCli
}

type InItem struct {
	Index    int
	Type     string
	BaseData interface{}
}

func (ii *InItem) Populate() *InTypes {

	return &InTypes{}
}
func (cl *Client) GetConfig() ([]byte, error) {
	req := &Empty{}
	t := time.Now()

	resp, err := cl.client.RequestConfig(cl.ctx, req)
	if err != nil {
		return nil, err
	}
	cl.responseTime = time.Since(t)
	if resp.GetSuccess() {
		return resp.GetData(), nil
	} else {
		return resp.GetData(), yaperror.Error(yaperror.CONFIG_GET_ERROR, nil)
	}

}
func (cl *Client) Call(fncname string, strname string, params []InItem) (interface{}, error) {
	var ins []*InTypes
	t := time.Now()
	for _, itm := range params {
		data, _ := json.Marshal(itm.BaseData)
		ins = append(ins, &InTypes{
			Index: int32(itm.Index),
			In:    data,
			Type:  itm.Type,
		})
	}
	req := &FunctionRequest{
		In:       ins,
		Function: fncname,
		Struct:   strname,
	}
	var respItem interface{}
	resp, err := cl.client.CallFunction(cl.ctx, req)
	if err != nil {
		fmt.Println("Error", err)
		return respItem, err
	}
	cl.responseTime = time.Since(t)
	json.Unmarshal(resp.GetData(), &respItem)

	var errN error
	if resp.GetError() != nil {
		errN = errors.New(resp.GetError().Message)
	}
	return respItem, errN
}

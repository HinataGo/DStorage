// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: dbproxy.proto

package go_micro_service_proxy

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for DBProxyService service

type DBProxyService interface {
	// 请求执行sql动作
	ExecuteAction(ctx context.Context, in *ReqExec, opts ...client.CallOption) (*RespExec, error)
}

type dBProxyService struct {
	c    client.Client
	name string
}

func NewDBProxyService(name string, c client.Client) DBProxyService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.service.dbproxy"
	}
	return &dBProxyService{
		c:    c,
		name: name,
	}
}

func (c *dBProxyService) ExecuteAction(ctx context.Context, in *ReqExec, opts ...client.CallOption) (*RespExec, error) {
	req := c.c.NewRequest(c.name, "DBProxyService.ExecuteAction", in)
	out := new(RespExec)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DBProxyService service

type DBProxyServiceHandler interface {
	// 请求执行sql动作
	ExecuteAction(context.Context, *ReqExec, *RespExec) error
}

func RegisterDBProxyServiceHandler(s server.Server, hdlr DBProxyServiceHandler, opts ...server.HandlerOption) error {
	type dBProxyService interface {
		ExecuteAction(ctx context.Context, in *ReqExec, out *RespExec) error
	}
	type DBProxyService struct {
		dBProxyService
	}
	h := &dBProxyServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&DBProxyService{h}, opts...))
}

type dBProxyServiceHandler struct {
	DBProxyServiceHandler
}

func (h *dBProxyServiceHandler) ExecuteAction(ctx context.Context, in *ReqExec, out *RespExec) error {
	return h.DBProxyServiceHandler.ExecuteAction(ctx, in, out)
}

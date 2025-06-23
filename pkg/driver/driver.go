package driver

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

const (
	DefaultName = "ccp.sai.dev" // ccp = customcsiplugin
)

type Driver struct {
	name     string
	endpoint string
	region   string
	token    string
	csi.UnimplementedNodeServer
	csi.UnimplementedControllerServer
	csi.UnimplementedIdentityServer
	srv *grpc.Server
}

type InputParams struct {
	Name     string
	Endpoint string
	Region   string
	Token    string
}

func NewDriver(params InputParams) *Driver {
	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
		region:   params.Region,
		token:    params.Token,
	}
}

func (d *Driver) Run() error {
	url, err := url.Parse(d.endpoint)
	if err != nil {
		return fmt.Errorf("error parsing the url %s", err.Error())
	}
	if url.Scheme != "unix" {
		return fmt.Errorf("expected unix endpoint but received %s", url.Scheme)
	}
	grpcAddress := path.Join(url.Host, filepath.FromSlash(url.Path))
	if url.Host == "" {
		grpcAddress = filepath.FromSlash(url.Path)
	}

	// os.remove try to remove the file and return the error
	// our goal is to delete the file, if it exists and ignore if it is not
	// when we do remove we can get many errors, if err is file not exists
	// then we can proceed
	// if we use os.IsExist it will return true, if error indicates file exists
	// os.IsNotExist will return true, if error indicates file does not exist
	// for first time file might not be available and os.isExist will return false
	// and !os.IsExist() will make it true, which is not desired

	if err := os.Remove(grpcAddress); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing listen address %s", err.Error())
	}
	// how to start the server
	// https://grpc.io/docs/languages/go/basics/#starting-the-server

	listener, err := net.Listen(url.Scheme, grpcAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err.Error())
	}
	fmt.Println(listener)
	d.srv = grpc.NewServer()

	// https://github.com/container-storage-interface/spec/blob/master/lib/go/csi/csi_grpc.pb.go#L1334C1-L1334C5

	csi.RegisterNodeServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterIdentityServer(d.srv, d)

	return d.srv.Serve(listener)
}

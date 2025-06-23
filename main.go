package main

import (
	"flag"

	"github.com/kama001/customcsiplugin/pkg/driver"
	"k8s.io/klog"
	// "k8s.io/klog/v2"
)

func main() {
	endpoint := flag.String("endpoint", "default value", "endpoint of gRPC server")
	token := flag.String("token", "default value", "token of storage provider")
	region := flag.String("region", "eu-central-1", "region where volumes are going to be provisioned")

	flag.Parse()

	drv := driver.NewDriver(driver.InputParams{
		Name:     driver.DefaultName,
		Endpoint: *endpoint,
		Token:    *token,
		Region:   *region,
	})

	if err := drv.Run(); err != nil {
		klog.Errorf("error running the driver %s\n", err.Error())
	}
}

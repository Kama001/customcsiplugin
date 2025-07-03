package driver

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Driver) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.
	ControllerGetCapabilitiesResponse, error) {
	caps := []*csi.ControllerServiceCapability{}

	for _, c := range []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	} {
		caps = append(caps, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: c,
				},
			},
		})
	}
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (d *Driver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// get name (for idempotency as well)
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Create volume should be called with a required name")
	}

	// extract memory requirements
	const gb = 1024 * 1024 * 1024
	// fmt.Println("requested capacity:", req.CapacityRange.GetRequiredBytes())
	sizeGB := int64(math.Ceil(float64(req.CapacityRange.GetRequiredBytes()) / float64(gb)))
	if sizeGB < 1 {
		sizeGB = 1
	}

	// get capabilities
	if req.GetVolumeCapabilities() == nil || len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities has not beend specified!!")
	}
	// validate volume capabilities
	// make sure Access mode: block or file system specified inside PVC is suported by SP
	// make sure Access type: readwritebymany or something is supported by SP

	// validate credentials
	if d.access_key == "" || d.secret_key == "" || d.token == "" {
		return nil, status.Error(codes.InvalidArgument, "AWS credentials not provided in secrets")
	}

	// create AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(d.region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(d.access_key, d.secret_key, d.token),
		),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load AWS config: %v", err)
	}

	// create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// create the request struct
	out, err := ec2Client.CreateVolume(ctx, &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("eu-central-1a"),
		Size:             aws.Int32(int32(sizeGB)),
	})

	fmt.Println("allocated capacity in bytes", sizeGB*gb)
	fmt.Println("allocated capacity in GBs", sizeGB)
	fmt.Println("creating volume.......")
	time.Sleep(5 * time.Second)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create EBS volume: %v", err)
	}

	// response
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      *out.VolumeId,
			CapacityBytes: int64(sizeGB * gb),
		},
	}, nil
}

func (d *Driver) ControllerPublishVolume(context.Context, *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	fmt.Println("controller publish volume is called!!!")
	return nil, nil
}

// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package batch

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	batch_sdkv2 "github.com/aws/aws-sdk-go-v2/service/batch"
	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	batch_sdkv1 "github.com/aws/aws-sdk-go/service/batch"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory: newJobDefinitionDataSource,
			Name:    "Job Definition",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceJobQueue,
			Name:    "Job Queue",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceComputeEnvironment,
			TypeName: "aws_batch_compute_environment",
		},
		{
			Factory:  DataSourceJobQueue,
			TypeName: "aws_batch_job_queue",
		},
		{
			Factory:  DataSourceSchedulingPolicy,
			TypeName: "aws_batch_scheduling_policy",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceComputeEnvironment,
			TypeName: "aws_batch_compute_environment",
			Name:     "Compute Environment",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceJobDefinition,
			TypeName: "aws_batch_job_definition",
			Name:     "Job Definition",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceSchedulingPolicy,
			TypeName: "aws_batch_scheduling_policy",
			Name:     "Scheduling Policy",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Batch
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*batch_sdkv1.Batch, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return batch_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config[names.AttrEndpoint].(string))})), nil
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*batch_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return batch_sdkv2.NewFromConfig(cfg, func(o *batch_sdkv2.Options) {
		if endpoint := config[names.AttrEndpoint].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}

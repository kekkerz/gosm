package ec2

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"reflect"
	"strconv"
	"testing"
)

type mockDescribeInstancesAPI func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

func (m mockDescribeInstancesAPI) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return m(ctx, params, optFns...)
}

func TestDescribeInstances(t *testing.T) {
	cases := []struct {
		client  func(t *testing.T) EC2DescribeInstancesAPI
		filters FilterOptions
		expect  []types.Reservation
	}{
		{
			client: func(t *testing.T) EC2DescribeInstancesAPI {
				return mockDescribeInstancesAPI(func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
					t.Helper()
					if params == nil {
						t.Fatal("expect filters to not be nil")
					}

					if e, a := "bitwarden", params.Filters[0].Values[0]; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &ec2.DescribeInstancesOutput{
						Reservations: []types.Reservation{
							{
								Instances: []types.Instance{
									{
										InstanceId: aws.String("i-xxxx"),
									},
								},
							},
						},
					}, nil
				})
			},
			filters: FilterOptions{
				Name: "bitwarden",
			},
			expect: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-xxxx"),
						},
					},
				},
			},
		},
		{
			client: func(t *testing.T) EC2DescribeInstancesAPI {
				return mockDescribeInstancesAPI(func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
					t.Helper()
					if params == nil {
						t.Fatal("expect filters to not be nil")
					}

					if e, a := "i-xxxx", params.InstanceIds[0]; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &ec2.DescribeInstancesOutput{
						Reservations: []types.Reservation{
							{
								Instances: []types.Instance{
									{
										InstanceId: aws.String("i-xxxx"),
									},
								},
							},
						},
					}, nil
				})
			},
			filters: FilterOptions{
				InstanceId: "i-xxxx",
			},
			expect: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-xxxx"),
						},
					},
				},
			},
		},
		{
			client: func(t *testing.T) EC2DescribeInstancesAPI {
				return mockDescribeInstancesAPI(func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
					t.Helper()
					if params == nil {
						t.Fatal("expect filters to not be nil")
					}

					e := &[]types.Filter{}
					err := json.Unmarshal([]byte(`[{"Name": "tag:Name", "Values": ["bitwarden", "tailscale"]}]`), &e)
					if err != nil {
						t.Fatal("Error unmarshaling JSON: ", err)
					}

					if a := params.Filters; reflect.DeepEqual(e, a) {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &ec2.DescribeInstancesOutput{
						Reservations: []types.Reservation{
							{
								Instances: []types.Instance{
									{
										InstanceId: aws.String("i-xxxx"),
									},
								},
							},
						},
					}, nil
				})
			},
			filters: FilterOptions{
				Tags: `[{"Name": "tag:Name", "Values": ["bitwarden", "tailscale"]}]`,
			},
			expect: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId: aws.String("i-xxxx"),
						},
					},
				},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			instances := GetInstanceMetaData(tt.client(t), tt.filters)
			e := aws.ToString(tt.expect[0].Instances[0].InstanceId)
			a := aws.ToString(instances[0].Instances[0].InstanceId)
			if e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

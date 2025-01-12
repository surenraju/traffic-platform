package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestVpcModule(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform-modules/vpc",
		Vars: map[string]interface{}{
			"cidr_block": "10.0.0.0/16",
			"vpc_name":   "test-vpc",
			"region":     region,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get VPC ID
	vpcID := terraform.Output(t, terraformOptions, "vpc_id")
	assert.NotEmpty(t, vpcID, "VPC ID should not be empty")

	// Get Subnets
	subnets := getSubnetsForVpc(t, vpcID, region)
	assert.Len(t, subnets, 8, "VPC should have 8 subnets")

	// Classify Subnets
	publicSubnets := []string{}
	privateSubnets := []string{}
	tgwSubnets := []string{}

	for _, subnet := range subnets {
		routeTable := getRouteTableForSubnet(t, subnet.SubnetId, region)

		if hasInternetGatewayRoute(routeTable) {
			publicSubnets = append(publicSubnets, *subnet.SubnetId)
		} else if hasNatGatewayRoute(routeTable) {
			privateSubnets = append(privateSubnets, *subnet.SubnetId)
		} else if hasTgwAttachmentTag(subnet) {
			tgwSubnets = append(tgwSubnets, *subnet.SubnetId)
		}
	}

	// Assertions for public subnets
	assert.Len(t, publicSubnets, 2, "There should be 2 public subnets with routes to the Internet Gateway")

	// Assertions for private subnets
	assert.Len(t, privateSubnets, 4, "There should be 4 private subnets with routes to the NAT Gateway")

	// Assertions for Transit Gateway subnets
	assert.Len(t, tgwSubnets, 2, "There should be 2 subnets tagged for Transit Gateway attachments")
	for _, subnetID := range tgwSubnets {
		validateTgwSubnetTags(t, region, subnetID, map[string]string{
			"TransitGatewayAttachment": "true",
		})
	}
}

func getSubnetsForVpc(t *testing.T, vpcID, region string) []*ec2.Subnet {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to get subnets for VPC %s: %v", vpcID, err)
	}
	return output.Subnets
}

func getRouteTableForSubnet(t *testing.T, subnetID *string, region string) *ec2.RouteTable {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{subnetID},
			},
		},
	})
	if err != nil || len(output.RouteTables) == 0 {
		t.Fatalf("Failed to get route table for subnet %s: %v", *subnetID, err)
	}
	return output.RouteTables[0]
}

func validateTgwSubnetTags(t *testing.T, region, subnetID string, expectedTags map[string]string) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: []*string{aws.String(subnetID)},
	})
	if err != nil || len(output.Subnets) == 0 {
		t.Fatalf("Failed to get subnet %s for tag validation: %v", subnetID, err)
	}
	subnet := output.Subnets[0]

	for key, expectedValue := range expectedTags {
		actualValue := getTagValue(subnet.Tags, key)
		assert.Equal(t, expectedValue, actualValue, "Subnet %s should have tag %s with value %s", subnetID, key, expectedValue)
	}
}

func getTagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func hasInternetGatewayRoute(routeTable *ec2.RouteTable) bool {
	for _, route := range routeTable.Routes {
		if route.GatewayId != nil && *route.GatewayId != "" {
			return true
		}
	}
	return false
}

func hasNatGatewayRoute(routeTable *ec2.RouteTable) bool {
	for _, route := range routeTable.Routes {
		if route.NatGatewayId != nil && *route.NatGatewayId != "" {
			return true
		}
	}
	return false
}

func hasTgwAttachmentTag(subnet *ec2.Subnet) bool {
	for _, tag := range subnet.Tags {
		if *tag.Key == "TransitGatewayAttachment" && *tag.Value == "true" {
			return true
		}
	}
	return false
}

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
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

	// Validate Subnets
	subnets := aws.GetSubnetsForVpc(t, region, vpcID)
	assert.Len(t, subnets, 8, "VPC should have 8 subnets")

	// Classify Subnets
	publicSubnets := []string{}
	privateSubnets := []string{}
	tgwSubnets := []string{}

	for _, subnet := range subnets {
		routeTable := aws.GetRouteTableForSubnet(t, subnet.Id, region)

		if hasInternetGatewayRoute(routeTable) {
			publicSubnets = append(publicSubnets, subnet.Id)
		} else if hasNatGatewayRoute(routeTable) {
			privateSubnets = append(privateSubnets, subnet.Id)
		} else if hasTgwAttachmentTag(subnet) {
			tgwSubnets = append(tgwSubnets, subnet.Id)
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

func validateTgwSubnetTags(t *testing.T, region, subnetID string, expectedTags map[string]string) {
	subnet := aws.GetSubnetById(t, region, subnetID)
	for key, expectedValue := range expectedTags {
		actualValue, exists := subnet.Tags[key]
		assert.True(t, exists, "TGW subnet "+subnetID+" should have the tag "+key)
		assert.Equal(t, expectedValue, actualValue, "TGW subnet "+subnetID+" should have tag "+key+" with value "+expectedValue)
	}
}

func hasInternetGatewayRoute(routeTable aws.RouteTable) bool {
	for _, route := range routeTable.Routes {
		if route.GatewayId != nil && *route.GatewayId != "" {
			return true
		}
	}
	return false
}

func hasNatGatewayRoute(routeTable aws.RouteTable) bool {
	for _, route := range routeTable.Routes {
		if route.NatGatewayId != nil && *route.NatGatewayId != "" {
			return true
		}
	}
	return false
}

func hasTgwAttachmentTag(subnet aws.Subnet) bool {
	return subnet.Tags["TransitGatewayAttachment"] == "true"
}

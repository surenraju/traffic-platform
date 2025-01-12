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
			"cidr_block":  "10.0.0.0/16",
			"vpc_name":    "test-vpc",
			"region":      region,
			"environment": "prod",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get VPC ID
	vpcID := terraform.Output(t, terraformOptions, "vpc_id")
	assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
	t.Logf("VPC ID: %s", vpcID)

	// Get Subnets
	subnets := getSubnetsForVpc(t, vpcID, region)
	assert.Len(t, subnets, 8, "VPC should have 8 subnets")
	t.Logf("Subnets: %v", subnets)

	// Classify Subnets
	publicSubnets := []string{}
	privateSubnets := []string{}
	tgwSubnets := []string{}

	for _, subnet := range subnets {
		t.Logf("Checking Subnet: %s", *subnet.SubnetId)

		if hasTgwAttachmentTag(subnet) {
			t.Logf("Subnet %s classified as TGW attachment", *subnet.SubnetId)
			tgwSubnets = append(tgwSubnets, *subnet.SubnetId)
			continue // Skip route table checks for TGW subnets
		}

		routeTable := getRouteTableForSubnet(t, subnet.SubnetId, region)
		if routeTable == nil {
			t.Logf("Subnet %s has no associated route table", *subnet.SubnetId)
			continue
		}

		t.Logf("Route Table for Subnet %s: %+v", *subnet.SubnetId, routeTable)

		if hasInternetGatewayRoute(routeTable) && *subnet.MapPublicIpOnLaunch {
			t.Logf("Subnet %s classified as public", *subnet.SubnetId)
			publicSubnets = append(publicSubnets, *subnet.SubnetId)
		} else if hasNatGatewayRoute(routeTable) {
			t.Logf("Subnet %s classified as private", *subnet.SubnetId)
			privateSubnets = append(privateSubnets, *subnet.SubnetId)
		} else {
			t.Logf("Subnet %s could not be classified", *subnet.SubnetId)
		}
	}

	// Assertions for public subnets
	t.Log("Validating public subnets...")
	assert.Len(t, publicSubnets, 2, "There should be 2 public subnets with routes to the Internet Gateway")
	t.Logf("Public Subnets: %v", publicSubnets)

	for _, subnet := range publicSubnets {
		t.Logf("Validating route table for public subnet %s", subnet)
		routeTable := getRouteTableForSubnet(t, aws.String(subnet), region)
		assert.NotNil(t, routeTable, "Public subnet %s should have an associated route table", subnet)
		assert.True(t, hasInternetGatewayRoute(routeTable), "Public subnet %s should have a route to the Internet Gateway", subnet)
	}

	// Assertions for private subnets
	t.Log("Validating private subnets...")
	assert.Len(t, privateSubnets, 4, "There should be 4 private subnets with routes to the NAT Gateway")
	t.Logf("Private Subnets: %v", privateSubnets)

	for _, subnet := range privateSubnets {
		t.Logf("Validating route table for private subnet %s", subnet)
		routeTable := getRouteTableForSubnet(t, aws.String(subnet), region)
		assert.NotNil(t, routeTable, "Private subnet %s should have an associated route table", subnet)
		assert.True(t, hasNatGatewayRoute(routeTable), "Private subnet %s should have a route to the NAT Gateway", subnet)
	}

	// Assertions for TGW subnets
	t.Log("Validating Transit Gateway subnets...")
	assert.Len(t, tgwSubnets, 2, "There should be 2 subnets tagged for Transit Gateway attachments")
	t.Logf("TGW Subnets: %v", tgwSubnets)

	for _, subnetID := range tgwSubnets {
		t.Logf("Validating tags for TGW Subnet: %s", subnetID)
		validateTgwSubnetTags(t, region, subnetID, map[string]string{
			"TransitGatewayAttachment": "true",
		})
	}
	t.Log("Test completed successfully.")
}

func getSubnetsForVpc(t *testing.T, vpcID, region string) []*ec2.Subnet {
	t.Logf("Fetching subnets for VPC ID: %s", vpcID)
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
	t.Logf("Fetching route table for Subnet ID: %s", *subnetID)
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
		t.Logf("No route table found for Subnet ID: %s (this is expected for TGW subnets)", *subnetID)
		return nil // Return nil if no route table is found
	}
	return output.RouteTables[0]
}

func hasInternetGatewayRoute(routeTable *ec2.RouteTable) bool {
	if routeTable == nil {
		return false
	}
	for _, route := range routeTable.Routes {
		if route.GatewayId != nil && *route.GatewayId != "" {
			return true
		}
	}
	return false
}

func hasNatGatewayRoute(routeTable *ec2.RouteTable) bool {
	if routeTable == nil {
		return false
	}
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

func validateTgwSubnetTags(t *testing.T, region, subnetID string, expectedTags map[string]string) {
	t.Logf("Validating tags for Subnet ID: %s", subnetID)
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: []*string{aws.String(subnetID)},
	})
	if err != nil || len(output.Subnets) == 0 {
		t.Fatalf("Failed to get subnet %s for tag validation: %v", subnetID, err)
	}
	subnet := output.Subnets[0]

	t.Logf("Retrieved tags for Subnet %s: %v", subnetID, subnet.Tags)

	for key, expectedValue := range expectedTags {
		actualValue := getTagValue(subnet.Tags, key)
		if expectedValue != actualValue {
			t.Logf("Tag validation is going to fail. Key: %s, Expected: %s, Actual: %s", key, expectedValue, actualValue)
		}
		t.Logf("Validating tag: %s, Expected: %s, Actual: %s", key, expectedValue, actualValue)
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

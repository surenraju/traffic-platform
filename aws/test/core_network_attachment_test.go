package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCoreNetworkAttachmentModule(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	// Create VPC and Transit Gateway first
	vpcOptions := &terraform.Options{
		TerraformDir: "../terraform-modules/vpc",
		Vars: map[string]interface{}{
			"cidr_block": "10.0.0.0/16",
			"vpc_name":   "test-vpc",
			"region":     region,
		},
	}
	defer terraform.Destroy(t, vpcOptions)
	terraform.InitAndApply(t, vpcOptions)

	tgwOptions := &terraform.Options{
		TerraformDir: "../terraform-modules/transit-gateway",
		Vars: map[string]interface{}{
			"description":     "Test Transit Gateway",
			"amazon_side_asn": 64512,
			"name":            "test-tgw",
		},
	}
	defer terraform.Destroy(t, tgwOptions)
	terraform.InitAndApply(t, tgwOptions)

	// Extract outputs
	vpcID := terraform.Output(t, vpcOptions, "vpc_id")
	tgwID := terraform.Output(t, tgwOptions, "transit_gateway_id")
	subnetIDs := terraform.OutputList(t, vpcOptions, "tgw_attachment_subnet_ids")
	routeTableIDs := terraform.OutputList(t, vpcOptions, "private_route_table_ids")

	// Test Core Network Attachment
	coreOptions := &terraform.Options{
		TerraformDir: "../terraform-modules/core-network-attachment",
		Vars: map[string]interface{}{
			"vpc_id":             vpcID,
			"transit_gateway_id": tgwID,
			"subnet_ids":         subnetIDs,
			"route_table_ids":    routeTableIDs,
		},
	}
	defer terraform.Destroy(t, coreOptions)
	terraform.InitAndApply(t, coreOptions)

	// Verify TGW Attachment
	for _, subnetID := range subnetIDs {
		attachment := getTransitGatewayAttachmentForSubnet(t, region, tgwID, subnetID)
		assert.NotNil(t, attachment, "Transit Gateway Attachment should exist for subnet "+subnetID)
	}

	// Verify Routes
	for _, routeTableID := range routeTableIDs {
		routeTable := getRouteTable(t, region, routeTableID)
		assertRouteExists(t, routeTable, "10.0.0.0/8", tgwID)
	}
}

func getTransitGatewayAttachmentForSubnet(t *testing.T, region, tgwID, subnetID string) *ec2.TransitGatewayVpcAttachment {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeTransitGatewayVpcAttachments(&ec2.DescribeTransitGatewayVpcAttachmentsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("transit-gateway-id"),
				Values: []*string{aws.String(tgwID)},
			},
			{
				Name:   aws.String("subnet-id"),
				Values: []*string{aws.String(subnetID)},
			},
		},
	})
	if err != nil || len(output.TransitGatewayVpcAttachments) == 0 {
		t.Fatalf("Failed to find Transit Gateway VPC Attachment for subnet %s: %v", subnetID, err)
	}
	return output.TransitGatewayVpcAttachments[0]
}

func getRouteTable(t *testing.T, region, routeTableID string) *ec2.RouteTable {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(routeTableID)},
	})
	if err != nil || len(output.RouteTables) == 0 {
		t.Fatalf("Failed to find route table %s: %v", routeTableID, err)
	}
	return output.RouteTables[0]
}

func assertRouteExists(t *testing.T, routeTable *ec2.RouteTable, destinationCidrBlock, tgwID string) {
	found := false
	for _, route := range routeTable.Routes {
		if *route.DestinationCidrBlock == destinationCidrBlock && route.TransitGatewayId != nil && *route.TransitGatewayId == tgwID {
			found = true
			break
		}
	}
	assert.True(t, found, "Route table should have a route to %s via Transit Gateway %s", destinationCidrBlock, tgwID)
}

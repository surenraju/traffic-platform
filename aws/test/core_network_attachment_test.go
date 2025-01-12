package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
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
		attachment := aws.GetTransitGatewayVpcAttachment(t, region, map[string]string{
			"subnet-id": subnetID,
		})
		assert.NotNil(t, attachment, "Transit Gateway Attachment should exist for subnet "+subnetID)
	}

	// Verify Routes
	for _, routeTableID := range routeTableIDs {
		routes := aws.GetRouteTable(t, region, routeTableID)
		assert.Contains(t, routes.Routes, "10.0.0.0/8", "Route table should have route to 10.0.0.0/8 via Transit Gateway")
	}
}

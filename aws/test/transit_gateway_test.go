package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTransitGatewayModule(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform-modules/transit-gateway",
		Vars: map[string]interface{}{
			"description":     "Test Transit Gateway",
			"amazon_side_asn": 64512,
			"name":            "test-tgw",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get Transit Gateway ID
	tgwID := terraform.Output(t, terraformOptions, "transit_gateway_id")
	assert.NotEmpty(t, tgwID, "Transit Gateway ID should not be empty")

	// Validate Transit Gateway exists in AWS
	tgw := aws.GetTransitGatewayById(t, tgwID, region)
	assert.NotNil(t, tgw, "Transit Gateway should exist")
	assert.Equal(t, tgw.AmazonSideAsn, int64(64512), "Transit Gateway ASN should match")
}

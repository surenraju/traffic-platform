package test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTransitGatewayModule(t *testing.T) {
	t.Parallel()
	_, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

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
	tgw := getTransitGatewayByID(t, tgwID, region)
	assert.NotNil(t, tgw, "Transit Gateway should exist")
	assert.Equal(t, int64(64512), *tgw.Options.AmazonSideAsn, "Transit Gateway ASN should match")
	assert.Equal(t, "Test Transit Gateway", *tgw.Description, "Transit Gateway description should match")
}

func getTransitGatewayByID(t *testing.T, tgwID string, region string) *ec2.TransitGateway {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)

	output, err := ec2Client.DescribeTransitGateways(&ec2.DescribeTransitGatewaysInput{
		TransitGatewayIds: []*string{aws.String(tgwID)},
	})
	if err != nil || len(output.TransitGateways) == 0 {
		t.Fatalf("Failed to get Transit Gateway %s: %v", tgwID, err)
	}
	return output.TransitGateways[0]
}

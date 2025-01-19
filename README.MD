# Automating a Basic Control Plane for Networking as a Service

This reposittory demonstrates a proof of concept (POC) control plane for Networking as a Service using automation tools like Terraform and custom scripts. While this POC is based on AWS, the approach can be extended to other cloud providers such as Azure and GCP. The POC provides two primary capabilities:

- **VPC as a Service**

This feature provisions VPCs with an optimized AZ and subnet design. Users can choose a VPC template, and the control plane provisions the VPC while applying tags for seamless integration into global networking.

This feature provisions VPCs with an optimized AZ and subnet design. Users can choose a VPC template, and the control plane provisions the VPC while applying tags for seamless integration into global networking.

- **Global Networking as a Service**

The control plane automatically discovers VPCs tagged for networking, attaches them to a Transit Gateway, and configures routing. The solution can be extended to include additional capabilities like traffic inspection, Route 53 hosted zones, VPN connectivity, multi-region extensions, and multi-cloud integrations.

## Steps to Implement the Control Plane POC

### Clone the Repository

The automation code is hosted on GitHub. Clone the repository to your local machine:

```sh
git clone https://github.com/surenraju/traffic-platform
cd traffic-platform
```

### Set Up the Environment

Build a Docker environment with required dependencies to run the automation code.

```sh
docker build . -t docker-terragrunt-local:latest
docker run -it -v $(pwd):/traffic-platform docker-terragrunt-local:latest /bin/bash
```

For simplicity, this POC uses the devopsinfra/docker-terragrunt image. In production, a CI/CD system like GitHub Actions can replace this setup.

### Terraform Modules in the POC

- **aws/terraform-modules/vpc**: Provisions VPCs with configurable AZ and subnet designs.
- **aws/terraform-modules/transit-gateway**: Automates the provisioning of AWS Transit Gateway.
- **aws/terraform-modules/core-network-attachment**: Attaches VPCs to the Transit Gateway.
- **aws/terraform-modules/reachability-analyzer**: Configures AWS Reachability Analyzer to verify network connectivity.

#### Lifecycle Hooks to Prevent Destructive Actions

To protect critical resources, lifecycle hooks are used:

```sh
lifecycle {
    prevent_destroy = true
}
```

### Testing Terraform Modules

The POC uses the `terratest` framework to validate individual modules and their interactions. Run the following commands from the `aws` folder.

```sh
go mod init traffic-platform
go mod tidy
go test -v ./test
```

### VPC Provisioning

Users provide VPC configurations in a Kubernetes Resource Model (KRM) like YAML format:

```yaml
# config/vpc1.yaml
apiVersion: trafficplatform.aws/v1
kind: VPC
metadata:
  name: vpc1
spec:
  account_id: "acc1"
  vpc_name: vpc1
  cidr_block: 10.10.0.0/16
  region: us-east-1
  environment: prod
```

Transform the configuration into Terraform input by running:

```sh
pip install -r resources/requirements.txt
python3 resources/main.py --resource=VPC
```

#### Validation The Plan

To ensure that only valid and intended changes reach production, `Conftest` policies can be implemented in the pipeline. These policies act as guardrails, detecting and blocking unauthorized or destructive actions such as resource deletions.

Here’s an example `Conftest` policy that prevents deletion operations:

```rego
package main

deny[msg] {
    resource := input.resource_changes[_]
    actions := resource.change.actions
    "delete" == actions[_]
    msg = sprintf("Delete operation detected for resource: %s", [resource.address])
}
```

How It Works?

- The policy inspects the resource_changes in the `Terraform plan` output.
- If a "delete" action is detected for any resource, the policy generates a denial message.
- This ensures no critical resource is accidentally deleted during automation.

```sh
terragrunt run-all plan --terragrunt-json-out-dir=output --terragrunt-working-dir=live/acc1/VPC
conftest test live/acc1/VPC/output/tfplan.json --policy policy/
```

#### Apply the Terraform plan to provision VPCs

```sh
terragrunt run-all apply --terragrunt-working-dir=live/acc1/VPC
```

### Creating the Transit Gateway as Part of the Control Plane

As part of the control plane automation, the Transit Gateway (TGW) must be provisioned to enable seamless global networking. This step ensures VPCs and other network components are attached to the TGW for centralized management. Follow these steps:

Run the following Terragrunt command to apply the Terraform configuration for the Transit Gateway setup:

```sh
# Run the pan and conftest policies
terragrunt run-all plan --terragrunt-json-out-dir=output --terragrunt-working-dir=live/acc1/TransitGateway
conftest test live/acc1/TransitGateway/output/tfplan.json --policy policy/
# Run terragrunt apply
terragrunt run-all apply --terragrunt-working-dir live/acc1/TransitGateway
```

### Global Networking with Transit Gateway

Lets configure the Transit Gateway for global networking. The setup requires initial configuration:

```yaml
# config/tgw.yaml
apiVersion: trafficplatform.aws/v1
kind: TransitGateway
metadata:
  name: tgw-01
spec:
  transit_gateways:
  - region: us-east-1
    id: tgw-008a46bb27ebdc169
```

#### Input for Global Networking as a Service

To enable Global Networking as a Service in the control plane, users need to define a configuration for attaching VPCs to the Transit Gateway. This is achieved by creating a YAML configuration file in the Kubernetes Resource Model (KRM) format, which specifies the account and region for discovering VPCs.

Users define the VPCs to be attached to the global network in the following YAML file:

```yaml
# config/core-net-att.yaml
apiVersion: trafficplatform.aws/v1
kind: CoreNetworkAttachment
metadata:
  name: acc1-core-network-attachment
spec:
  account_id: "acc1"
  region: "us-east-1"

```

#### Transform the Configuration into Terraform Inputs

Run the following command to convert the KRM model into Terraform inputs:

```sh
# Transform the TransitGateway KRM model to terragrunt.hcl
python3 resources/main.py --resource=CoreNetworkAttachment
```

This step performs the following tasks:

- Discovers VPCs, subnets, and route tables tagged for Transit Gateway attachment.
- Converts the input YAML into Terraform-compatible configurations.

#### Validate the Terraform Plan

Before applying the configuration, generate and validate the Terraform plan:

```sh
terragrunt run-all plan --terragrunt-json-out-dir=output --terragrunt-working-dir=live/acc1/CoreNetworkAttachment
conftest test live/acc1/CoreNetworkAttachment/output/tfplan.json --policy policy/
```

#### Apply the Terraform Configuration

Finally, apply the configuration to attach the discovered VPCs to the Transit Gateway:

```sh
terragrunt run-all apply --terragrunt-working-dir=live/acc1/CoreNetworkAttachment
```

This operation:
  - Attaches the tagged VPCs to the Transit Gateway
  - Configures route tables to enable seamless global connectivity

#### Validation with Reachability Analyzer

To ensure network connectivity remains intact during automation, the control plane integrates AWS Reachability Analyzer for validation. This tool tests connectivity between VPCs and verifies that the core network configurations, such as Transit Gateway attachments and route tables, are functioning as expected.

```sh
terragrunt run-all apply --terragrunt-working-dir=live/acc1/ReachabilityAnalyzer
```

Integrate the test results into your monitoring and alerting system to ensure immediate action on failures:
  - Use tools like Amazon CloudWatch Alarms or third-party solutions to notify the Networking SRE team when a test fails.
  - Automate a rollback or additional diagnostics in response to detected issues.
# Automating a Basic Control Plane for Networking as a Service

This reposittory demonstrates a proof of concept (POC) control plane for Networking as a Service using automation tools like Terraform and custom scripts. While this POC is based on AWS, the approach can be extended to other cloud providers such as Azure and GCP. The POC provides two primary capabilities:

- VPC as a Service:
This feature provisions VPCs with an optimized AZ and subnet design. Users can choose a VPC template, and the control plane provisions the VPC while applying tags for seamless integration into global networking.
- Global Networking as a Service:
The control plane automatically discovers VPCs tagged for networking, attaches them to a Transit Gateway, and configures routing. The solution can be extended to include additional capabilities like traffic inspection, Route 53 hosted zones, VPN connectivity, multi-region extensions, and multi-cloud integrations.












docker run -it -v $(pwd):/traffic-platform-root docker-terragrunt-local:latest /bin/bash

git clone https://{TOKEN}@github.com/surenraju/traffic-platform.git
git config --global email



export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"

cd aws
pip install -r resources/requirements.txt
python3 resources/main.py --resource=VPC
python3 resources/main.py --resource=CoreNetworkAttachment

terragrunt run-all init
terragrunt run-all plan
terragrunt run-all apply

terragrunt run-all apply --terragrunt-working-dir live/acc1/VPC
terragrunt run-all apply --terragrunt-working-dir live/acc1/TransitGateway
terragrunt run-all apply --terragrunt-working-dir live/acc1/CoreNetworkAttachment

terragrunt run-all apply --terragrunt-working-dir live/acc1/ReachabilityAnalyzer

terragrunt run-all plan --terragrunt-json-out-dir=output --terragrunt-working-dir=live/acc1/VPC

Validation:

conftest test live/acc1/VPC/output/tfplan.json --policy policy/

Test:

go mod init traffic-platform
go mod tidy
go test -v ./test



Debug
export TF_PLUGIN_TIMEOUT=5m
rm -rf ~/.terraform.d/plugins
rm -rf .terraform
rm -rf .terragrunt-cache

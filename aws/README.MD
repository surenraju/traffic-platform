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

terragrunt run-all plan \
  --terragrunt-json-out-dir=output \
  --terragrunt-working-dir=live/acc1/VPC

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
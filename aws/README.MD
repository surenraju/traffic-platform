docker run -it -v $(pwd):/traffic-platform-copy devopsinfra/docker-terragrunt:aws-latest /bin/bash

git clone https://{TOKEN}@github.com/surenraju/traffic-platform.git
git config --global email

cd aws
pip install -r resources/requirements.txt
python resources/main.py

export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"

terragrunt run-all init
terragrunt run-all plan
terragrunt run-all apply

cd ../vpc1
terragrunt init
terragrunt plan
terragrunt apply

cd ../vpc2
terragrunt init
terragrunt plan
terragrunt apply


Debug
export TF_PLUGIN_TIMEOUT=5m
rm -rf ~/.terraform.d/plugins
rm -rf .terraform
rm -rf .terragrunt-cache
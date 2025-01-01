docker run -it -v $(pwd):/traffic-platform-copy devopsinfra/docker-terragrunt:aws-latest /bin/bash

cd aws
pip install -r resources/requirements.txt
python resources/main.py


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
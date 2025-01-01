import os
import yaml
from pathlib import Path
from trafficplatform_aws_v1_vpc_processor import VPCProcessor
from trafficplatform_aws_v1_account_processor import AccountProcessor
from trafficplatform_aws_v1_vpc import VPCResource
from trafficplatform_aws_v1_account import AccountResource

# Dictionary to keep track of resources by kind
resources_by_kind = {
    "Account": [],
    "VPC": [],
}

RESOURCE_MODELS = {
    "trafficplatform.aws/v1": {
        "Account": AccountResource,
        "VPC": VPCResource,
    },
}

PROCESSORS = {
    "trafficplatform.aws/v1": {
        "Account": AccountProcessor,
        "VPC": VPCProcessor,
    },
}

def process_yaml_file(file_path):
    """Process a single YAML file and group resources by kind."""
    # Load the YAML file
    with open(file_path, 'r') as file:
        resource_data = yaml.safe_load(file)

    kind = resource_data.get("kind")
    api_version = resource_data.get("apiVersion")

    # Extract required fields for the resource model
    metadata = resource_data["metadata"]
    spec = resource_data["spec"]

    # Get the correct resource model class based on api_version and kind
    if api_version not in RESOURCE_MODELS:
        print(f"Unsupported apiVersion: {api_version}")
        return

    if kind not in RESOURCE_MODELS[api_version]:
        print(f"Unsupported kind: {kind}")
        return

    resource_model_class = RESOURCE_MODELS[api_version][kind]

    # Initialize the resource model instance
    resource_instance = resource_model_class(api_version, kind, metadata, spec)

    # Add the resource instance to resources_by_kind
    resources_by_kind[kind].append(resource_instance)

def process_directory(input_dir):
    """Process all YAML/YML files in the given directory."""
    for root, _, files in os.walk(input_dir):
        for file in files:
            if file.endswith(".yaml") or file.endswith(".yml"):
                file_path = os.path.join(root, file)
                print(f"Processing file: {file_path}")
                process_yaml_file(file_path)

def main():
    current_dir = os.getcwd()
    input_dir = os.path.join(current_dir, "config")

    # Process all YAML/YML files in the directory and group resources by kind
    process_directory(input_dir)

    # Process resources for each kind using the corresponding processors
    if resources_by_kind["Account"]:
        account_processor = AccountProcessor(resources_by_kind["Account"])
        account_processor.process()

    if resources_by_kind["VPC"]:
        vpc_processor = VPCProcessor(resources_by_kind["VPC"])
        vpc_processor.process()

if __name__ == "__main__":
    main()
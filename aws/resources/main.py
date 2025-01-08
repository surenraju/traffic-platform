import os
import yaml
import argparse
from pathlib import Path
from trafficplatform_aws_v1_vpc_processor import VPCProcessor
from trafficplatform_aws_v1_core_network_attachment_processor import CoreNetworkAttachmentProcessor
from trafficplatform_aws_v1_vpc import VPCResource
from trafficplatform_aws_v1_core_network_attachment import CoreNetworkAttachmentResource

# Dictionary to keep track of resources by kind
resources_by_kind = {
    "CoreNetworkAttachment": [],
    "VPC": [],
}

RESOURCE_MODELS = {
    "trafficplatform.aws/v1": {
        "CoreNetworkAttachment": CoreNetworkAttachmentResource,
        "VPC": VPCResource,
    },
}

PROCESSORS = {
    "trafficplatform.aws/v1": {
        "CoreNetworkAttachment": CoreNetworkAttachmentProcessor,
        "VPC": VPCProcessor,
    },
}

def process_yaml_file(file_path, resource_type=None):
    """Process a single YAML file and group resources by kind."""
    # Load the YAML file
    with open(file_path, 'r') as file:
        resource_data = yaml.safe_load(file)

    kind = resource_data.get("kind")
    api_version = resource_data.get("apiVersion")

    # Skip processing if the kind does not match the specified resource type
    if resource_type and kind != resource_type:
        print(f"Skipping file: {file_path}")
        return

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


def process_directory(input_dir, resource_type=None):
    """Process all YAML/YML files in the given directory for the specified resource type."""
    for root, _, files in os.walk(input_dir):
        for file in files:
            if file.endswith(".yaml") or file.endswith(".yml"):
                file_path = os.path.join(root, file)
                process_yaml_file(file_path, resource_type)

def main():
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description="Process YAML files and group resources by kind.")
    parser.add_argument("--resource", type=str, help="Specify the resource type to process (e.g., CoreNetworkAttachment, VPC).")
    args = parser.parse_args()

    current_dir = os.getcwd()
    input_dir = os.path.join(current_dir, "config")

    print(f"Processing resource type '{args.resource}'")
    # Process all YAML/YML files in the directory for the specified resource type
    process_directory(input_dir, args.resource)

    # Process resources for each kind using the corresponding processors
    if resources_by_kind["CoreNetworkAttachment"]:
        core_network_attachment_processor = CoreNetworkAttachmentProcessor(resources_by_kind["CoreNetworkAttachment"])
        core_network_attachment_processor.process()

    if resources_by_kind["VPC"]:
        vpc_processor = VPCProcessor(resources_by_kind["VPC"])
        vpc_processor.process()

if __name__ == "__main__":
    main()
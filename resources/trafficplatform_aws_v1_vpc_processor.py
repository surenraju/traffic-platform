from base_processor import BaseProcessor
import json

class VPCProcessor(BaseProcessor):
    def validate(self):
        """Validate the list of VPC resources for duplicates and conflicting CIDR blocks."""
        seen_names = set()
        cidr_blocks = set()

        for resource in self.resources:
            # Print the resource as a formatted JSON
            print(f"Processing resource: {json.dumps(vars(resource), indent=2)}")

            # Accessing attributes via dot notation, not subscripting
            resource_name = resource.metadata["name"]  # Accessing name via attribute
            cidr_block = resource.spec.get("cidr_block")  # Accessing spec attribute

            # Check for duplicate resource names
            if resource_name in seen_names:
                print(f"Duplicate name found: {resource_name}")
                return False
            seen_names.add(resource_name)

            # Check for conflicting CIDR blocks
            if cidr_block in cidr_blocks:
                print(f"Conflicting CIDR block found: {cidr_block}")
                return False
            cidr_blocks.add(cidr_block)

        return True

    def validate_resource(self, resource):
        """Validate a single VPC resource by calling resource's validate method."""
        return resource.validate()  # Calls the validate method of the resource model

    def transform_resource(self, resource):
        """Transform a single VPC resource into the Terragrunt configuration by calling the resource's transform method."""
        return resource.transform()  # Calls the transform method of the resource model
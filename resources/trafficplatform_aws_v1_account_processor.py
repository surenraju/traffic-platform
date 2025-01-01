# account_processor.py
from base_processor import BaseProcessor

class AccountProcessor(BaseProcessor):
    def validate(self):
        """Validate the list of Account resources for duplicates."""
        seen_names = set()

        for resource in self.resources:
            resource_name = resource["metadata"]["name"]

            # Check for duplicate resource names
            if resource_name in seen_names:
                print(f"Duplicate name found: {resource_name}")
                return False
            seen_names.add(resource_name)

        return True

    def validate_resource(self, resource):
        """Validate a single Account resource by calling resource's validate method."""
        return resource.validate()  # Calls the validate method of the resource model

    def transform_resource(self, resource):
        """Transform a single Account resource into the Terragrunt configuration by calling the resource's transform method."""
        return resource.transform()  # Calls the transform method of the resource model
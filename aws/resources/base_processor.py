from pathlib import Path

class BaseProcessor:
    def __init__(self, resources):
        self.resources = resources

    def process(self):
        """Process the resources: validate and transform them."""
        if not self.validate():
            return None
        
        # Loop through resources, validate, transform, and write them to the filesystem
        for resource in self.resources:
            if not self.validate_resource(resource):
                print(f"Validation failed for resource {resource.metadata['name']}")
                continue
            
            terragrunt_content = self.transform_resource(resource)
            if terragrunt_content:
                self.write_to_filesystem(resource, terragrunt_content)
            else:
                print(f"Failed to render terragrunt.hcl for {resource.metadata['name']}")

    def validate(self):
        """Validate the list of resources. This will be overridden by each processor."""
        raise NotImplementedError("Validation method must be implemented in a subclass.")

    def validate_resource(self, resource):
        """Validate a single resource. This method can be overridden in a subclass if needed."""
        return resource.validate()  # Calls the validate method of the resource model

    def transform_resource(self, resource):
        """Transform a single resource into the desired output format. This will be overridden by each processor."""
        return resource.transform()  # Calls the transform method of the resource model

    def write_to_filesystem(self, resource, terragrunt_content):
        """Write the transformed content to the filesystem."""
        resource_name = resource.metadata["name"]  # Access name attribute
        account_id = resource.spec.get("account_id")  # Access account_id attribute
        kind = resource.kind  # Access kind attribute
        live_dir = Path(f"live/{account_id}/{kind}/{resource_name}")
        live_dir.mkdir(parents=True, exist_ok=True)

        terragrunt_file_path = live_dir / "terragrunt.hcl"
        with open(terragrunt_file_path, "w") as terragrunt_file:
            terragrunt_file.write(terragrunt_content)
        
        print(f"Generated terragrunt.hcl for {kind} {resource_name} at {terragrunt_file_path}")
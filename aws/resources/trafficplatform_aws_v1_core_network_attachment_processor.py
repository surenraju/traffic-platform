from base_processor import BaseProcessor
from vpc_utils import find_vpcs_with_tgw_attachment, find_tgw_id
import json
from jinja2 import Template
from pathlib import Path


class CoreNetworkAttachmentProcessor(BaseProcessor):
    def process(self):
        """
        Process the resources: validate and transform them. Generate one terragrunt.hcl
        for each VPC found in the specified account and region. Additionally, generate
        reachability analysis terragrunt files for each VPC pair.
        """
        if not self.validate():
            return None

        processed_vpcs = []  # Global list of all processed VPCs

        for resource in self.resources:
            if not self.validate_resource(resource):
                print(f"Validation failed for resource {resource.metadata['name']}")
                continue

            try:
                # Find VPCs with TGW attachment in the specified account and region
                account_id = resource.spec.get("account_id")
                region = resource.spec.get("region")
                vpcs = find_vpcs_with_tgw_attachment(account_id, region)

                if not vpcs:
                    raise ValueError(
                        f"No VPCs with Transit Gateway attachments found for account "
                        f"{account_id} in region {region}."
                    )

                # Get Transit Gateway ID
                transit_gateway_id = find_tgw_id(region)
                if not transit_gateway_id:
                    raise ValueError(
                        f"No Transit Gateway found for region {region}. "
                        "Please check your TGW configuration."
                    )

                # Render and write one terragrunt.hcl file per VPC
                for vpc in vpcs:
                    terragrunt_content = self.transform_resource(resource, vpc, transit_gateway_id)
                    if terragrunt_content:
                        self.write_to_filesystem(resource, vpc, terragrunt_content)
                        processed_vpcs.append(vpc)  # Add VPC to the processed list
                    else:
                        print(f"Failed to render terragrunt.hcl for VPC {vpc['vpc_id']}")

            except ValueError as e:
                print(f"Error processing resource {resource.metadata['name']}: {e}")

        # Call the reachability analysis function
        self.process_network_path_analysis(processed_vpcs, account_id)

    def validate(self):
        """Validate the list of CoreNetworkAttachment resources."""
        seen_names = set()

        for resource in self.resources:
            resource_name = resource.metadata["name"]

            # Check for duplicate resource names
            if resource_name in seen_names:
                print(f"Duplicate name found: {resource_name}")
                return False
            seen_names.add(resource_name)

        return True

    def validate_resource(self, resource):
        """Validate a single CoreNetworkAttachment resource."""
        return resource.validate()

    def transform_resource(self, resource, vpc, transit_gateway_id):
        """
        Transform a single CoreNetworkAttachment resource into terragrunt.hcl content.

        Args:
            resource: The resource being processed.
            vpc: The VPC data dictionary.
            transit_gateway_id: The Transit Gateway ID for the region.

        Returns:
            str: The rendered terragrunt.hcl content.
        """
       
        spec = {
            "vpc_id": vpc["vpc_id"],
            "subnet_ids": vpc.get("subnet_ids", []),
            "route_table_ids": vpc.get("route_table_ids", []),
            "transit_gateway_id": transit_gateway_id,
        }

        # Validate required fields
        if not spec["subnet_ids"]:
            raise ValueError(f"No subnet IDs found for VPC {vpc['vpc_id']}.")
        if not spec["route_table_ids"]:
            raise ValueError(f"No route table IDs found for VPC {vpc['vpc_id']}.")

        # Render the Jinja2 template
        template_file = resource.TEMPLATE_FILE
        with open(template_file, 'r') as file:
            template_content = file.read()

        template = Template(template_content)
        return template.render(**spec)

    def write_to_filesystem(self, resource, vpc, terragrunt_content):
        """
        Write the rendered terragrunt.hcl content to the filesystem.

        Args:
            resource: The resource being processed.
            vpc: The VPC data dictionary.
            terragrunt_content: The rendered terragrunt.hcl content.
        """
        account_id = resource.spec.get("account_id")
        kind = resource.kind
        vpc_name = vpc.get("vpc_name") or vpc["vpc_id"]
        live_dir = Path(f"live/{account_id}/{kind}/{vpc_name}")
        live_dir.mkdir(parents=True, exist_ok=True)

        terragrunt_file_path = live_dir / "terragrunt.hcl"
        with open(terragrunt_file_path, "w") as terragrunt_file:
            terragrunt_file.write(terragrunt_content)

        print(f"Generated terragrunt.hcl for {kind} {resource.metadata['name']} at {terragrunt_file_path}")

    def process_network_path_analysis(self, vpcs, account_id):
        """
        Generate reachability analysis terragrunt files for all VPC pairs.

        Args:
            vpcs: List of processed VPCs.
            account_id: Account ID for the resources.
        """
        for source_vpc in vpcs:
            for target_vpc in vpcs:
                if source_vpc["vpc_id"] == target_vpc["vpc_id"]:
                    continue  # Skip pairing a VPC with itself

                # Generate args for the reachability analyzer
                args = {
                    "source_vpc_id": source_vpc["vpc_id"],
                    "target_vpc_id": target_vpc["vpc_id"]
                }

                # Use VPC names for the directory naming
                source_vpc_name = source_vpc.get("vpc_name", source_vpc["vpc_id"])
                target_vpc_name = target_vpc.get("vpc_name", target_vpc["vpc_id"])
                vpc_pair_name = f"{source_vpc_name}_to_{target_vpc_name}"

                # Render the reachability analyzer terragrunt.hcl content
                template_file = 'resources/templates/reachability_analyzer_terragrunt.hcl.j2'
                with open(template_file, 'r') as file:
                    template_content = file.read()

                template = Template(template_content)
                content = template.render(**args)

                # Write the rendered content to the filesystem
                live_dir = Path(f"live/{account_id}/ReachabilityAnalyzer/{vpc_pair_name}")
                live_dir.mkdir(parents=True, exist_ok=True)

                terragrunt_file_path = live_dir / "terragrunt.hcl"
                with open(terragrunt_file_path, "w") as terragrunt_file:
                    terragrunt_file.write(content)

                print(f"Generated reachability analysis terragrunt.hcl for {vpc_pair_name} at {terragrunt_file_path}")
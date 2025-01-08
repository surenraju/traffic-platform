from base_resource import BaseResource
from vpc_utils import find_vpcs_with_tgw_attachment, find_tgw_id


class CoreNetworkAttachmentResource(BaseResource):
    SCHEMA_FILE = "resources/schemas/core_network_attachment_schema.json"
    TEMPLATE_FILE = "resources/templates/core_network_attachment_terragrunt.hcl.j2"

    def validate(self):
        """
        Validate the resource against the schema and ensure the region has a valid Transit Gateway,
        and that VPCs, subnets, and route tables are properly configured.
        """
        # Validate against the JSON schema
        if not super().validate(self.SCHEMA_FILE):
            return False

        # Validate that the region has a valid Transit Gateway
        transit_gateway_id = find_tgw_id(self.spec["region"])
        if not transit_gateway_id:
            print(
                f"Validation error: No Transit Gateway found for region {self.spec['region']}. "
                "Please check your TGW configuration."
            )
            return False

        # Find VPCs with TGW attachment
        vpcs = find_vpcs_with_tgw_attachment(self.spec["account_id"], self.spec["region"])
        if not vpcs:
            print(
                f"Validation error: No VPCs with Transit Gateway attachments found for account "
                f"{self.spec['account_id']} in region {self.spec['region']}."
            )
            return False

        # Validate that subnet IDs and route table IDs are available
        for vpc in vpcs:
            if not vpc.get("subnet_ids"):
                print(f"Validation error: No subnet IDs found for VPC {vpc['vpc_id']}.")
                return False

            if not vpc.get("route_table_ids"):
                print(f"Validation error: No route table IDs found for VPC {vpc['vpc_id']}.")
                return False

        return True

    def transform(self):
        """
        Transform the resource into the desired output format using the Jinja2 template.
        """
        raise NotImplementedError("Transform is not supported for CoreNetworkAttachmentResource.")
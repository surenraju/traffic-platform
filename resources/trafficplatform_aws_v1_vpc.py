from base_resource import BaseResource


class VPCResource(BaseResource):
    SCHEMA_FILE = "resources/schemas/vpc_schema.json"
    TEMPLATE_FILE = "resources/templates/vpc_terragrunt.hcl.j2"

    def validate(self):
        return super().validate(self.SCHEMA_FILE)

    def transform(self):
        return super().transform(self.TEMPLATE_FILE)
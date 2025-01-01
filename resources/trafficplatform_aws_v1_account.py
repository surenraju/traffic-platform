from base_resource import BaseResource


class AccountResource(BaseResource):
    SCHEMA_FILE = "resources/schemas/account_schema.json"
    TEMPLATE_FILE = "resources/templates/account_terragrunt.hcl.j2"

    def validate(self):
        return super().validate(self.SCHEMA_FILE)

    def transform(self):
        raise NotImplementedError("Transformation is not supported for AccountResource.")
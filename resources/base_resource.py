import json
import yaml
from jsonschema import validate, ValidationError
from jinja2 import Template


class BaseResource:
    def __init__(self, api_version, kind, metadata, spec):
        self.apiVersion = api_version
        self.kind = kind
        self.metadata = metadata
        self.spec = spec

    def validate(self, schema_file):
        """
        Validate the resource against a JSON schema.
        """
        with open(schema_file, 'r') as schema:
            schema_data = json.load(schema)

        resource_data = {
            "apiVersion": self.apiVersion,
            "kind": self.kind,
            "metadata": self.metadata,
            "spec": self.spec
        }

        try:
            validate(instance=resource_data, schema=schema_data)
            return True
        except ValidationError as e:
            print(f"Validation error: {e.message}")
            return False

    @classmethod
    def from_json(cls, json_str):
        """
        Load the resource from a JSON string.
        """
        data = json.loads(json_str)
        return cls(data["apiVersion"], data["kind"], data["metadata"], data["spec"])

    @classmethod
    def from_yaml(cls, yaml_str):
        """
        Load the resource from a YAML string.
        """
        data = yaml.safe_load(yaml_str)
        return cls(data["apiVersion"], data["kind"], data["metadata"], data["spec"])

    def transform(self, template_file):
        """
        Transform the resource using a Jinja2 template.
        """
        with open(template_file, 'r') as template_file:
            template_content = template_file.read()

        template = Template(template_content)
        return template.render(**self.spec)
    

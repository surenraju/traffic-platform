import boto3
import os
import yaml

def find_vpcs_with_tgw_attachment(account_id, region):
    """Find all VPCs, subnets, and route tables with TransitGatewayAttachment=true."""
    ec2 = boto3.client('ec2', region_name=region)

    vpcs = ec2.describe_vpcs(
        Filters=[{"Name": "tag:TransitGatewayAttachment", "Values": ["true"]}]
    )["Vpcs"]

    results = []
    for vpc in vpcs:
        vpc_id = vpc["VpcId"]

        # Fetch VPC name from the tags
        vpc_name = None
        if "Tags" in vpc:
            for tag in vpc["Tags"]:
                if tag["Key"] == "Name":
                    vpc_name = tag["Value"]
                    break
        vpc_name = vpc_name or vpc_id  # Use vpc_id as fallback if no Name tag is found

        # Fetch subnets
        subnets = ec2.describe_subnets(
            Filters=[
                {"Name": "vpc-id", "Values": [vpc_id]},
                {"Name": "tag:TransitGatewayAttachment", "Values": ["true"]}
            ]
        )["Subnets"]
        subnet_ids = [subnet["SubnetId"] for subnet in subnets]

        # Fetch route tables
        route_tables = ec2.describe_route_tables(
            Filters=[
                {"Name": "vpc-id", "Values": [vpc_id]},
                {"Name": "tag:TransitGatewayAttachment", "Values": ["true"]}
            ]
        )["RouteTables"]
        route_table_ids = [rt["RouteTableId"] for rt in route_tables]

        results.append({
            "vpc_id": vpc_id,
            "vpc_name": vpc_name, 
            "subnet_ids": subnet_ids,
            "route_table_ids": route_table_ids,
        })

    return results


def find_tgw_id(region="eu-west-1", tgw_file="tgw.yaml"):
    """
    Find the Transit Gateway (TGW) ID for the given region.

    Args:
        region (str): The region to look for the TGW ID.
        tgw_file (str): The name of the YAML file containing TGW specifications. Default is 'tgw.yaml'.

    Returns:
        str: The TGW ID for the given region if found, otherwise None.
    """
    try:
        # Get the path to the tgw.yaml file in the "config" directory
        current_dir = os.getcwd()
        input_dir = os.path.join(current_dir, "config")
        tgw_file_path = os.path.join(input_dir, tgw_file)

        # Open and parse the tgw.yaml file
        with open(tgw_file_path, 'r') as file:
            tgw_data = yaml.safe_load(file)

        # Check if the spec and transit_gateways keys exist
        transit_gateways = tgw_data.get("spec", {}).get("transit_gateways", [])
        if not transit_gateways:
            print(f"No transit gateways found in {tgw_file_path}.")
            return None

        # Search for the region in the list of transit gateways
        for tgw in transit_gateways:
            if tgw.get("region") == region:
                return tgw.get("id")

        print(f"No Transit Gateway found for region: {region}")
        return None

    except FileNotFoundError:
        print(f"File {tgw_file_path} not found.")
        return None
    except yaml.YAMLError as e:
        print(f"Error parsing YAML file {tgw_file_path}: {e}")
        return None
#!/usr/bin/env python3
"""
PostgreSQL AGE Operator API - Customer Management Example
This script demonstrates how to manage customer instances via the API
"""

import argparse
import json
import os
import sys
import time
from typing import Dict, List, Optional

import requests


class PostgresAGEAPIClient:
    """Client for PostgreSQL AGE Operator API"""

    def __init__(self, api_url: str, api_key: str):
        self.api_url = api_url.rstrip('/')
        self.headers = {
            'Authorization': f'Bearer {api_key}',
            'Content-Type': 'application/json'
        }
        self.session = requests.Session()
        self.session.headers.update(self.headers)

    def health_check(self) -> Dict:
        """Check API server health"""
        response = self.session.get(f'{self.api_url}/health')
        response.raise_for_status()
        return response.json()

    def create_customer(self, name: str, display_name: str = None,
                       image_tag: str = 'latest', config: Dict = None,
                       labels: Dict = None) -> Dict:
        """Create a new customer instance"""
        payload = {
            'name': name,
            'displayName': display_name or name,
            'imageTag': image_tag,
            'config': config or {},
            'labels': labels or {}
        }

        response = self.session.post(
            f'{self.api_url}/api/v1/customers',
            json=payload
        )

        if response.status_code == 400:
            error_data = response.json()
            print(f"❌ Validation failed: {error_data['error']['message']}")
            if 'details' in error_data['error']:
                for detail in error_data['error']['details']:
                    print(f"   {detail['field']}: {detail['message']}")
            return None

        response.raise_for_status()
        return response.json()

    def get_customer(self, name: str) -> Dict:
        """Get customer instance details"""
        response = self.session.get(f'{self.api_url}/api/v1/customers/{name}')
        if response.status_code == 404:
            return None
        response.raise_for_status()
        return response.json()

    def list_customers(self, page: int = 1, page_size: int = 10) -> Dict:
        """List all customer instances"""
        params = {'page': page, 'pageSize': page_size}
        response = self.session.get(
            f'{self.api_url}/api/v1/customers',
            params=params
        )
        response.raise_for_status()
        return response.json()

    def update_customer(self, name: str, updates: Dict) -> Dict:
        """Update customer instance"""
        response = self.session.put(
            f'{self.api_url}/api/v1/customers/{name}',
            json=updates
        )
        if response.status_code == 404:
            return None
        response.raise_for_status()
        return response.json()

    def delete_customer(self, name: str) -> Dict:
        """Delete customer instance"""
        response = self.session.delete(f'{self.api_url}/api/v1/customers/{name}')
        if response.status_code == 404:
            return None
        response.raise_for_status()
        return response.json()

    def get_customer_status(self, name: str) -> Dict:
        """Get customer instance status"""
        response = self.session.get(f'{self.api_url}/api/v1/customers/{name}/status')
        if response.status_code == 404:
            return None
        response.raise_for_status()
        return response.json()

    def wait_for_customer_ready(self, name: str, timeout: int = 300) -> bool:
        """Wait for customer instance to be ready"""
        start_time = time.time()
        while time.time() - start_time < timeout:
            status_response = self.get_customer_status(name)
            if not status_response:
                print(f"❌ Customer {name} not found")
                return False

            status = status_response['data']
            print(f"⏳ Customer {name} status: {status['phase']}")

            if status['ready']:
                print(f"✅ Customer {name} is ready!")
                return True

            if status['phase'] == 'Failed':
                print(f"❌ Customer {name} failed: {status.get('message', 'Unknown error')}")
                return False

            time.sleep(10)

        print(f"⏰ Timeout waiting for customer {name} to be ready")
        return False


def create_environment_config(environment: str) -> Dict:
    """Create configuration based on environment"""
    configs = {
        'development': {
            'resources': {
                'requests': {'cpu': '100m', 'memory': '256Mi'},
                'limits': {'cpu': '500m', 'memory': '1Gi'}
            },
            'highAvailability': False,
            'backupEnabled': False,
            'monitoringEnabled': False
        },
        'staging': {
            'resources': {
                'requests': {'cpu': '250m', 'memory': '512Mi'},
                'limits': {'cpu': '1000m', 'memory': '2Gi'}
            },
            'highAvailability': True,
            'backupEnabled': True,
            'monitoringEnabled': True
        },
        'production': {
            'resources': {
                'requests': {'cpu': '500m', 'memory': '1Gi'},
                'limits': {'cpu': '2000m', 'memory': '4Gi'}
            },
            'highAvailability': True,
            'backupEnabled': True,
            'monitoringEnabled': True
        }
    }
    return configs.get(environment, configs['development'])


def cmd_create(client: PostgresAGEAPIClient, args) -> None:
    """Create customer command"""
    print(f"🚀 Creating customer: {args.name}")

    config = create_environment_config(args.environment)
    labels = {
        'environment': args.environment,
        'managed-by': 'python-api-client'
    }

    if args.team:
        labels['team'] = args.team

    result = client.create_customer(
        name=args.name,
        display_name=args.display_name,
        image_tag=args.image_tag,
        config=config,
        labels=labels
    )

    if result:
        print("✅ Customer created successfully!")
        print(json.dumps(result['data'], indent=2))

        if args.wait:
            print("\n⏳ Waiting for customer to be ready...")
            client.wait_for_customer_ready(args.name)


def cmd_list(client: PostgresAGEAPIClient, args) -> None:
    """List customers command"""
    result = client.list_customers(page=args.page, page_size=args.page_size)

    if result['data']['customers']:
        print(f"📋 Customers (Page {result['data']['page']}/{(result['data']['total'] + args.page_size - 1) // args.page_size}):")
        print()

        for customer in result['data']['customers']:
            status_icon = "✅" if customer['status']['ready'] else "⏳"
            print(f"{status_icon} {customer['name']}")
            print(f"   Display Name: {customer.get('displayName', 'N/A')}")
            print(f"   Status: {customer['status']['phase']}")
            print(f"   Namespace: {customer['namespace']}")
            print(f"   Image Tag: {customer.get('imageTag', 'N/A')}")
            print(f"   Created: {customer['createdAt']}")
            print()

        print(f"Total: {result['data']['total']} customers")
    else:
        print("📭 No customers found")


def cmd_get(client: PostgresAGEAPIClient, args) -> None:
    """Get customer command"""
    result = client.get_customer(args.name)

    if result:
        print(f"📋 Customer: {args.name}")
        print(json.dumps(result['data'], indent=2))
    else:
        print(f"❌ Customer '{args.name}' not found")


def cmd_update(client: PostgresAGEAPIClient, args) -> None:
    """Update customer command"""
    updates = {}

    if args.image_tag:
        updates['imageTag'] = args.image_tag

    if args.display_name:
        updates['displayName'] = args.display_name

    if args.environment:
        updates['config'] = create_environment_config(args.environment)

    if not updates:
        print("❌ No updates specified")
        return

    print(f"🔄 Updating customer: {args.name}")
    result = client.update_customer(args.name, updates)

    if result:
        print("✅ Customer updated successfully!")
        print(json.dumps(result['data'], indent=2))
    else:
        print(f"❌ Customer '{args.name}' not found")


def cmd_delete(client: PostgresAGEAPIClient, args) -> None:
    """Delete customer command"""
    if not args.confirm:
        confirm = input(f"⚠️  Are you sure you want to delete customer '{args.name}'? (y/N): ")
        if confirm.lower() != 'y':
            print("❌ Deletion cancelled")
            return

    print(f"🗑️ Deleting customer: {args.name}")
    result = client.delete_customer(args.name)

    if result:
        print("✅ Customer deleted successfully!")
    else:
        print(f"❌ Customer '{args.name}' not found")


def cmd_status(client: PostgresAGEAPIClient, args) -> None:
    """Get customer status command"""
    result = client.get_customer_status(args.name)

    if result:
        status = result['data']
        status_icon = "✅" if status['ready'] else "⏳"
        print(f"{status_icon} Customer: {args.name}")
        print(f"   Phase: {status['phase']}")
        print(f"   Ready: {status['ready']}")
        print(f"   Replicas: {status.get('replicas', 0)}/{status.get('readyReplicas', 0)}")
        if status.get('message'):
            print(f"   Message: {status['message']}")
        print(f"   Last Updated: {status['lastUpdated']}")
    else:
        print(f"❌ Customer '{args.name}' not found")


def main():
    """Main function"""
    parser = argparse.ArgumentParser(
        description='PostgreSQL AGE Operator API Client'
    )
    parser.add_argument(
        '--api-url',
        default=os.getenv('API_URL', 'http://localhost:8080'),
        help='API server URL'
    )
    parser.add_argument(
        '--api-key',
        default=os.getenv('API_KEY'),
        help='API key for authentication'
    )

    subparsers = parser.add_subparsers(dest='command', help='Available commands')

    # Create command
    create_parser = subparsers.add_parser('create', help='Create customer instance')
    create_parser.add_argument('name', help='Customer name')
    create_parser.add_argument('--display-name', help='Display name')
    create_parser.add_argument('--image-tag', default='latest', help='Image tag')
    create_parser.add_argument('--environment', default='development',
                              choices=['development', 'staging', 'production'],
                              help='Environment type')
    create_parser.add_argument('--team', help='Team label')
    create_parser.add_argument('--wait', action='store_true',
                              help='Wait for customer to be ready')

    # List command
    list_parser = subparsers.add_parser('list', help='List customer instances')
    list_parser.add_argument('--page', type=int, default=1, help='Page number')
    list_parser.add_argument('--page-size', type=int, default=10, help='Page size')

    # Get command
    get_parser = subparsers.add_parser('get', help='Get customer instance')
    get_parser.add_argument('name', help='Customer name')

    # Update command
    update_parser = subparsers.add_parser('update', help='Update customer instance')
    update_parser.add_argument('name', help='Customer name')
    update_parser.add_argument('--display-name', help='New display name')
    update_parser.add_argument('--image-tag', help='New image tag')
    update_parser.add_argument('--environment',
                              choices=['development', 'staging', 'production'],
                              help='New environment configuration')

    # Delete command
    delete_parser = subparsers.add_parser('delete', help='Delete customer instance')
    delete_parser.add_argument('name', help='Customer name')
    delete_parser.add_argument('--confirm', action='store_true',
                              help='Skip confirmation prompt')

    # Status command
    status_parser = subparsers.add_parser('status', help='Get customer status')
    status_parser.add_argument('name', help='Customer name')

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    if not args.api_key:
        print("❌ API key is required. Set API_KEY environment variable or use --api-key")
        sys.exit(1)

    try:
        client = PostgresAGEAPIClient(args.api_url, args.api_key)

        # Test connection
        client.health_check()

        # Execute command
        commands = {
            'create': cmd_create,
            'list': cmd_list,
            'get': cmd_get,
            'update': cmd_update,
            'delete': cmd_delete,
            'status': cmd_status
        }

        commands[args.command](client, args)

    except requests.exceptions.ConnectionError:
        print(f"❌ Failed to connect to API server at {args.api_url}")
        print("   Make sure the API server is running and accessible")
        sys.exit(1)
    except requests.exceptions.HTTPError as e:
        if e.response.status_code == 401:
            print("❌ Authentication failed. Check your API key")
        else:
            print(f"❌ API request failed: {e}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Unexpected error: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()
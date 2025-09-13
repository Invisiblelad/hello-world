import os
import requests

RANCHER_URL = os.environ.get('RANCHER_URL')
RANCHER_TOKEN = os.environ.get('RANCHER_TOKEN')

if not RANCHER_URL or not RANCHER_TOKEN:
    print("Please set RANCHER_URL and RANCHER_TOKEN environment variables.")
    exit(1)

headers = {
    'Authorization': f'Bearer {RANCHER_TOKEN}',
    'Accept': 'application/json'
}

def get_clusters():
    url = f"{RANCHER_URL}/v3/clusters"
    resp = requests.get(url, headers=headers, verify=False)
    resp.raise_for_status()
    return resp.json().get('data', [])

def get_kubeconfig(cluster_id):
    url = f"{RANCHER_URL}/v3/clusters/{cluster_id}?action=generateKubeconfig"
    resp = requests.post(url, headers=headers, verify=False)
    resp.raise_for_status()
    return resp.json().get('config', None)

def main():
    clusters = get_clusters()
    for cluster in clusters:
        cluster_id = cluster['id']
        cluster_name = cluster['name']
        print(f"Fetching kubeconfig for {cluster_name} ({cluster_id})...")
        kubeconfig = get_kubeconfig(cluster_id)
        if kubeconfig:
            filename = f"{cluster_name}.yaml"
            with open(filename, 'w') as f:
                f.write(kubeconfig)
            print(f"Saved kubeconfig to {filename}")
        else:
            print(f"Failed to get kubeconfig for {cluster_name}")

if __name__ == "__main__":
    main()

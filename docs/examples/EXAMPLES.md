# Harbor CLI Examples

## Authentication and Configuration

### First-time Setup

```bash
# Interactive setup
$ hrbcli config init
Harbor URL: https://harbor.example.com
Username: admin
Password: ********
Verify configuration? (y/n): y
✓ Configuration saved to ~/.hrbcli.yaml
✓ Successfully connected to Harbor

# Manual configuration
$ hrbcli config set harbor_url https://harbor.example.com
$ hrbcli config set username admin
$ export HARBOR_PASSWORD=secretpassword
```

### Multiple Harbor Instances

```bash
# Production Harbor
$ hrbcli --config ~/.hrbcli-prod.yaml project list

# Development Harbor
$ hrbcli --config ~/.hrbcli-dev.yaml project list

# Using environment variables
$ HARBOR_URL=https://dev.harbor.example.com hrbcli project list
```

## Project Management

### Create Projects with Different Settings

```bash
# Simple project creation
$ hrbcli project create my-app

# Public project with quotas
$ hrbcli project create shared-images \
    --public \
    --storage-limit 50G \
    --member-limit 100

# Project with security settings
$ hrbcli project create secure-app \
    --enable-content-trust \
    --prevent-vulnerable \
    --severity-threshold high \
    --auto-scan
```

### Manage Project Members

```bash
# Add user as developer
$ hrbcli project member add my-app john --role developer

# Add user as admin
$ hrbcli project member add my-app alice --role admin

# List project members
$ hrbcli project member list my-app

# Remove member
$ hrbcli project member remove my-app john
```

## Repository Operations

### Working with Repositories

```bash
# List all repositories with sizes
$ hrbcli repo list my-app --detail

NAME                SIZE      TAGS  PULLS  LAST MODIFIED
my-app/frontend     1.2 GB    15    1234   2 hours ago
my-app/backend      856 MB    23    5678   1 hour ago
my-app/database     2.1 GB    8     910    3 days ago

# Search repositories
$ hrbcli repo list my-app --filter "*frontend*"

# Get repository info
$ hrbcli repo info my-app/frontend
```

### Tag Management

```bash
# List tags with vulnerability status
$ hrbcli repo tags my-app/frontend --detail

TAG      SIZE    VULNERABILITIES  SCAN STATUS  CREATED
latest   245 MB  High: 2, Med: 5  Finished     2 hours ago
v2.1.0   245 MB  High: 0, Med: 3  Finished     1 day ago
v2.0.0   238 MB  High: 5, Med: 8  Finished     1 week ago

# Delete old tags
$ hrbcli repo delete my-app/frontend:v1.0.0
$ hrbcli repo delete my-app/frontend:v1.1.0

# Bulk delete tags
$ hrbcli repo tags my-app/frontend --filter "v1.*" -o json | \
    jq -r '.[] | .name' | \
    xargs -I {} hrbcli repo delete my-app/frontend:{}
```

## Security Scanning

### Scan Artifacts

```bash
# Scan single artifact
$ hrbcli artifact scan my-app/frontend:latest

# Scan all artifacts in repository
$ hrbcli artifact scan my-app/frontend --all

# Wait for scan to complete
$ hrbcli artifact scan my-app/frontend:latest --wait
```

### View Vulnerabilities

```bash
# List vulnerabilities
$ hrbcli artifact vulnerabilities my-app/frontend:latest

SEVERITY  CVE             PACKAGE         VERSION  FIXED VERSION
Critical  CVE-2021-12345  openssl         1.0.1    1.0.2
High      CVE-2021-12346  libcurl         7.1.0    7.2.0
High      CVE-2021-12347  nginx           1.18.0   1.19.0

# Export vulnerability report
$ hrbcli artifact vulnerabilities my-app/frontend:latest -o json > vulns.json

# Check if image is safe to deploy
$ hrbcli artifact vulnerabilities my-app/frontend:latest --severity high
$ echo $?  # Exit code 0 if no high/critical vulnerabilities
```

## Replication

### Set Up Replication

```bash
# Create push replication to remote Harbor
$ hrbcli replication create prod-sync \
    --source my-app \
    --destination https://prod-harbor.example.com \
    --destination-namespace production \
    --trigger "manual,schedule:0 2 * * *" \
    --filter "name:*/release-*"

# Create pull replication from Docker Hub
$ hrbcli replication create dockerhub-mirror \
    --source https://hub.docker.com \
    --source-filter "nginx,alpine,ubuntu" \
    --destination mirror \
    --direction pull \
    --trigger "schedule:0 */6 * * *"
```

### Monitor Replication

```bash
# List executions
$ hrbcli replication executions prod-sync

ID    STATUS      START TIME           END TIME
123   Succeeded   2024-01-10 02:00:00  2024-01-10 02:15:00
122   Failed      2024-01-09 02:00:00  2024-01-09 02:05:00

# Get execution details
$ hrbcli replication execution 122

# Get execution logs
$ hrbcli replication logs 122
```

## Proxy Cache

### Set Up a Proxy Cache Project

```bash
# Create a registry endpoint for Docker Hub
hrbcli registry create dockerhub --type docker-hub --url https://hub.docker.com

# Create the proxy cache project referencing that registry
hrbcli project create mycache \
    --proxy-cache \
    --registry-name dockerhub \
    --proxy-speed -1
```

The CLI resolves the registry ID automatically when `--registry-name` is
provided. There is currently no option to set a TTL or expiry time for cached
content.

## Automation Scripts

### Promote Images Through Environments

```bash
#!/bin/bash
# promote.sh - Promote image from dev to staging to prod

IMAGE=$1
VERSION=$2

# Scan in dev
echo "Scanning image in development..."
hrbcli artifact scan dev/${IMAGE}:${VERSION} --wait

# Check vulnerabilities
if hrbcli artifact vulnerabilities dev/${IMAGE}:${VERSION} --severity high; then
    echo "✓ No high/critical vulnerabilities found"
else
    echo "✗ High/critical vulnerabilities found, aborting"
    exit 1
fi

# Copy to staging
echo "Promoting to staging..."
hrbcli artifact copy dev/${IMAGE}:${VERSION} staging/${IMAGE}:${VERSION}

# Tag as staging-latest
hrbcli artifact tag staging/${IMAGE}:${VERSION} staging-latest

# After testing, promote to production
read -p "Promote to production? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    hrbcli artifact copy staging/${IMAGE}:${VERSION} prod/${IMAGE}:${VERSION}
    hrbcli artifact tag prod/${IMAGE}:${VERSION} latest
    echo "✓ Promoted to production"
fi
```

### Cleanup Old Images

```bash
#!/bin/bash
# cleanup.sh - Remove old images keeping last N tags

PROJECT=$1
REPO=$2
KEEP=5

# Get all tags sorted by creation date
TAGS=$(hrbcli repo tags ${PROJECT}/${REPO} -o json | \
    jq -r 'sort_by(.created) | reverse | .[].name')

# Keep the latest N tags, delete the rest
echo "$TAGS" | tail -n +$((KEEP+1)) | while read tag; do
    echo "Deleting ${PROJECT}/${REPO}:${tag}"
    hrbcli repo delete ${PROJECT}/${REPO}:${tag} --force
done
```

### Generate Security Report

```bash
#!/bin/bash
# security-report.sh - Generate security report for all images

OUTPUT="security-report-$(date +%Y%m%d).csv"
echo "Project,Repository,Tag,Critical,High,Medium,Low" > $OUTPUT

hrbcli project list -o json | jq -r '.[].name' | while read project; do
    hrbcli repo list $project -o json | jq -r '.[].name' | while read repo; do
        hrbcli repo tags $repo -o json | jq -r '.[].name' | while read tag; do
            VULNS=$(hrbcli artifact vulnerabilities ${repo}:${tag} -o json | \
                jq -r '[.[] | .severity] | group_by(.) | map({(.[0]): length}) | add')
            
            CRITICAL=$(echo $VULNS | jq -r '.Critical // 0')
            HIGH=$(echo $VULNS | jq -r '.High // 0')
            MEDIUM=$(echo $VULNS | jq -r '.Medium // 0')
            LOW=$(echo $VULNS | jq -r '.Low // 0')
            
            echo "$project,$repo,$tag,$CRITICAL,$HIGH,$MEDIUM,$LOW" >> $OUTPUT
        done
    done
done

echo "Report saved to $OUTPUT"
```

## Advanced Usage

### Using with CI/CD

```yaml
# .gitlab-ci.yml example
variables:
  HARBOR_URL: https://harbor.example.com
  
stages:
  - build
  - scan
  - deploy

build:
  stage: build
  script:
    - docker build -t ${HARBOR_URL}/my-app/backend:${CI_COMMIT_SHA} .
    - docker push ${HARBOR_URL}/my-app/backend:${CI_COMMIT_SHA}

scan:
  stage: scan
  script:
    - hrbcli artifact scan my-app/backend:${CI_COMMIT_SHA} --wait
    - hrbcli artifact vulnerabilities my-app/backend:${CI_COMMIT_SHA} --severity high
  allow_failure: false

deploy:
  stage: deploy
  script:
    - hrbcli artifact tag my-app/backend:${CI_COMMIT_SHA} latest
    - kubectl set image deployment/backend backend=${HARBOR_URL}/my-app/backend:latest
  only:
    - main
```

### JSON Processing with jq

```bash
# Get total storage used by project
$ hrbcli project get my-app -o json | jq '.current_usage.storage'

# List projects over quota
$ hrbcli project list -o json | \
    jq '.[] | select(.current_usage.storage > .quota.storage) | .name'

# Find images without recent scans
$ hrbcli repo list my-app -o json | \
    jq '.[] | select(.scan_overview.scan_status != "finished") | .name'

# Export artifact list with specific fields
$ hrbcli artifact list my-app/backend -o json | \
    jq '.[] | {digest: .digest, tags: .tags, size: .size, created: .created}'
```

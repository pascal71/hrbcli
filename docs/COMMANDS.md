# Harbor CLI Command Reference

## Global Flags

These flags are available for all commands:

```
--config string        Config file (default $HOME/.hrbcli.yaml)
--debug               Enable debug output
--harbor-url string   Harbor server URL
--insecure            Skip TLS certificate verification
--no-color            Disable colored output
-o, --output string   Output format (table|json|yaml) (default "table")
--password string     Harbor password
--username string     Harbor username
```

## Commands

### Project Management

#### `hrbcli project list`

List all projects accessible to the user.

```bash
# List all projects
hrbcli project list

# List projects with details
hrbcli project list --detail

# Output as JSON
hrbcli project list -o json

# Filter by name
hrbcli project list --name-filter "prod*"
```

#### `hrbcli project create`

Create a new project.

```bash
# Create a public project
hrbcli project create myproject --public

# Create with storage quota (in bytes, or use K, M, G, T)
hrbcli project create myproject --storage-limit 10G

# Create with member limit
hrbcli project create myproject --member-limit 50
```

#### `hrbcli project delete`

Delete a project.

```bash
# Delete a project
hrbcli project delete myproject

# Force delete without confirmation
hrbcli project delete myproject --force
```

#### `hrbcli project update`

Update project settings.

```bash
# Make project public
hrbcli project update myproject --public=true

# Update storage quota
hrbcli project update myproject --storage-limit 20G

# Enable content trust
hrbcli project update myproject --enable-content-trust
```

### Registry Management

#### `hrbcli registry list`

List registry endpoints configured in Harbor.

```bash
# List all registries
hrbcli registry list

# Search by name
hrbcli registry list --query docker
```

#### `hrbcli registry create`

Create a new registry endpoint. Use `--interactive` to be prompted for values.

```bash
# Create Docker Hub endpoint
hrbcli registry create dockerhub --type docker-hub --url https://hub.docker.com

# Interactive mode
hrbcli registry create --interactive
```

#### `hrbcli registry get`

Show details of a registry endpoint.

```bash
hrbcli registry get 1
```

#### `hrbcli registry update`

Update an existing registry configuration.

```bash
hrbcli registry update 1 --description "Updated"
```

#### `hrbcli registry delete`

Delete a registry endpoint.

```bash
hrbcli registry delete 1
```

#### `hrbcli registry ping`

Verify connectivity with a registry endpoint.

```bash
hrbcli registry ping --url https://registry.example.com --type docker-registry
```

#### `hrbcli registry adapters`

List available registry adapters.

```bash
hrbcli registry adapters
```

#### `hrbcli registry adapter-info`

Show detailed information about an adapter type.

```bash
hrbcli registry adapter-info docker-hub
```

### Repository Management

#### `hrbcli repo list`

List repositories in a project.

```bash
# List all repositories in a project
hrbcli repo list myproject

# List with details (size, tags count)
hrbcli repo list myproject --detail

# Filter by name
hrbcli repo list myproject --filter "app*"
```

#### `hrbcli repo get`

Show details for a repository.

```bash
# Get repository information
hrbcli repo get myproject/myapp

# Include creation/update timestamps
hrbcli repo get myproject/myapp --detail
```

#### `hrbcli repo delete`

Delete a repository.

```bash
# Delete entire repository
hrbcli repo delete myproject/myapp

# Delete specific tag
hrbcli repo delete myproject/myapp:v1.0.0

# Force delete without confirmation
hrbcli repo delete myproject/myapp --force
```

#### `hrbcli repo tags`

List tags for a repository.

```bash
# List all tags
hrbcli repo tags myproject/myapp

# List with details (size, scan status)
hrbcli repo tags myproject/myapp --detail

# Filter by tag name
hrbcli repo tags myproject/myapp --filter "v1.*"
```

### Artifact Management

#### `hrbcli artifact list`

List artifacts in a repository or across a project.

```bash
# List all artifacts in a repository
hrbcli artifact list myproject/myapp

# List all artifacts in a project
hrbcli artifact list myproject

# Include vulnerability summary
hrbcli artifact list myproject/myapp --with-scan-overview

# Include labels and extra details
hrbcli artifact list myproject/myapp --with-label --detail
```

#### `hrbcli artifact get`

Display details for a specific artifact.

```bash
# Show artifact information
hrbcli artifact get myproject/myapp:latest
```

#### `hrbcli artifact scan`

Trigger vulnerability scan.

```bash
# Scan specific artifact
hrbcli artifact scan myproject/myapp@sha256:abc123

# Scan all artifacts in repository
hrbcli artifact scan myproject/myapp --all

# Wait for scan to finish
hrbcli artifact scan myproject/myapp:latest --wait
```

#### `hrbcli artifact vulnerabilities`

Show vulnerability report for an artifact. Use `--summary` for an overview with counts by severity, or `--severity` to fail if vulnerabilities of that level or higher exist. The report can be saved to a file using `--file`.


```bash
# Show vulnerabilities
hrbcli artifact vulnerabilities myproject/myapp:latest

# Fail if high or critical vulns found
hrbcli artifact vulnerabilities myproject/myapp:latest --severity high

# Save report to file
hrbcli artifact vulnerabilities myproject/myapp:latest --file vulns.json -o json
```

#### `hrbcli artifact sbom`

Display the SBOM report for an artifact. Use `--file` to save the report locally.

```bash
hrbcli artifact sbom myproject/myapp:latest -o json

# Download SBOM to a file
hrbcli artifact sbom myproject/myapp:latest --file sbom.json
```

#### `hrbcli artifact copy`

Copy artifacts between projects.

```bash
# Copy specific artifact
hrbcli artifact copy myproject/myapp:v1.0 targetproject/myapp:v1.0

# Copy with all tags
hrbcli artifact copy myproject/myapp targetproject/myapp --all-tags
```

### Scanner

#### `hrbcli scanner running`

Show running scans in a project or repository.

```bash
hrbcli scanner running myproject
hrbcli scanner running myproject/myrepo
```

#### `hrbcli scanner scan`

Trigger vulnerability scan for all artifacts in a project or repository.

```bash
# Scan all repositories in project
hrbcli scanner scan myproject

# Scan a single repository
hrbcli scanner scan myproject/myrepo
```



#### `hrbcli scanner reports`

Retrieve vulnerability or SBOM reports for artifacts in a project or repository. When used with `--summary`, displays counts of vulnerabilities by severity for each artifact. Use `--output-dir` to download the reports for all artifacts to a directory. Results are sorted by severity (critical, high, medium, low, total, repository) by default; use `--sort` and `--reverse` to change ordering.


```bash
# Vulnerability summary for a project
hrbcli scanner reports myproject --summary

# SBOM reports for a repository
hrbcli scanner reports myproject/myrepo --type sbom

# Download all vulnerability reports in a project
hrbcli scanner reports myproject --output-dir reports
```

### Label Management

#### `hrbcli label list`

List labels.

```bash
hrbcli label list
```

#### `hrbcli label create`

Create a new label.

```bash
hrbcli label create mylabel --scope g
```

### User Management

#### `hrbcli user list`

List Harbor users.

```bash
# List all users
hrbcli user list

# Search by username
hrbcli user list --search "john"

# List with details
hrbcli user list --detail
```

#### `hrbcli user create`

Create a new user.

```bash
# Create user
hrbcli user create john --email john@example.com --realname "John Doe"

# Create admin user
hrbcli user create admin-user --email admin@example.com --admin

# With specific password
hrbcli user create john --email john@example.com --password "secretpass"
```

#### `hrbcli user delete`

Delete a user.

```bash
# Delete user
hrbcli user delete john

# Force delete
hrbcli user delete john --force
```

### System Administration

#### `hrbcli system info`

Get system information.

```bash
# Get general info
hrbcli system info

# Include storage info
hrbcli system info --with-storage

# Output as YAML
hrbcli system info -o yaml
```

#### `hrbcli system statistics`

Show general Harbor statistics.

```bash
# Display statistics
hrbcli system statistics

# Output as JSON
hrbcli system statistics -o json
```

#### `hrbcli system health`

Check system health.

```bash
# Check overall health
hrbcli system health

# Check specific component
hrbcli system health --component core
hrbcli system health --component jobservice
```

#### `hrbcli system config get`


Get Harbor system configuration.

```bash
# Show all configuration
hrbcli system config get

# Show specific value

hrbcli system config get auth_mode
```

#### `hrbcli system config set`


Update Harbor system configuration.

```bash
hrbcli system config set read_only true

```

#### `hrbcli system gc`

Manage garbage collection.

```bash
# Schedule garbage collection
hrbcli system gc schedule

# Get GC history
hrbcli system gc history

# Get GC job details
hrbcli system gc status <job-id>
```

### Replication

#### `hrbcli replication list`

List replication policies.

```bash
# List all policies
hrbcli replication list

# Filter by name
hrbcli replication list --name-filter "prod*"
```

#### `hrbcli replication create`

Create replication policy.

```bash
# Create push-based replication
hrbcli replication create prod-sync \
  --source myproject \
  --destination https://harbor2.example.com \
  --destination-namespace myproject

# Create pull-based replication
hrbcli replication create prod-pull \
  --source https://harbor2.example.com/myproject \
  --destination myproject \
  --direction pull
```

#### `hrbcli replication get`

Show details of a replication policy.

```bash
hrbcli replication get 1
```

#### `hrbcli replication delete`

Delete a replication policy.

```bash
hrbcli replication delete 1
```

#### `hrbcli replication executions`

List executions of a policy.

```bash
hrbcli replication executions 1
```

#### `hrbcli replication execution`

Show execution statistics.

```bash
hrbcli replication execution 10
```

#### `hrbcli replication logs`

Show logs for an execution.

```bash
hrbcli replication logs 10
```

#### `hrbcli replication statistics`

Aggregate statistics across executions.

```bash
hrbcli replication statistics
```

#### `hrbcli replication execute`

Execute replication manually. The replication policy can be specified either by
ID or by name using `--policy-name`.

```bash
# Execute replication by name
hrbcli replication execute prod-sync

# Execute using the flag
hrbcli replication execute --policy-name prod-sync

# Dry run
hrbcli replication execute prod-sync --dry-run
```

#### `hrbcli distribution providers <project>`

List distribution providers configured for a project.

```bash
hrbcli distribution providers myproject
```

#### `hrbcli distribution policies <project>`

List distribution policies defined in a project.

```bash
hrbcli distribution policies myproject
```

#### `hrbcli distribution policy <project> <name>`

Show details of a distribution policy.

```bash
hrbcli distribution policy myproject mypolicy
```

### Job Service

#### `hrbcli jobservice dashboard`

Display job service worker pools, workers and job queues.

```bash
hrbcli jobservice dashboard
```

### Configuration

#### `hrbcli config init`

Initialize configuration interactively.

```bash
hrbcli config init
```

#### `hrbcli config set`

Set configuration values.

```bash
# Set Harbor URL
hrbcli config set harbor_url https://harbor.example.com

# Set default output format
hrbcli config set output_format json

# Set default project
hrbcli config set default_project library
```

#### `hrbcli config get`

Get configuration values.

```bash
# Get specific value
hrbcli config get harbor_url

# Get all values
hrbcli config get
```

#### `hrbcli config list`

List all configuration.

```bash
hrbcli config list
```

### Shell Completion

Generate completion scripts for your shell.

```bash
# Bash completion
hrbcli completion bash > /etc/bash_completion.d/hrbcli

# Zsh completion
hrbcli completion zsh > _hrbcli
```

### Version Information

Display the CLI version and build details.

```bash
hrbcli version
```

## Examples

### Complete Workflow Examples

#### Setting up a new project

```bash
# Create project
hrbcli project create production --public --storage-limit 100G

# Add user to project
hrbcli project member add production john --role developer

# Create replication from dev to production
hrbcli replication create dev-to-prod \
  --source development \
  --destination production \
  --filter "name=release/*"
```

#### Repository management workflow

```bash
# List repositories
hrbcli repo list myproject

# Check tags
hrbcli repo tags myproject/webapp

# Scan for vulnerabilities
hrbcli artifact scan myproject/webapp:latest

# Copy to production
hrbcli artifact copy myproject/webapp:v1.2.3 production/webapp:v1.2.3

# Clean up old tags
hrbcli repo delete myproject/webapp:old-version
```

#### Security scanning workflow

```bash
# Scan all repositories in project
for repo in $(hrbcli repo list myproject -o json | jq -r '.[].name'); do
  hrbcli artifact scan "$repo" --all
done

# Check scan results
hrbcli artifact list myproject/webapp --with-scan-overview

# List all artifacts in the project with details
hrbcli artifact list myproject --detail

# Export vulnerability report
hrbcli artifact vulnerabilities myproject/webapp:latest -o json > vulns.json
```

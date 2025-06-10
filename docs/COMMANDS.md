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

List artifacts in a repository.

```bash
# List all artifacts
hrbcli artifact list myproject/myapp

# List with vulnerabilities
hrbcli artifact list myproject/myapp --with-scan-overview

# List with labels
hrbcli artifact list myproject/myapp --with-label
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
```

#### `hrbcli artifact vulnerabilities`

Show vulnerability report for an artifact. Use `--severity` to fail if vulnerabilities of that level or higher exist.

```bash
# Show vulnerabilities
hrbcli artifact vulnerabilities myproject/myapp:latest

# Fail if high or critical vulns found
hrbcli artifact vulnerabilities myproject/myapp:latest --severity high
```

#### `hrbcli artifact sbom`

Display the SBOM report for an artifact.

```bash
hrbcli artifact sbom myproject/myapp:latest -o json
```

#### `hrbcli artifact copy`

Copy artifacts between projects.

```bash
# Copy specific artifact
hrbcli artifact copy myproject/myapp:v1.0 targetproject/myapp:v1.0

# Copy with all tags
hrbcli artifact copy myproject/myapp targetproject/myapp --all-tags
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

#### `hrbcli system backup`

Create a backup of Harbor data and configuration. Docker must be available on the host running the command.

```bash
# Backup Harbor using default settings
hrbcli system backup

# Store the backup under /backups
hrbcli system backup --dir /backups
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

Execute replication manually.

```bash
# Execute replication
hrbcli replication execute prod-sync

# Dry run
hrbcli replication execute prod-sync --dry-run
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

# Export vulnerability report
hrbcli artifact vulnerabilities myproject/webapp:latest -o json > vulns.json
```

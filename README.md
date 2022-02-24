```yaml
project:
  name: priority-project
services:
  - priority: 1
    hostname: db
    type: redis@6
    mode: NON_HA
  - priority: 2
    hostname: app
    type: go@1
    mode: NON_HA
    ports:
      - port: 8080
    buildFromGit: https://github.com/zeropsio/test-priority-project-app
    enableSubdomainAccess: true
```
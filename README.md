# test-server

## Examples

### Examples on a local environment

```shell
# Run a server using examples/todo-api.yaml
$ go run main.go -config examples/todo-api.yaml &

# Request to /todo
$ curl -s 'http://localhost:8080/todo' | jq
[
  {
    "id": 1,
    "status": "In Progress",
    "title": "Take A Tour of Go"
  },
  {
    "id": 2,
    "status": "To Do",
    "title": "Read Effective Go"
  }
]

# Request to /todo/1
$ curl -s 'http://localhost:8080/todo/1' | jq
{
  "id": 1,
  "status": "In Progress",
  "title": "Take A Tour of Go"
}

# Request to /todo/2
$ curl -s 'http://localhost:8080/todo/2' | jq
{
  "id": 2,
  "status": "To Do",
  "title": "Read Effective Go"
}
```

### Examples on a Docker environment

```shell
# Set a config file
$ CONFIG_FILE=./examples/todo-api-on-docker.yaml

# Build Dockerfile
$ docker image build -t test-server-on-docker:latest --build-arg CONFIG_FILE=${CONFIG_FILE} -f dockerfiles/test-server.dockerfile .

# Start a docker container with port forwarding to run a server
$ docker container run -p 8080:8080 -d --name test-server test-server-on-docker:latest

# Request to /todo
$ curl -s 'http://localhost:8080/todo' | jq
[
  {
    "id": 1,
    "status": "In Progress",
    "title": "Take A Tour of Go"
  },
  {
    "id": 2,
    "status": "To Do",
    "title": "Read Effective Go"
  }
]

# Request to /todo/1
$ curl -s 'http://localhost:8080/todo/1' | jq
{
  "id": 1,
  "status": "In Progress",
  "title": "Take A Tour of Go"
}

# Request to /todo/2
$ curl -s 'http://localhost:8080/todo/2' | jq
{
  "id": 2,
  "status": "To Do",
  "title": "Read Effective Go"
}
```

### Examples on a Kubernetes environment

```shell
# Set a config file
$ CONFIG_FILE=./examples/todo-api.yaml

# Build Dockerfile
$ docker image build -t test-server-on-kubernetes:latest --build-arg CONFIG_FILE=${CONFIG_FILE} -f dockerfiles/test-server.dockerfile .

# Apply Kubernetes manifests into default namespace
$ kubectl apply -f kubernetes-manifests/

# Execute kubectl proxy
$ kubectl proxy &
[1] nnnn
Starting to serve on 127.0.0.1:8001

# Request to /todo
$ curl -s 'http://localhost:8001/api/v1/namespaces/default/services/test-server-service:8080/proxy/todo' | jq
[
  {
    "id": 1,
    "status": "In Progress",
    "title": "Take A Tour of Go"
  },
  {
    "id": 2,
    "status": "To Do",
    "title": "Read Effective Go"
  }
]

# Request to /todo/1
$ curl -s 'http:/localhost:8001/api/v1/namespaces/default/services/test-server-service:8080/proxy/todo/1' | jq
{
  "id": 1,
  "status": "In Progress",
  "title": "Take A Tour of Go"
}

# Request to /todo/2
$ curl -s 'http://localhost:8001/api/v1/namespaces/default/services/test-server-service:8080/proxy/todo/2' | jq
{
  "id": 2,
  "status": "To Do",
  "title": "Read Effective Go"
}
```

## Unit Tests

```shell
go test .
```

## Acceptance Tests

### Acceptance Tests using testdata/config.yaml on a local environment

```shell
# Run a server using testdata/config.yaml
go run main.go -config testdata/config.yaml

# Activate venv
source venv/bin/activate

# Run acceptance tests
tavern-ci tavern/test_config.tavern.yaml

# Deactivate venv
deactivate
```

### Acceptance Tests using testdata/configWithDefaultValues.yaml

```shell
# Run a server using testdata/configWithDefaultValues.yaml
go run main.go -config testdata/configWithDefaultValues.yaml

# Activate venv
source venv/bin/activate

# Run acceptance tests
tavern-ci tavern/test_config_with_default_values.tavern.yaml

# Deactivate venv
deactivate
```

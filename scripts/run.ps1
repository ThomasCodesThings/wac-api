param (
    $command
)

if (-not $command) {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:DEPARTMENT_API_ENVIRONMENT = "Development"
$env:DEPARTMENT_API_PORT = "8080"
$env:DEPARTMENT_API_MONGODB_USERNAME = "root"
$env:DEPARTMENT_API_MONGODB_PASSWORD = "neUhaDnes"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "start" {
        mongo up --detach
        go run ${ProjectRoot}/cmd/department-api-service
        mongo down
    }
    default {
        throw "Unknown command: $command"
    }
}

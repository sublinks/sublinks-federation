# Sublinks Federation

This is the Federation service for the Sublinks project. It's built using Go.

Together with the [Sublinks Core](https://github.com/sublinks/sublinks-api) and [Sublinks Frontend](https://github.com/sublinks/sublinks-frontend) it's creating a federated link aggregation and microblogging platform.

## Contributing

### Developer Guidelines

[CONTRIBUTING.md](CONTRIBUTING.md)

### Feature Requests / Bugs

Please post any feature requests or bug reports in the repository's [Issues section](https://github.com/sublinks/sublinks-federation/issues).

## Local Dev

### Install pre-requisites:

- `go install golang.org/x/vuln/cmd/govulncheck@latest`
- `go install -v github.com/go-critic/go-critic/cmd/gocritic@latest`
- Install [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
- Install [pre-commit](https://pre-commit.com/#installation)
- `pre-commit install`

### Docker Setup (for dependencies)

To use this Docker setup you first need to complete a few steps:

- Authenticate towards GitHub's container registry(GHCR)
  - Please follow [this guide](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-with-a-personal-access-token-classic) to create an access token and authenticate before running the below command.
- Run `git submodule update --init --recursive` to pull down all submodules
- docker-compose -f docker/docker-compose.federation.yaml up [-d]
  - -d to run in detached mode (don't include the [] brackets, they indicate being optional)

### Run

- Copy .env-sample to .env
- Replace values in .env file
- Run `go run ./cmd/`
- Open [localhost:8080](http://localhost:8080/)

## Running in Production

_**NOTE**: Sublinks is still in early development. This text is added for future use and should not be taken as an indication that Sublinks is ready for use in a real production environment_

A docker image is generated as an artifact of this repository and it is the preferred way of running the service. Environment variables should be passed to the image using standard Docker workflows (using a docker-compose.yaml is the preferred solution)

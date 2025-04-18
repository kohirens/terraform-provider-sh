# Kohirens Terraform Provider Sh

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.10
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
2. cd into the repository
3. Run `docker compose up`
4. Open another terminal and log into the container:
   ```shell
    docker exec -it terraform-provider-sh-terraform-1 sh
   ```
5. Run `go mod tidy .`
6. Run unit and acceptance test with: `TF_ACC=1 go test -count=1 -v ./...`
7. Run the example with:
   ```shell
   go build .
   go install .
   cd examples/provider
   terraform plan
   cd ../data-sources/sh_vars
   terraform plan
   ```
NOTE: You don't run terraform init when in development mode.


## Using the provider

```terraform
terraform {
  required_providers {
    sh = {
      source = "kohirens/sh"
    }
  }
}

provider "sh" {}

# Grab shell environment variables for use in a plan and apply.
data "sh_vars" "example" {
  required = true # throw an error and fail if any values are missing.
  names = [
    "APP_CLIENT_ID",
    "APP_SECRET_ID",
    "APP_SECRET_KEY",
    "CONTENT_URL",
    "HOME",
  ]
  # It would be nice if this module could encrypt them and return that. but it does not do that right now.
}

module "webapp" {
  source = "github.com/kohirens/terraform-aws-web?ref=3.0.0"
  //...
  lf_environment_vars = data.sh_vars.values
}

output "home_dir" {
  value = data.sh_vars.values["HOME"]
}
```

## Developing the Provider

If you wish to work on the provider, the faster method is to have Docker
installed so you can run the container to develop the appl. Or you'll need
[Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named data source page

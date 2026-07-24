# terraform-provider-ketryx

This is a Terraform provider for managing Ketryx. Currently it just manages Ketryx projects.

## Development

In the project root, create a `.env` file with a Ketryx API key enabled for
permissions to manage Ketryx resources.

```
KETRYX_API_KEY=***********************
```

Add the following to `${HOME}/.terraformrc` to make Terraform look for a local
build of this provider.

```
provider_installation {
  dev_overrides {
    "registry.terraform.io/beaconbio/ketryx" = "/path/to/go/bin"
  }

  direct {}
}
```

`/path/to/go/bin` should be where `go install` puts Go binaries. By default
this is `${GOPATH}/bin`.

Afterwards, you can begin development and changes may be tested against
Terraform code after updating the provider binary with `go install` on the
project root.

Task may be used to test code under `examples/provider`.

```
task watch          # Will start a watch loop to build the provider

task example-apply  # Executes `terraform apply`
```

Run `task --list-all` for a full list.

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


## Commits

This project uses conventional commits to automatically generate a changelog and
[semantic version numbers](http://semver.org/). Please make sure all commits to
main adhere to the [spec](https://www.conventionalcommits.org/en/v1.0.0/#specification).

## Releases

A new release PR will be created automatically from releasable commits (eg:
prefixed with `feat:` for a minor version bump or `fix:` for a patch bump) using
[Release Please](https://github.com/googleapis/release-please). Once the PR is
merged, a release will be created and [goreleaser](https://goreleaser.com/) will
build the project for various platforms and add the assets to the release.

## Atlantis

This provider is built into our Atlantis image. If you release a new version,
you should [bump the version](https://github.com/beacon-biosignals/infra/blob/main/beacon-images/atlantis/Dockerfile#L76)
in the Atlantis dockerfile and rebuild the image.

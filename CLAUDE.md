# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Task runner is used for automation (requires [Task](https://taskfile.dev)):

```bash
task watch        # go install with file watching (watches **/*.go, 500ms interval)
task example      # terraform plan && terraform apply in examples/provider/
task example-cleanup  # terraform destroy in examples/provider/
```

Direct Go commands:
```bash
go build ./...
go test ./...
go test ./internal/provider/...  # test a specific package
```

Environment variables are loaded from `.env` (create from `.env.example` if present). `KETRYX_API_KEY` is required for the provider to authenticate.

## Architecture

This is a [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) provider for the Ketryx API, published at `registry.terraform.io/beacon-biosignals/ketryx`.

**Two layers:**

1. **`internal/provider/`** — Terraform-facing layer. Implements the Plugin Framework interfaces (`provider.Provider`, `resource.Resource`, `datasource.DataSource`). Handles schema definition, plan/state marshaling, and calls the API client.

2. **`pkg/ketryx/`** — API client layer. Plain HTTP client for `https://app.ketryx.com`. Owns all request/response structs and API method implementations. Authenticates via `Authorization: Bearer <api_key>`.

**Provider registration** (`provider.go`): `ketryxProvider` reads `api_key` from config or `KETRYX_API_KEY` env var, constructs a `ketryx.Client`, and passes it to resources/data sources via `Configure()`.

**Current resources and data sources:**
- `ketryx_project` resource — Create/Delete/ImportState implemented; Read and Update are stubs
- `ketryx_project` data source — stub (not implemented)
- `ketryx_projects` data source — lists all projects

**Naming convention:** A resource/data source named `ketryx_foo` lives in `internal/provider/foo_resource.go` or `foo_data_source.go`, with corresponding API methods in `pkg/ketryx/client.go`.

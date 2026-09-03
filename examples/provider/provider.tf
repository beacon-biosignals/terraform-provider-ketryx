#
# This is a test space for engaging and using the Ketryx provider.

terraform {
  required_providers {
    ketryx = {
      source = "registry.terraform.io/beacon-biosignals/ketryx"
    }
  }
}

provider "ketryx" {}

data "ketryx_projects" "projects" {}

resource "ketryx_project" "test" {
  name = "test-project"
}

output "ketryx_projects" {
  value = data.ketryx_projects.projects
}

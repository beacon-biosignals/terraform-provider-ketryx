terraform {
  required_providers {
    ketryx = {
      source = "registry.terraform.io/beaconbio/ketryx"
    }
  }
}

provider "ketryx" {}

data "ketryx_projects" "projects" {}

output "ketryx_projects" {
  value = data.ketryx_projects.projects
}

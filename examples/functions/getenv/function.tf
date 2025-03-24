terraform {
  required_providers {
    sh = {
      source = "registry.terraform.io/kohirens/sh"
    }
  }
}

provider "sh" {}

output "home_dir" {
  value = provider::sh::getenv("HOME")
}

terraform {
  required_providers {
    sh = {
      source = "kohirens.com/edu/sh"
    }
  }
}

provider "sh" {}

data "sh_vars" "example" {
  names = [
    "HOME",
    "HOSTNAME",
    "AWS_PROFILE"
  ]
}

output "values" {
  value = data.sh_vars.example.values
}
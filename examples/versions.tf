terraform {
  required_version = ">= 0.13.0"
  required_providers {
    local = "~> 1.2"
    ignition = {
      source  = "e-breuninger/ignition"
      version = "~> 1.0.0"
    }
  }
}

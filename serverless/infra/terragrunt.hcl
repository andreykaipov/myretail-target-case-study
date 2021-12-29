include "root" {
  path = find_in_parent_folders()
}

locals {
  cf_zone_name = get_env("CLOUDFLARE_ZONE_NAME", "")
  username     = get_env("MYRETAIL_USERNAME")
  password     = get_env("MYRETAIL_PASSWORD")
}

inputs = {
  cf_zone_name   = local.cf_zone_name
  myretail_creds = "${local.username}:${local.password}"
}

remote_state {
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }

  backend = "http"
  config = {
    username       = get_env("TF_BACKEND_USERNAME")
    password       = get_env("TF_BACKEND_PASSWORD")
    address        = "https://tf.kaipov.com/myretail"
    lock_address   = "https://tf.kaipov.com/myretail"
    unlock_address = "https://tf.kaipov.com/myretail"
  }
}

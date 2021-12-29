locals {
  dir = get_parent_terragrunt_dir()
}

terraform {
  before_hook "before_hook" {
    commands = ["apply", "plan"]
    execute = ["sh", "-c", <<EOF
      cd "${local.dir}"
      deno bundle src/index.js dist/worker.js
    EOF
    ]
  }
}

inputs = {
  myretail_worker_script = "${local.dir}/dist/worker.js"
}

retry_max_attempts       = 3
retry_sleep_interval_sec = 10

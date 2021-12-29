terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 3.0"
    }
  }
}

variable "myretail_worker_script" {
  type        = string
  description = "The path to our bundled JS worker."
}

variable "myretail_creds" {
  type        = string
  description = "Basic auth credentials our script expects for restricted routes."
}

variable "cf_zone_name" {
  type        = string
  description = "The Cloudflare zone name so we have a custom route for the myRetail API."
  default     = ""
}

resource "cloudflare_workers_kv_namespace" "myretail" {
  title = "myretail"
}

resource "cloudflare_worker_script" "myretail" {
  name    = "myretail"
  content = file(var.myretail_worker_script)

  kv_namespace_binding {
    name         = cloudflare_workers_kv_namespace.myretail.title
    namespace_id = cloudflare_workers_kv_namespace.myretail.id
  }

  secret_text_binding {
    name = "creds"
    text = var.myretail_creds
  }
}

data "cloudflare_zone" "zone" {
  count = var.cf_zone_name == "" ? 0 : 1
  name  = var.cf_zone_name
}

resource "cloudflare_worker_route" "myretail" {
  count       = var.cf_zone_name == "" ? 0 : 1
  zone_id     = data.cloudflare_zone.zone[0].zone_id
  pattern     = "myretail.${data.cloudflare_zone.zone[0].name}/*"
  script_name = cloudflare_worker_script.myretail.name
}

// The following record must exist in our zone for the above custom route to
// work properly. See
// https://developers.cloudflare.com/workers/platform/routes#subdomains-must-have-a-dns-record.
resource "cloudflare_record" "myretail" {
  count   = var.cf_zone_name == "" ? 0 : 1
  zone_id = data.cloudflare_zone.zone[0].zone_id
  name    = "myretail"
  type    = "AAAA"
  value   = "100::"
  proxied = true
  ttl     = 1
}

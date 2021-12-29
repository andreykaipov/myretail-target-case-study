## setup

### development

Local development is via Wrangler. Our JavaScript is built and bundled with
Deno. Deployment is Terraform.

### deployment

Setup our Cloudflare and Terraform backend environment variables. The Terraform
credentials are for my own Terraform backend at `https://tf.kaipov.com`, so one
will either have to remove the `remote_state` stanza from
[`infra/terragrunt.hcl`](infra/terragrunt.hcl) entirely to use local state, or
use a custom backend.

```console
❯ export CLOUDFLARE_API_TOKEN=redacted
❯ export CLOUDFLARE_ACCOUNT_ID=redacted
❯ export TF_BACKEND_USERNAME=redacted
❯ export TF_BACKEND_PASSWORD=redacted
```

Configure some more environment variables. The zone name is optional if you'd
would like to create a custom route for it under an owned Cloudflare zone.
Otherwise, `myRetail` would get deployed to your Worker's subdomain, i.e.
`myretail.your-subdomain.workers.dev`.

```console
❯ export CLOUDFLARE_ZONE_NAME=kaipov.com
❯ export MYRETAIL_USERNAME=admin
❯ export MYRETAIL_PASSWORD=password
```

Then just deploy with Terragrunt. Some of the following Terraform output is
omitted for brevity:

```console
❯ TERRAGRUNT_WORKING_DIR=infra terragrunt apply -auto-approve
INFO[0000] Executing hook: before_hook
Check file:///home/andrey/projects/personal/target.com-case-study/serverless/src/index.js
Bundle file:///home/andrey/projects/personal/target.com-case-study/serverless/src/index.js
Emit "dist/worker.js" (9.17KB)
2021/12/27 13:17:21 [DEBUG] LOCK https://tf.kaipov.com/myretail
Acquiring state lock. This may take a few moments...
2021/12/27 13:17:21 [DEBUG] GET https://tf.kaipov.com/myretail

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # cloudflare_record.myretail[0] will be created
  + resource "cloudflare_record" "myretail" {

  # cloudflare_worker_route.myretail[0] will be created
  + resource "cloudflare_worker_route" "myretail" {

  # cloudflare_worker_script.myretail will be created
  + resource "cloudflare_worker_script" "myretail" {

  # cloudflare_workers_kv_namespace.myretail will be created
  + resource "cloudflare_workers_kv_namespace" "myretail" {

Plan: 4 to add, 0 to change, 0 to destroy.
cloudflare_workers_kv_namespace.myretail: Creating...
cloudflare_record.myretail[0]: Creating...
cloudflare_workers_kv_namespace.myretail: Creation complete after 0s [id=46de4da64abe4612bb306643141b232b]
cloudflare_worker_script.myretail: Creating...
cloudflare_record.myretail[0]: Creation complete after 1s [id=53715b591902be4ef7629274db2e0cb0]
cloudflare_worker_script.myretail: Creation complete after 1s [id=myretail]
cloudflare_worker_route.myretail[0]: Creating...
cloudflare_worker_route.myretail[0]: Creation complete after 1s [id=21cbbe6799034bf9bcbf1f2cb33a1052]
2021/12/27 13:17:24 [DEBUG] POST https://tf.kaipov.com/myretail?ID=d05b2b11-89b8-10ff-e7e0-2c9c35ba6191
2021/12/27 13:17:24 [DEBUG] UNLOCK https://tf.kaipov.com/myretail

Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

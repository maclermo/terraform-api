# terraform-api

## Goal

Runs Terraform plans from a remote server, mimics Terraform Cloud.

## Technical specs

Uses:

1. Go 1.17+
2. Gin-gonic
3. gruntwork.io/terratest

## How to use

Build and run server.

Zip the following file as an example:

```terraform
terraform {
  required_version = ">= 0.12.26"
}

variable "hello" {
  type    = string
  default = "you"
}

output "hello_who" {
  value = "Hi, hi, hi, ${var.hello}!"
}
```

Send a POST to the server with your tf files as a zipped file.

Workspace is MANDATORY!

```bash
# curl --location --request POST "127.0.0.1:8080/apply" \
# --form "terraform=@\"/Users/maclermo/terraform/api-test/main.tf.zip\"" \
# --form "workspace=\"home\"" \
# --form "vars=\"{\\\"hello\\\": \\\"mother\\\"}\""

{
    "id": "9259aa55-475d-4dc8-84fe-cea8cf63cf88"
}

```

You will get a request ID, from which you can get the output:

```bash
# curl --location --request GET "127.0.0.1:8080/output/9259aa55-475d-4dc8-84fe-cea8cf63cf88"

{
    "hello_who": "Hi, hi, hi, mother!"
}
```

## In development

A nice frontend for the API is in the works...

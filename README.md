Terraform Provider Lightstep
==================

This repository is a [Terraform](https://www.terraform.io) provider for [Lightstep](https://lightstep.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.12

Building The Provider
---------------------

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

Using the provider
----------------------
1. Download and install [Terraform](https://www.terraform.io/downloads.html)
2. Move `terraform-provider-lightstep` executable to the root folder where `main.tf` is. 
  
```shell
  $ mv terraform-provider-lightstep /path/to/main.tf/folder
``` 
3. Execute `terraform plan`
4. `terraform init`


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

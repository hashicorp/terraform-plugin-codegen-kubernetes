# Terraform Plugin Framework Code Generator for Kubernetes 

> :warning: This is an experimental project to aid development of the [Terraform Kubernetes Provider](https://github.com/hashicorp/terraform-provider-kubernetes) and is not ready for production use.

## Overview

This repository contains tools for generating code that uses the Terraform Plugin Framework to implement Kubernetes resources. 

This project contains:

- A [code generator](./internal/generator/) that templates Terraform Plugin Framework code.
- An [autocrud](./autocrud) package for encoding between Terraform model types and Kubernetes unstructured objects, and sending them to the Kubernetes API using the client-go dynamic client.  

## Dependencies

This project depends upon having an installation of [hashicorp/terraform-plugin-codegen-openapi](https://github.com/hashicorp/terraform-plugin-codegen-openapi).

## Usage

This tool is used as a binary and can be installed by running `make install` at the top level. 

## License

Refer to [Mozilla Public License v2.0](./LICENSE).

## Experimental Status

By using the software in this repository (the "Software"), you acknowledge that: (1) the Software is still in development, may change, and has not been released as a commercial product by HashiCorp and is not currently supported in any way by HashiCorp; (2) the Software is provided on an "as-is" basis, and may include bugs, errors, or other issues; (3) the Software is NOT INTENDED FOR PRODUCTION USE, use of the Software may result in unexpected results, loss of data, or other unexpected results, and HashiCorp disclaims any and all liability resulting from use of the Software; and (4) HashiCorp reserves all rights to make all decisions about the features, functionality and commercial release (or non-release) of the Software, at any time and without any obligation or liability whatsoever.

Terraform Custom Provider for Tencent Cloud
===========================================

This Terraform custom provider is designed for own use case scenario.

Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- | ---- | ----|
| >= 1.3.x	| 0.1.1	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.3.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Local Installation
------------------

1. Run make file `make install-local-custom-provider` to install the provider under ~/.terraform.d/plugins.

2. The provider source should be change to the path that configured in the *Makefile*:

    ```
    terraform {
      required_providers {
        st-tencentcloud = {
          source = "example.local/myklst/st-tencentcloud"
        }
      }
    }

    provider "st-tencentcloud" {
      region = "ap-hongkong"
    }
    ```

Why Custom Provider
-------------------

This custom provider exists due to some of the resources and data sources in the
official Tencent Cloud Terraform provider may not fulfill the requirements of some
scenario. The reason behind every resources and data sources are stated as below:

### Resources

### Data Sources

- **st-tencentcloud_cloud_load_balancers**

  The tags parameter of Tencent Cloud API
  [*DescribeLoadBalancers*](https://www.tencentcloud.com/document/product/214/1261)
  will return all load balancers when any one of the tags are matched. This may
  be a problem when the user wants to match exactly all given tags, therefore
  this data source will filter once more after listing the load balancers
  from Tencent Cloud API to match all the given tags.

  The example bahaviors of Tencent Cloud API *DescribeLoadBalancers*:

  | Load Balancer   | Tags                                            | Given tags: { "location": "office" "env": "test" }          |
  |-----------------|-------------------------------------------------|-------------------------------------------------------------|
  | load-balancer-A | { "location": "office" "env" : "test" }         | Matched (work as expected)                                  |
  | load-balancer-B | { "location": "office" "env" : "prod" }         | Matched (should not be matched as the `env` is prod)          |

References
----------

- Website: https://www.terraform.io
- Tencent Cloud official Terraform provider: https://github.com/tencentcloudstack/terraform-provider-tencentcloud

- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework

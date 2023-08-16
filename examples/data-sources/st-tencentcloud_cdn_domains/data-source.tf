data "st-tencentcloud_cdn_domains" "cdn_domains" {
  domain = "www.test.com"
}

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

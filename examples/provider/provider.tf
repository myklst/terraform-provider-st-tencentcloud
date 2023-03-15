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

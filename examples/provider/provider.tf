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

resource "st-tencentcloud_cam_group_membership" "test" {
  group_id = 271589
  user_id = 15565624
}

resource "st-tencentcloud_cam_group_membership" "test2" {
  group_id = 271589
  user_id = 15565625
}

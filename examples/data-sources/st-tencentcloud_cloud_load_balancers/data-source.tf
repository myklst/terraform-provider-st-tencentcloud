provider "st-tencentcloud" {
  alias  = "clb"
  region = "ap-hongkong"
}

data "st-tencentcloud_cloud_load_balancers" "clbs" {
  provider = st-tencentcloud.clb

  id   = "load_balancer_id"
  name = "load_balancer_name"
  tags = {
    "app"           = "web"
    "env"           = "basic"
  }
}

output "cloud_load_balancers" {
  value = data.st-tencentcloud_cloud_load_balancers.clbs
}

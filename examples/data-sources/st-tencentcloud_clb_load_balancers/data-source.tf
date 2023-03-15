provider "st-tencentcloud" {
  alias  = "clb"
  region = "ap-hongkong"
}

data "st-tencentcloud_clb_load_balancers" "clbs" {
  provider = st-tencentcloud.clb

  id   = "load_balancer_id"
  name = "load_balancer_name"
  tags = {
    "app" = "web"
    "env" = "basic"
  }
}

output "clb_load_balancers" {
  value = data.st-tencentcloud_clb_load_balancers.clbs
}

data "st-tencentcloud_clb_load_balancers" "clbs" {
  tags = {
    "app" = "crond"
    "env" = "test"
  }
}

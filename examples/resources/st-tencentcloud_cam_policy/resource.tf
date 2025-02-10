resource "st-tencentcloud_cam_policy" "name" {
  user_name         = "devopsuser01"
  attached_policies = ["QcloudAAFullAccess", "QcloudCamFullAccess"]
}

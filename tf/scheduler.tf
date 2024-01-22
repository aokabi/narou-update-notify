resource "google_cloud_scheduler_job" "job" {
  name     = "narou-notify-2"
  description = "narou-notify cloud functionを定期実行する"
  region   = "asia-northeast1"
  time_zone = "Asia/Tokyo"
  schedule = "0 * * * *"

  pubsub_target {
    topic_name = "projects/main-349812/topics/narou-notify"
    data = base64encode("test")
  }
}
resource "google_cloud_scheduler_job" "job" {
  name     = "narou-notify"
  description = "narou-notify cloud functionを定期実行する"
  region   = "asia-northeast1"
  time_zone = "Asia/Tokyo"
  schedule = "0 * * * *"

  retry_config {
    retry_count = 2
    max_retry_duration = "60s"
    min_backoff_duration = "5s"
    max_backoff_duration = "3600s"
    max_doublings = 5
  }

  pubsub_target {
    topic_name = var.topic_name
    data = base64encode("test")
  }
}
resource "google_cloud_scheduler_job" "job" {
  name     = "narou-notify"
  region   = "asia-northeast1"
  schedule = "0 * * * *"

  pubsub_target {
    topic_name = projects/main-349812/topics/narou-notify
  }
}
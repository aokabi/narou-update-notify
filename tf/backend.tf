terraform {
    backend "gcs" {
        bucket = "main-349812-bucket-tfstate"
        prefix = "terraform/state"
    }
}
resource "google_firebase_project" "firebase" {
  provider = google-beta
  project  = var.project_id

  depends_on = [google_project_service.enabled]
}

resource "google_firebase_android_app" "sunrise_alarm" {
  provider     = google-betsunrise_alarm
  project      = var.project_id
  display_name = "Sunrise Alarm"
  package_name = "com.example.sunriseappnew"
}

# Output information for Android app
data "google_firebase_android_app_config" "default" {
  provider   = google-beta
  project    = var.project_id
  app_id     = google_firebase_android_app.default.app_id
}

output "firebase_android_config" {
  value = data.google_firebase_android_app_config.default.config_file_contents
}

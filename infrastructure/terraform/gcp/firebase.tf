resource "google_firebase_project" "firebase" {
  provider   = google-beta
  project    = var.project_id

  depends_on = [google_project_service.enabled]
}

resource "google_firebase_android_app" "sunrise_alarm" {
  provider     = google-beta
  project      = var.project_id
  display_name = "Sunrise Alarm"
  package_name = "com.example.sunriseappnew"

  depends_on = [google_firebase_project.firebase]
}

# Output information for Android app
data "google_firebase_android_app_config" "default" {
  provider   = google-beta
  project    = var.project_id
  app_id     = google_firebase_android_app.sunrise_alarm.app_id
}

# This is the firebase config to be stored in the application files
output "firebase_android_config" {
  value = data.google_firebase_android_app_config.default.config_file_contents
}

# Firebase project ID
# output "firebase_project_id" {
#   value = google_firebase_project.firebase.project_id
# }

# Copies the firebase config to the application level directory
resource "null_resource" "copy_firebase_config" {
  depends_on = [data.google_firebase_android_app_config.default]

  provisioner "local-exec" {
    command = <<-EOT
      terraform output -raw firebase_android_config | base64 -d > "$(git rev-parse --show-toplevel)/app/google-services.json"
    EOT
  }
}

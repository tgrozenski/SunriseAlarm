# This is a service account to enable lambda to access GCP services
resource "google_service_account" "lambda_service_account" {
  account_id   = "lambda-account-832493223"
  display_name = "Lambda Service Account"
}

# This attaches the firebase.admin permission to lambda_service_account
resource "google_project_iam_member" "lambda_firebase_admin" {
  project = var.project_id
  role    = "roles/firebase.admin"
  member  = "serviceAccount:${google_service_account.lambda_service_account.email}"
}

# Generate a key for lambda_service_account
resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.lambda_service_account.name
}

# The key for lambda to use for firebase.admin privilages
output "service_account_key" {
  description = "The service account key"
  sensitive = true
  value = google_service_account_key.mykey.private_key
}

# The service account email
output "service_account_email" {
  value = google_service_account.lambda_service_account.email
}

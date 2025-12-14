# This is a service account to enable lambda to access GCP services
resource "google_service_account" "lambda_service_account" {
  account_id   = "lambda-account-832493223"
  display_name = "Lambda Service Account"
}

resource "google_project_iam_member" "lambda_firebase_admin" {
  project = var.project_id
  role    = "roles/firebase.admin"
  member  = "serviceAccount:${google_service_account.lambda_service_account.email}"
}


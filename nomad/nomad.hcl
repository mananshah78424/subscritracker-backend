job "subscription-checker" {
  datacenters = ["dc1"]
  type        = "batch"

  # Run daily at 2:00 AM UTC
  periodic {
    cron             = "0 2 * * *"
    prohibit_overlap = true
    time_zone        = "UTC"
  }

  group "subscription-check" {
    count = 1

    task "check-subscriptions" {
      driver = "docker"

      config {
        image = "yourusername/subscritracker-backend:latest"
        
        # Command to run the subscription check
        command = "/app/subscription-checker"
      }

      # Environment variables
      env {
        NODE_ENV = "production"
        # Add other required environment variables here
        # GOOGLE_CLIENT_ID = "${GOOGLE_CLIENT_ID}"
        # GOOGLE_CLIENT_SECRET = "${GOOGLE_CLIENT_SECRET}"
        # DB_HOST = "${DB_HOST}"
        # DB_PORT = "${DB_PORT}"
        # DB_USER = "${DB_USER}"
        # DB_PASSWORD = "${DB_PASSWORD}"
        # DB_NAME = "${DB_NAME}"
      }

      # Resource constraints
      resources {
        cpu    = 100
        memory = 128
      }

      # Restart policy
      restart {
        attempts = 3
        interval = "30m"
        delay    = "15s"
        mode     = "fail"
      }

      # Logging
      logging {
        type = "file"
        config {
          path = "/var/log/nomad/subscription-check.log"
        }
      }
    }
  }
}

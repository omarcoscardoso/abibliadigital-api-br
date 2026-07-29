#!/usr/bin/env bash
set -euo pipefail

# Configuration
SERVICE_NAME="abibliadigital-api-br"
REGION="us-central1"
MEMORY="128Mi"
MIN_INSTANCES="0"
MAX_INSTANCES="10"

echo "=== Deploying $SERVICE_NAME to GCP Cloud Run ==="

# Check gcloud CLI availability
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud CLI is not installed."
    echo "Please install Google Cloud SDK: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

PROJECT_ID=$(gcloud config get-value project 2>/dev/null || echo "")
if [ -z "$PROJECT_ID" ]; then
    echo "Error: GCP Project ID not configured. Run 'gcloud config set project <PROJECT_ID>'"
    exit 1
fi

IMAGE_URI="gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest"

echo "1. Building Docker image: $IMAGE_URI"
gcloud builds submit --tag "$IMAGE_URI" .

echo "2. Deploying to Cloud Run ($SERVICE_NAME)"
gcloud run deploy "$SERVICE_NAME" \
    --image "$IMAGE_URI" \
    --platform managed \
    --region "$REGION" \
    --memory "$MEMORY" \
    --min-instances "$MIN_INSTANCES" \
    --max-instances "$MAX_INSTANCES" \
    --allow-unauthenticated

echo "=== Deployment completed successfully! ==="

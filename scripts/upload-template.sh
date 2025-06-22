#!/bin/bash

# Script to upload templates to the template registry
# Usage: ./upload-tempecho "Template upload completed successfully!"
echo "Template: $TEMPLATE_NAME"
echo "Version: $VERSION"
echo "Description: $DESCRIPTION"
echo "Category: $CATEGORY"
echo "Features: $FEATURES"
if [ -n "$PREVIEW_IMAGE_PATH" ]; then
    echo "Preview image: $PREVIEW_IMAGE_PATH"
fi [template_name] [version] [description] [preview_image_path] [category] [features]

set -e

TEMPLATE_NAME=${1:-"template_1"}
VERSION=${2:-$(date +"%Y%m%d-%H%M%S")}
DESCRIPTION=${3:-"Template upload via script"}
PREVIEW_IMAGE_PATH=${4:-""}
CATEGORY=${5:-"ecommerce"}
FEATURES=${6:-"responsive,modern,sveltekit"}
TEMPLATE_REGISTRY_URL=${TEMPLATE_REGISTRY_URL:-"http://localhost:9001"}

echo "Uploading template: $TEMPLATE_NAME version $VERSION"
echo "Category: $CATEGORY"
echo "Features: $FEATURES"
if [ -n "$PREVIEW_IMAGE_PATH" ]; then
    echo "Preview image: $PREVIEW_IMAGE_PATH"
fi

# Change to the template directory
TEMPLATE_DIR="./templates/$TEMPLATE_NAME"
if [ ! -d "$TEMPLATE_DIR" ]; then
    echo "Error: Template directory $TEMPLATE_DIR not found"
    exit 1
fi

cd "$TEMPLATE_DIR"

# Check if build directory exists
if [ ! -d "build" ]; then
    echo "Error: Build directory not found. Please run 'npm run build' first"
    exit 1
fi

# Create a temporary tarball
TEMP_FILE=$(mktemp /tmp/template-upload-XXXXXX.tar.gz)
echo "Creating tarball..."
tar -czf "$TEMP_FILE" build

# Validate preview image if provided
if [ -n "$PREVIEW_IMAGE_PATH" ]; then
    if [ ! -f "$PREVIEW_IMAGE_PATH" ]; then
        echo "Error: Preview image file not found: $PREVIEW_IMAGE_PATH"
        rm "$TEMP_FILE"
        exit 1
    fi
    
    # Check if it's an image file by extension
    case "${PREVIEW_IMAGE_PATH,,}" in
        *.jpg|*.jpeg|*.png|*.gif|*.webp)
            echo "Preview image validated: $PREVIEW_IMAGE_PATH"
            ;;
        *)
            echo "Warning: Preview image should be a common image format (jpg, png, gif, webp)"
            ;;
    esac
fi

# Upload to template registry
echo "Uploading to template registry..."
if [ -n "$PREVIEW_IMAGE_PATH" ]; then
    # Upload with preview image
    curl -X POST \
        -F "template_name=$TEMPLATE_NAME" \
        -F "version=$VERSION" \
        -F "description=$DESCRIPTION" \
        -F "category=$CATEGORY" \
        -F "features=$FEATURES" \
        -F "assets=@$TEMP_FILE" \
        -F "preview_image=@$PREVIEW_IMAGE_PATH" \
        "$TEMPLATE_REGISTRY_URL/templates/upload"
else
    # Upload without preview image
    curl -X POST \
        -F "template_name=$TEMPLATE_NAME" \
        -F "version=$VERSION" \
        -F "description=$DESCRIPTION" \
        -F "category=$CATEGORY" \
        -F "features=$FEATURES" \
        -F "assets=@$TEMP_FILE" \
        "$TEMPLATE_REGISTRY_URL/templates/upload"
fi

# Clean up
rm "$TEMP_FILE"

echo "Template upload completed successfully!"
echo "Template: $TEMPLATE_NAME"
echo "Version: $VERSION"
echo "Description: $DESCRIPTION"
if [ -n "$PREVIEW_IMAGE_PATH" ]; then
    echo "Preview image: $PREVIEW_IMAGE_PATH"
fi

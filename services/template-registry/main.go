package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	s3Client   *s3.Client
	bucketName string
	publicURL  string
	rootCtx    = context.TODO()
	logger     *zap.Logger
)

// TemplateRegistry manages pre-built template assets and metadata
type TemplateRegistry struct {
	TemplateName string   `json:"template_name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Features     []string `json:"features"`
}

// TemplateUploadRequest represents a template upload request
type TemplateUploadRequest struct {
	TemplateName string `json:"template_name" form:"template_name"`
	Version      string `json:"version" form:"version"`
	Description  string `json:"description" form:"description"`
	Force        bool   `json:"force" form:"force"`
	Assets       string `json:"assets" form:"assets" swaggertype:"string" format:"binary"`
	PreviewImage string `json:"preview_image,omitempty" form:"preview_image" swaggertype:"string" format:"binary"`
}

// TemplateVersion represents a template version with metadata
type TemplateVersion struct {
	Version     string    `json:"version"`
	Description string    `json:"description"`
	UploadTime  time.Time `json:"upload_time"`
	AssetCount  int       `json:"asset_count"`
	TotalSize   int64     `json:"total_size"`
	GitCommit   string    `json:"git_commit,omitempty"`
	BuildID     string    `json:"build_id"`
	Status      string    `json:"status"` // "available", "uploading", "failed"
	Manifest    string    `json:"manifest"`
}

// TemplateManifest contains metadata about template assets
type TemplateManifest struct {
	Version      string            `json:"version"`
	TemplateName string            `json:"template_name"`
	Description  string            `json:"description"`
	UploadTime   time.Time         `json:"upload_time"`
	Assets       []AssetInfo       `json:"assets"`
	TotalSize    int64             `json:"total_size"`
	AssetCount   int               `json:"asset_count"`
	GitCommit    string            `json:"git_commit,omitempty"`
	BuildID      string            `json:"build_id"`
	Checksum     string            `json:"checksum"`
	Metadata     map[string]string `json:"metadata"`
	ThumbnailURL string            `json:"thumbnail_url,omitempty"`
}

// AssetInfo represents information about a single asset
type AssetInfo struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	Checksum     string    `json:"checksum"`
	LastModified time.Time `json:"last_modified"`
}

// Template represents a template with its metadata for API responses
type Template struct {
	Name         string    `json:"name"`
	Title        string    `json:"title,omitempty"`
	Version      string    `json:"version"`
	Description  string    `json:"description"`
	Category     string    `json:"category,omitempty"`
	Features     []string  `json:"features,omitempty"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	PreviewURL   *string   `json:"preview_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func init() {
	// Initialize zap logger (production vs development)
	var err error
	if os.Getenv("ENVIRONMENT") == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	// Load environment variables (non-fatal if missing .env)
	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found, proceeding with existing environment")
	}

	accessKeyID := os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("CLOUDFLARE_R2_ENDPOINT")
	bucketName = os.Getenv("CLOUDFLARE_R2_BUCKET_NAME")
	publicURL = os.Getenv("CLOUDFLARE_R2_PUBLIC_URL")

	if accessKeyID == "" || secretAccessKey == "" || endpoint == "" || bucketName == "" || publicURL == "" {
		logger.Fatal("missing required Cloudflare R2 environment variables",
			zap.Bool("has_access_key", accessKeyID != ""),
			zap.Bool("has_secret", secretAccessKey != ""),
			zap.Bool("has_endpoint", endpoint != ""),
			zap.Bool("has_bucket", bucketName != ""),
			zap.Bool("has_public_url", publicURL != ""),
		)
	}

	cfg, err := config.LoadDefaultConfig(rootCtx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.Fatal("failed to load AWS config", zap.Error(err))
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	logger.Info("template registry initialized", zap.String("bucket", bucketName))
}

func main() {
	defer func() { _ = logger.Sync() }()

	r := mux.NewRouter()

	// Basic request logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			// Wrap ResponseWriter to capture status
			rw := &statusRecorder{ResponseWriter: w, status: 200}
			next.ServeHTTP(rw, req)
			latency := time.Since(start)
			logger.Info("http_request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", rw.status),
				zap.Duration("latency", latency),
				zap.String("remote_ip", req.RemoteAddr),
				zap.String("user_agent", req.UserAgent()),
			)
		})
	})

	// Template registry endpoints
	r.HandleFunc("/templates", listTemplatesHandler).Methods("GET")
	r.HandleFunc("/templates/{template_name}", getTemplateHandler).Methods("GET")
	r.HandleFunc("/templates/{template_name}/versions", getTemplateVersionsHandler).Methods("GET")
	r.HandleFunc("/templates/{template_name}/latest", getLatestVersionHandler).Methods("GET")
	r.HandleFunc("/templates/{template_name}/versions/{version}", getVersionHandler).Methods("GET")
	r.HandleFunc("/templates/upload", uploadTemplateHandler).Methods("POST")
	r.HandleFunc("/templates/{template_name}/versions/{version}/download", downloadTemplateHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	logger.Info("starting template registry service", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatal("server exited", zap.Error(err))
	}
}

// statusRecorder helps capture HTTP status codes for logging
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// uploadTemplateHandler handles template upload requests
// @Summary Upload a template
// @Description Upload a new template with assets and optional preview image
// @Tags templates
// @Accept multipart/form-data
// @Produce json
// @Param template_name formData string true "Template name"
// @Param version formData string false "Template version (auto-generated if not provided)"
// @Param description formData string false "Template description"
// @Param category formData string false "Template category (e.g., web, mobile, desktop)"
// @Param features formData string false "Template features (comma-separated list)"
// @Param force formData boolean false "Force upload even if version exists"
// @Param assets formData file true "Template assets archive (tar.gz)"
// @Param preview_image formData file false "Preview image for template (PNG, JPG, WebP, GIF)"
// @Success 200 {object} map[string]interface{} "Template uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /templates/upload [post]
func uploadTemplateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	templateName := r.FormValue("template_name")
	version := r.FormValue("version")
	description := r.FormValue("description")
	category := r.FormValue("category")
	featuresStr := r.FormValue("features")
	force := r.FormValue("force") == "true"

	// Parse features from comma-separated string
	var features []string
	if featuresStr != "" {
		rawFeatures := strings.Split(featuresStr, ",")
		for _, feature := range rawFeatures {
			feature = strings.TrimSpace(feature)
			if feature != "" {
				features = append(features, feature)
			}
		}
	}

	if templateName == "" {
		http.Error(w, "Template name is required", http.StatusBadRequest)
		return
	}

	if version == "" {
		version = generateVersionFromTime()
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("assets")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the preview image (optional)
	var previewImageFile multipart.File
	var previewImageHandler *multipart.FileHeader
	previewImageFile, previewImageHandler, err = r.FormFile("preview_image")
	hasPreviewImage := err == nil && previewImageFile != nil
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Failed to get preview image file", http.StatusBadRequest)
		return
	}
	if previewImageFile != nil {
		defer previewImageFile.Close()
	}

	logger.Info("uploading template", zap.String("template", templateName), zap.String("version", version), zap.String("filename", handler.Filename), zap.Bool("force", force))
	if hasPreviewImage {
		logger.Info("preview image included", zap.String("file", previewImageHandler.Filename))
	}

	// Check if template version already exists (unless forced)
	if !force {
		exists, err := templateVersionExists(templateName, version)
		if err != nil {
			logger.Warn("failed to check existing template version", zap.Error(err), zap.String("template", templateName), zap.String("version", version))
		} else if exists {
			writeJSONResponse(w, map[string]interface{}{
				"status":  "skipped",
				"message": "Template version already exists",
				"version": version,
			})
			return
		}
	}

	// Upload template to R2
	registry := &TemplateRegistry{
		TemplateName: templateName,
		Version:      version,
		Description:  description,
		Category:     category,
		Features:     features,
	}

	if err := registry.UploadTemplate(file, handler.Size, previewImageFile, previewImageHandler); err != nil {
		logger.Error("template upload failed", zap.Error(err), zap.String("template", templateName), zap.String("version", version))
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Update latest pointer
	if err := updateLatestPointer(templateName, version); err != nil {
		logger.Warn("failed to update latest version pointer", zap.Error(err), zap.String("template", templateName), zap.String("version", version))
		// Don't fail the upload for this
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":        "success",
		"message":       "Template uploaded successfully",
		"template_name": templateName,
		"version":       version,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	})
}

// UploadTemplate uploads a template archive to R2
func (tr *TemplateRegistry) UploadTemplate(file io.Reader, size int64, previewImageFile multipart.File, previewImageHandler *multipart.FileHeader) error {
	logger.Info("starting template upload", zap.String("template", tr.TemplateName), zap.String("version", tr.Version))

	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "template_upload_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract archive to temp directory
	if err := extractArchive(file, tempDir); err != nil {
		return fmt.Errorf("failed to extract archive: %v", err)
	}

	// Find build directory in extracted files
	buildDir := filepath.Join(tempDir, "build")
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return fmt.Errorf("build directory not found in uploaded archive")
	}

	// Generate manifest
	manifest, err := tr.generateManifest(buildDir)
	if err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	// Upload preview image if provided
	if previewImageFile != nil && previewImageHandler != nil {
		thumbnailURL, err := tr.uploadPreviewImage(previewImageFile, previewImageHandler)
		if err != nil {
			logger.Warn("failed to upload preview image", zap.Error(err), zap.String("template", tr.TemplateName), zap.String("version", tr.Version))
			// Don't fail the whole upload for preview image issues
		} else {
			manifest.ThumbnailURL = thumbnailURL
			logger.Info("preview image uploaded", zap.String("url", thumbnailURL), zap.String("template", tr.TemplateName), zap.String("version", tr.Version))
		}
	}

	// Upload assets to R2
	if err := tr.uploadAssetsToR2(buildDir, manifest); err != nil {
		return fmt.Errorf("failed to upload assets to R2: %w", err)
	}

	logger.Info("template upload complete", zap.String("template", tr.TemplateName), zap.String("version", tr.Version))
	return nil
}

// extractArchive extracts a tar.gz archive to the specified directory
func extractArchive(src io.Reader, dest string) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
	}

	return nil
}

// generateManifest generates a manifest for the template
func (tr *TemplateRegistry) generateManifest(buildDir string) (*TemplateManifest, error) {
	var assets []AssetInfo
	var totalSize int64

	err := filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(buildDir, path)

		// Calculate checksum
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return err
		}
		checksum := hex.EncodeToString(hasher.Sum(nil))

		// Determine content type
		contentType := "application/octet-stream"
		if ext := filepath.Ext(path); ext != "" {
			switch ext {
			case ".html":
				contentType = "text/html"
			case ".css":
				contentType = "text/css"
			case ".js":
				contentType = "application/javascript"
			case ".json":
				contentType = "application/json"
			case ".png":
				contentType = "image/png"
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".svg":
				contentType = "image/svg+xml"
			case ".woff", ".woff2":
				contentType = "font/woff"
			}
		}

		assets = append(assets, AssetInfo{
			Path:         strings.ReplaceAll(relPath, "\\", "/"),
			Size:         info.Size(),
			ContentType:  contentType,
			Checksum:     checksum,
			LastModified: info.ModTime(),
		})

		totalSize += info.Size()
		return nil
	})

	if err != nil {
		return nil, err
	}

	buildID := fmt.Sprintf("build_%s_%s_%d", tr.TemplateName, tr.Version, time.Now().Unix())

	manifest := &TemplateManifest{
		Version:      tr.Version,
		TemplateName: tr.TemplateName,
		Description:  tr.Description,
		UploadTime:   time.Now(),
		Assets:       assets,
		TotalSize:    totalSize,
		AssetCount:   len(assets),
		BuildID:      buildID,
		Metadata:     make(map[string]string),
	}

	// Store category and features in metadata
	if tr.Category != "" {
		manifest.Metadata["category"] = tr.Category
	}
	if len(tr.Features) > 0 {
		manifest.Metadata["features"] = strings.Join(tr.Features, ",")
	}
	// Store title based on template name
	manifest.Metadata["title"] = formatTemplateTitle(tr.TemplateName)

	// Calculate overall checksum
	manifestJSON, _ := json.Marshal(manifest.Assets)
	hasher := sha256.New()
	hasher.Write(manifestJSON)
	manifest.Checksum = hex.EncodeToString(hasher.Sum(nil))

	return manifest, nil
}

// uploadAssetsToR2 uploads template assets to R2
func (tr *TemplateRegistry) uploadAssetsToR2(buildDir string, manifest *TemplateManifest) error {
	// Upload individual assets
	for _, asset := range manifest.Assets {
		filePath := filepath.Join(buildDir, asset.Path)
		s3Key := fmt.Sprintf("%s/%s/%s", tr.TemplateName, tr.Version, asset.Path)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open asset %s: %v", asset.Path, err)
		}

		_, err = s3Client.PutObject(rootCtx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(s3Key),
			Body:        file,
			ContentType: aws.String(asset.ContentType),
		})
		file.Close()

		if err != nil {
			return fmt.Errorf("failed to upload asset %s: %v", asset.Path, err)
		}

		logger.Debug("uploaded asset", zap.String("key", s3Key))
	}

	// Upload manifest
	manifestJSON, _ := json.Marshal(manifest)
	manifestKey := fmt.Sprintf("%s/%s/manifest.json", tr.TemplateName, tr.Version)

	_, err := s3Client.PutObject(rootCtx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(manifestKey),
		Body:        strings.NewReader(string(manifestJSON)),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		return fmt.Errorf("failed to upload manifest: %v", err)
	}

	logger.Info("manifest uploaded", zap.String("key", manifestKey), zap.String("template", tr.TemplateName), zap.String("version", tr.Version))
	return nil
}

// uploadPreviewImage uploads a preview image to R2 and returns the public URL
func (tr *TemplateRegistry) uploadPreviewImage(previewImageFile multipart.File, previewImageHandler *multipart.FileHeader) (string, error) {
	// Determine content type based on file extension
	contentType := "image/jpeg" // default
	if ext := filepath.Ext(previewImageHandler.Filename); ext != "" {
		switch strings.ToLower(ext) {
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".webp":
			contentType = "image/webp"
		case ".gif":
			contentType = "image/gif"
		}
	}

	// Generate S3 key for the preview image
	// Use timestamp to avoid conflicts and ensure uniqueness
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(previewImageHandler.Filename)
	s3Key := fmt.Sprintf("template-previews/%s/%s/preview_%d%s", tr.TemplateName, tr.Version, timestamp, fileExt)

	// Upload the preview image
	_, err := s3Client.PutObject(rootCtx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(s3Key),
		Body:        previewImageFile,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload preview image: %v", err)
	}

	// Construct the public URL
	// Use the public CDN URL directly
	publicPreviewURL := fmt.Sprintf("%s/%s", publicURL, s3Key)

	logger.Info("uploaded preview image", zap.String("key", s3Key))
	return publicPreviewURL, nil
}

// listTemplatesHandler lists all available templates
// @Summary List all templates
// @Description Get a list of all available templates with metadata
// @Tags templates
// @Produce json
// @Success 200 {array} Template "Templates retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /templates [get]
func listTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := listAvailableTemplatesWithMetadata()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Return templates array directly to match Swagger documentation
	writeJSONResponse(w, templates)
}

// getTemplateHandler gets information about a specific template
func getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	versions, err := getTemplateVersions(templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get template: %v", err), http.StatusInternalServerError)
		return
	}

	if len(versions) == 0 {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":        "success",
		"template_name": templateName,
		"versions":      versions,
		"version_count": len(versions),
	})
}

// getTemplateVersionsHandler gets all versions for a template
func getTemplateVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	versions, err := getTemplateVersions(templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get template versions: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"versions": versions,
		"count":    len(versions),
	})
}

// getLatestVersionHandler gets the latest version for a template
func getLatestVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	version, err := getLatestVersion(templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get latest version: %v", err), http.StatusInternalServerError)
		return
	}

	if version == nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":  "success",
		"version": version,
	})
}

// getVersionHandler gets a specific version of a template
func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]
	version := vars["version"]

	manifest, err := getTemplateManifest(templateName, version)
	if err != nil {
		http.Error(w, "Template version not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"manifest": manifest,
	})
}

// downloadTemplateHandler provides download functionality for templates
func downloadTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]
	version := vars["version"]

	// Get template manifest to verify it exists
	manifest, err := getTemplateManifest(templateName, version)
	if err != nil {
		http.Error(w, "Template version not found", http.StatusNotFound)
		return
	}

	// For now, just return the manifest and asset list
	// In a real implementation, you might want to create a zip file
	writeJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"template": templateName,
		"version":  version,
		"manifest": manifest,
		"message":  "Use individual asset URLs to download assets",
	})
}

// Helper functions

func templateVersionExists(templateName, version string) (bool, error) {
	manifestKey := fmt.Sprintf("%s/%s/manifest.json", templateName, version)

	_, err := s3Client.HeadObject(rootCtx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(manifestKey),
	})

	if err != nil {
		// Check if error is "not found"
		return false, nil
	}

	return true, nil
}

func updateLatestPointer(templateName, version string) error {
	latestKey := fmt.Sprintf("%s/latest", templateName)

	_, err := s3Client.PutObject(rootCtx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(latestKey),
		Body:        strings.NewReader(version),
		ContentType: aws.String("text/plain"),
	})

	return err
}

func listAvailableTemplates() ([]string, error) {
	prefix := ""

	resp, err := s3Client.ListObjectsV2(rootCtx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, err
	}

	var templates []string
	for _, commonPrefix := range resp.CommonPrefixes {
		templateName := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
		templateName = strings.TrimSuffix(templateName, "/")
		if templateName != "" && templateName != "template-previews" {
			templates = append(templates, templateName)
		}
	}

	sort.Strings(templates)
	return templates, nil
}

func listAvailableTemplatesWithMetadata() ([]Template, error) {
	prefix := ""

	resp, err := s3Client.ListObjectsV2(rootCtx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, err
	}

	var templates []Template
	for _, commonPrefix := range resp.CommonPrefixes {
		templateName := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
		templateName = strings.TrimSuffix(templateName, "/")
		if templateName != "" && templateName != "template-previews" {
			// Get latest version manifest to retrieve metadata
			latestManifest, err := getLatestManifest(templateName)
			if err != nil {
				// If no manifest available, create basic template entry
				templates = append(templates, Template{
					Name:        templateName,
					Title:       formatTemplateTitle(templateName),
					Version:     "unknown",
					Description: "No description available",
					Category:    "web",
					Features:    []string{"responsive"},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				continue
			}

			// Extract thumbnail URL and preview URL if available
			var thumbnailURL *string
			var previewURL *string
			if latestManifest.ThumbnailURL != "" {
				thumbnailURL = &latestManifest.ThumbnailURL
				previewURL = &latestManifest.ThumbnailURL // Use same URL for preview for now
			}

			// Create template object with metadata matching Swagger schema
			template := Template{
				Name:         templateName,
				Title:        formatTemplateTitle(templateName),
				Version:      latestManifest.Version,
				Description:  latestManifest.Description,
				Category:     "web",                            // Default category
				Features:     []string{"responsive", "modern"}, // Default features
				ThumbnailURL: thumbnailURL,
				PreviewURL:   previewURL,
				CreatedAt:    latestManifest.UploadTime,
				UpdatedAt:    latestManifest.UploadTime,
			}

			// Override with metadata from template manifest if available
			if category, exists := latestManifest.Metadata["category"]; exists && category != "" {
				template.Category = category
			}
			if featuresStr, exists := latestManifest.Metadata["features"]; exists && featuresStr != "" {
				// Parse features from comma-separated string
				features := strings.Split(featuresStr, ",")
				for i, feature := range features {
					features[i] = strings.TrimSpace(feature)
				}
				template.Features = features
			}
			if title, exists := latestManifest.Metadata["title"]; exists && title != "" {
				template.Title = title
			}

			templates = append(templates, template)
		}
	}

	// Sort templates by name
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

func getTemplateVersions(templateName string) ([]TemplateVersion, error) {
	prefix := fmt.Sprintf("%s/", templateName)

	resp, err := s3Client.ListObjectsV2(rootCtx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, err
	}

	var versions []TemplateVersion
	for _, commonPrefix := range resp.CommonPrefixes {
		versionDir := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
		versionDir = strings.TrimSuffix(versionDir, "/")

		if versionDir == "" || versionDir == "latest" {
			continue
		}

		// Try to get manifest for this version
		manifest, err := getTemplateManifest(templateName, versionDir)
		if err != nil {
			// If no manifest, create a basic version entry
			versions = append(versions, TemplateVersion{
				Version:    versionDir,
				Status:     "available",
				UploadTime: time.Now(), // This would ideally come from object metadata
			})
			continue
		}

		versions = append(versions, TemplateVersion{
			Version:     manifest.Version,
			Description: manifest.Description,
			UploadTime:  manifest.UploadTime,
			AssetCount:  manifest.AssetCount,
			TotalSize:   manifest.TotalSize,
			GitCommit:   manifest.GitCommit,
			BuildID:     manifest.BuildID,
			Status:      "available",
			Manifest:    fmt.Sprintf("%s/%s/manifest.json", templateName, versionDir),
		})
	}

	// Sort versions by upload time (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].UploadTime.After(versions[j].UploadTime)
	})

	return versions, nil
}

func getLatestVersion(templateName string) (*TemplateVersion, error) {
	versions, err := getTemplateVersions(templateName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, nil
	}

	return &versions[0], nil
}

func getTemplateManifest(templateName, version string) (*TemplateManifest, error) {
	manifestKey := fmt.Sprintf("%s/%s/manifest.json", templateName, version)

	resp, err := s3Client.GetObject(rootCtx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(manifestKey),
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest TemplateManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func getLatestManifest(templateName string) (*TemplateManifest, error) {
	// First try to get the latest version
	latestKey := fmt.Sprintf("%s/latest", templateName)

	resp, err := s3Client.GetObject(rootCtx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(latestKey),
	})

	if err != nil {
		// If no latest pointer, get the most recent version
		versions, err := getTemplateVersions(templateName)
		if err != nil || len(versions) == 0 {
			return nil, fmt.Errorf("no versions found for template %s", templateName)
		}
		// versions are already sorted by upload time (newest first)
		return getTemplateManifest(templateName, versions[0].Version)
	}
	defer resp.Body.Close()

	// Read the version from the latest pointer
	versionBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	latestVersion := strings.TrimSpace(string(versionBytes))
	return getTemplateManifest(templateName, latestVersion)
}

func generateVersionFromTime() string {
	return fmt.Sprintf("v%s", time.Now().Format("20060102-150405"))
}

// formatTemplateTitle converts a template name to a human-readable title
func formatTemplateTitle(templateName string) string {
	// Replace underscores and hyphens with spaces
	title := strings.ReplaceAll(templateName, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")

	// Capitalize each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, map[string]interface{}{
		"status":    "healthy",
		"service":   "template-registry",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

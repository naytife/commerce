package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	s3Client            *s3.Client
	storesBucketName    string
	templatesBucketName string
	ctx                 = context.Background()
)

type StoreDeployer struct {
	ShopID       string `json:"shop_id"`
	Subdomain    string `json:"subdomain"`
	TemplateName string `json:"template_name"`
	Version      string `json:"version,omitempty"` // Empty for latest
}

type DeploymentRequest struct {
	ShopID       string            `json:"shop_id"`
	Subdomain    string            `json:"subdomain"`
	TemplateName string            `json:"template_name"`
	Version      string            `json:"version,omitempty"`
	DataOverride map[string]string `json:"data_override,omitempty"`
}

type DeploymentResponse struct {
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	ShopID     string    `json:"shop_id"`
	Subdomain  string    `json:"subdomain"`
	Template   string    `json:"template"`
	Version    string    `json:"version"`
	URL        string    `json:"url"`
	DeployedAt time.Time `json:"deployed_at"`
	AssetCount int       `json:"asset_count"`
	TotalSize  int64     `json:"total_size"`
	DeployTime string    `json:"deploy_time"`
}

type StoreData struct {
	Shop     ShopInfo    `json:"shop"`
	Products []Product   `json:"products"`
	Settings StoreConfig `json:"settings"`
}

type ShopInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Subdomain    string `json:"subdomain"`
	Currency     string `json:"currency"`
	Description  string `json:"description"`
	Logo         string `json:"logo"`
	ContactEmail string `json:"contact_email"`
	Theme        string `json:"theme"`
}

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	Slug        string   `json:"slug"`
	Available   bool     `json:"available"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
}

type StoreConfig struct {
	PaymentMethods []string          `json:"payment_methods"`
	ShippingZones  []string          `json:"shipping_zones"`
	TaxSettings    map[string]string `json:"tax_settings"`
	Analytics      map[string]string `json:"analytics"`
}

type TemplateManifest struct {
	Version   string            `json:"version"`
	BuildTime time.Time         `json:"build_time"`
	GitCommit string            `json:"git_commit"`
	Assets    []Asset           `json:"assets"`
	Metadata  map[string]string `json:"metadata"`
}

type Asset struct {
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Hash         string `json:"hash"`
	ContentType  string `json:"content_type"`
	CacheControl string `json:"cache_control"`
}

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Initialize R2 client
	accessKey := strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY"))
	endpoint := strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_ENDPOINT"))
	storesBucketName = strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"))
	templatesBucketName = strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_TEMPLATES_BUCKET_NAME"))

	if accessKey == "" || secretKey == "" || endpoint == "" || storesBucketName == "" || templatesBucketName == "" {
		log.Fatal("Missing required R2 environment variables")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	log.Printf("Store deployer initialized with stores bucket: %s, templates bucket: %s", storesBucketName, templatesBucketName)
}

func main() {
	r := mux.NewRouter()

	// Store deployment endpoints
	r.HandleFunc("/deploy", deployStoreHandler).Methods("POST")
	r.HandleFunc("/redeploy/{subdomain}", redeployStoreHandler).Methods("POST")
	r.HandleFunc("/update-data/{subdomain}", updateDataHandler).Methods("POST")
	r.HandleFunc("/status/{subdomain}", getDeploymentStatusHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9003"
	}

	log.Printf("Store deployer service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func (sd *StoreDeployer) DeployStore() (*DeploymentResponse, error) {
	startTime := time.Now()
	log.Printf("Starting deployment for shop %s (%s) using template %s", sd.ShopID, sd.Subdomain, sd.TemplateName)

	// 1. Get template version to use
	version, err := sd.resolveTemplateVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template version: %v", err)
	}
	sd.Version = version

	// 2. Get template manifest
	manifest, err := sd.getTemplateManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get template manifest: %v", err)
	}

	// 3. Fetch store data from backend
	storeData, err := sd.fetchStoreData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store data: %v", err)
	}

	// 4. Copy template assets to store path
	if err := sd.copyTemplateAssets(manifest); err != nil {
		return nil, fmt.Errorf("failed to copy template assets: %v", err)
	}

	// 5. Generate and upload store data files
	if err := sd.uploadStoreData(storeData); err != nil {
		return nil, fmt.Errorf("failed to upload store data: %v", err)
	}

	deployTime := time.Since(startTime)

	// 6. Create deployment response
	response := &DeploymentResponse{
		Status:     "success",
		Message:    "Store deployed successfully",
		ShopID:     sd.ShopID,
		Subdomain:  sd.Subdomain,
		Template:   sd.TemplateName,
		Version:    sd.Version,
		URL:        fmt.Sprintf("https://%s.naytife.com", sd.Subdomain),
		DeployedAt: time.Now().UTC(),
		AssetCount: len(manifest.Assets),
		TotalSize:  calculateTotalSize(manifest.Assets),
		DeployTime: deployTime.String(),
	}

	log.Printf("Successfully deployed store %s in %v", sd.Subdomain, deployTime)
	return response, nil
}

func (sd *StoreDeployer) resolveTemplateVersion() (string, error) {
	if sd.Version != "" {
		// Specific version requested
		return sd.Version, nil
	}

	// Get latest version from template registry
	templateRegistryURL := os.Getenv("TEMPLATE_REGISTRY_URL")
	if templateRegistryURL == "" {
		templateRegistryURL = "http://template-registry:9001"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/templates/%s/latest", templateRegistryURL, sd.TemplateName))
	if err != nil {
		return "", fmt.Errorf("failed to contact template registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("template not found: %s", sd.TemplateName)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	version, ok := response["version"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format from template registry")
	}

	versionStr, ok := version["version"].(string)
	if !ok {
		return "", fmt.Errorf("invalid version format from template registry")
	}

	return versionStr, nil
}

func (sd *StoreDeployer) getTemplateManifest() (*TemplateManifest, error) {
	manifestKey := fmt.Sprintf("%s/%s/manifest.json", sd.TemplateName, sd.Version)

	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(templatesBucketName),
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

func (sd *StoreDeployer) fetchStoreData() (*StoreData, error) {
	// In a real implementation, this would call the backend GraphQL API
	// For now, return mock data structure

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8002"
	}

	// TODO: Implement GraphQL query to fetch:
	// - Shop details
	// - Products with variants
	// - Store configuration/settings
	// - Payment methods
	// - Shipping configuration

	// Mock data for now
	storeData := &StoreData{
		Shop: ShopInfo{
			ID:           sd.ShopID,
			Name:         fmt.Sprintf("Store %s", sd.Subdomain),
			Subdomain:    sd.Subdomain,
			Currency:     "USD",
			Description:  "Amazing store description",
			Logo:         "",
			ContactEmail: fmt.Sprintf("contact@%s.com", sd.Subdomain),
			Theme:        sd.TemplateName,
		},
		Products: []Product{},
		Settings: StoreConfig{
			PaymentMethods: []string{"stripe", "paypal"},
			ShippingZones:  []string{"US", "CA"},
			TaxSettings:    map[string]string{"enabled": "true"},
			Analytics:      map[string]string{"ga_id": ""},
		},
	}

	return storeData, nil
}

func (sd *StoreDeployer) copyTemplateAssets(manifest *TemplateManifest) error {
	log.Printf("Copying %d template assets for store %s", len(manifest.Assets), sd.Subdomain)

	templatePath := fmt.Sprintf("%s/%s", sd.TemplateName, sd.Version)
	storePath := fmt.Sprintf("%s", sd.Subdomain)

	// Copy each asset from template to store location
	for _, asset := range manifest.Assets {
		sourceKey := fmt.Sprintf("%s/%s", templatePath, asset.Path)
		destKey := fmt.Sprintf("%s/%s", storePath, asset.Path)

		// Copy object from templates bucket to stores bucket
		_, err := s3Client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:      aws.String(storesBucketName),
			Key:         aws.String(destKey),
			CopySource:  aws.String(fmt.Sprintf("%s/%s", templatesBucketName, sourceKey)),
			ContentType: &asset.ContentType,
			// Set appropriate cache control for static assets
			CacheControl: aws.String("public, max-age=31536000, immutable"), // 1 year cache for static assets
		})

		if err != nil {
			return fmt.Errorf("failed to copy asset %s: %v", asset.Path, err)
		}

		log.Printf("Copied asset: %s -> %s", sourceKey, destKey)
	}

	log.Printf("Successfully copied template assets to store path")
	return nil
}

func (sd *StoreDeployer) uploadStoreData(storeData *StoreData) error {
	log.Printf("Uploading store data for %s", sd.Subdomain)

	// Generate data files
	dataFiles := map[string]interface{}{
		"shop.json":     storeData.Shop,
		"products.json": storeData.Products,
		"settings.json": storeData.Settings,
		"metadata.json": map[string]interface{}{
			"template_name":    sd.TemplateName,
			"template_version": sd.Version,
			"deployed_at":      time.Now().UTC().Format(time.RFC3339),
			"last_updated":     time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Upload each data file
	for filename, data := range dataFiles {
		if err := sd.uploadDataFile(filename, data); err != nil {
			return fmt.Errorf("failed to upload %s: %v", filename, err)
		}
	}

	log.Printf("Successfully uploaded store data files")
	return nil
}

func (sd *StoreDeployer) uploadDataFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/data/%s", sd.Subdomain, filename)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(storesBucketName),
		Key:          aws.String(key),
		Body:         strings.NewReader(string(jsonData)),
		ContentType:  aws.String("application/json"),
		CacheControl: aws.String("no-cache, no-store, must-revalidate"), // Data files should not be cached
	})

	return err
}

// HTTP Handlers

func deployStoreHandler(w http.ResponseWriter, r *http.Request) {
	var req DeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShopID == "" || req.Subdomain == "" || req.TemplateName == "" {
		http.Error(w, "Missing required fields: shop_id, subdomain, template_name", http.StatusBadRequest)
		return
	}

	deployer := &StoreDeployer{
		ShopID:       req.ShopID,
		Subdomain:    req.Subdomain,
		TemplateName: req.TemplateName,
		Version:      req.Version,
	}

	response, err := deployer.DeployStore()
	if err != nil {
		log.Printf("Deployment failed: %v", err)
		http.Error(w, fmt.Sprintf("Deployment failed: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, response)
}

func redeployStoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subdomain := vars["subdomain"]

	var req struct {
		TemplateName string `json:"template_name,omitempty"`
		Version      string `json:"version,omitempty"`
		ShopID       string `json:"shop_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShopID == "" {
		http.Error(w, "Missing required field: shop_id", http.StatusBadRequest)
		return
	}

	// Get current deployment info if template not specified
	if req.TemplateName == "" {
		currentInfo, err := getCurrentDeploymentInfo(subdomain)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get current deployment info: %v", err), http.StatusInternalServerError)
			return
		}
		req.TemplateName = currentInfo["template_name"].(string)
	}

	deployer := &StoreDeployer{
		ShopID:       req.ShopID,
		Subdomain:    subdomain,
		TemplateName: req.TemplateName,
		Version:      req.Version,
	}

	response, err := deployer.DeployStore()
	if err != nil {
		log.Printf("Redeployment failed: %v", err)
		http.Error(w, fmt.Sprintf("Redeployment failed: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, response)
}

func updateDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subdomain := vars["subdomain"]

	var req struct {
		ShopID string `json:"shop_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShopID == "" {
		http.Error(w, "Missing required field: shop_id", http.StatusBadRequest)
		return
	}

	// Get current deployment info
	currentInfo, err := getCurrentDeploymentInfo(subdomain)
	if err != nil {
		http.Error(w, fmt.Sprintf("Store not found or not deployed: %v", err), http.StatusNotFound)
		return
	}

	deployer := &StoreDeployer{
		ShopID:       req.ShopID,
		Subdomain:    subdomain,
		TemplateName: currentInfo["template_name"].(string),
		Version:      currentInfo["template_version"].(string),
	}

	// Fetch fresh store data
	storeData, err := deployer.fetchStoreData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch store data: %v", err), http.StatusInternalServerError)
		return
	}

	// Upload only data files (no asset copying needed)
	if err := deployer.uploadStoreData(storeData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update store data: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":     "success",
		"message":    "Store data updated successfully",
		"subdomain":  subdomain,
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func getDeploymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subdomain := vars["subdomain"]

	status, err := getDeploymentStatus(subdomain)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get deployment status: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, status)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, map[string]interface{}{
		"status":    "healthy",
		"service":   "store-deployer",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Helper functions

func getCurrentDeploymentInfo(subdomain string) (map[string]interface{}, error) {
	metadataKey := fmt.Sprintf("%s/data/metadata.json", subdomain)

	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(storesBucketName),
		Key:    aws.String(metadataKey),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metadata map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func getDeploymentStatus(subdomain string) (map[string]interface{}, error) {
	metadata, err := getCurrentDeploymentInfo(subdomain)
	if err != nil {
		return nil, err
	}

	// Check if store is accessible
	isAccessible := checkStoreAccessibility(subdomain)

	status := map[string]interface{}{
		"subdomain":        subdomain,
		"status":           "deployed",
		"template_name":    metadata["template_name"],
		"template_version": metadata["template_version"],
		"deployed_at":      metadata["deployed_at"],
		"last_updated":     metadata["last_updated"],
		"accessible":       isAccessible,
		"url":              fmt.Sprintf("https://%s.naytife.com", subdomain),
	}

	return status, nil
}

func checkStoreAccessibility(subdomain string) bool {
	// Check if main index.html exists
	indexKey := fmt.Sprintf("%s/index.html", subdomain)

	_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(storesBucketName),
		Key:    aws.String(indexKey),
	})

	return err == nil
}

func calculateTotalSize(assets []Asset) int64 {
	var total int64
	for _, asset := range assets {
		total += asset.Size
	}
	return total
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

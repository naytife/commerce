package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	s3Client            *s3.Client
	storesBucketName    string
	templatesBucketName string
	imagesBucketName    string
	httpClient          = &http.Client{Timeout: 15 * time.Second}
	rootCtx             = context.TODO()
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
	imagesBucketName = strings.TrimSpace(os.Getenv("CLOUDFLARE_R2_IMAGES_BUCKET_NAME"))

	if accessKey == "" || secretKey == "" || endpoint == "" || storesBucketName == "" || templatesBucketName == "" || imagesBucketName == "" {
		log.Fatal("Missing required R2 environment variables")
	}

	cfg, err := config.LoadDefaultConfig(rootCtx,
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
	r.HandleFunc("/cleanup/{subdomain}", cleanupStoreHandler).Methods("DELETE")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
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
		templateRegistryURL = "http://template-registry:8002"
	}

	reqURL := fmt.Sprintf("%s/templates/%s/latest", templateRegistryURL, sd.TemplateName)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
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

	resp, err := s3Client.GetObject(rootCtx, &s3.GetObjectInput{
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

// Change fetchStoreData to return map[string]interface{}
func (sd *StoreDeployer) fetchStoreData() (map[string]interface{}, error) {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8000"
	}
	backendURL = strings.TrimRight(backendURL, "/") + "/query"

	var wg sync.WaitGroup
	var shopResp, productsResp map[string]interface{}
	var shopErr, productsErr error

	wg.Add(2)

	// Shop query goroutine
	go func() {
		defer wg.Done()
		log.Printf("[DEBUG] About to send shop GraphQL request to %s", backendURL)
		shopQuery := `query GetShop { shop { id title defaultDomain contactPhone contactEmail address { address } whatsAppNumber whatsAppLink facebookLink instagramLink images { siteLogo { url altText } siteLogoDark { url altText } favicon { url altText } banner { url altText } bannerDark { url altText } coverImage { url altText } coverImageDark { url altText } } currencyCode about shopProductsCategory seoDescription seoKeywords seoTitle paymentMethods { id name provider enabled config { publishableKey testMode } } categories(first: 100) { edges { node { id slug title description images { banner { url altText } } } } } } }`
		shopReq := map[string]interface{}{
			"query":     shopQuery,
			"variables": map[string]interface{}{},
		}
		shopResp, shopErr = doGraphQLRequest(backendURL, shopReq, sd.Subdomain)
		if shopErr != nil {
			log.Printf("[ERROR] shop GraphQL request failed: %v", shopErr)
		} else {
			log.Printf("[DEBUG] shop GraphQL request succeeded")
		}
	}()

	// Products query goroutine
	go func() {
		defer wg.Done()
		log.Printf("[DEBUG] About to send products GraphQL request to %s", backendURL)
		productsQuery := `query GetProducts($first: Int) { products(first: $first) { edges { node { id productId slug title description attributes { title value } defaultVariant { id variationId price availableQuantity description isDefault attributes { title value } stockStatus } variants { id variationId price availableQuantity description isDefault attributes { title value } stockStatus } images { url altText } updatedAt createdAt } } pageInfo { hasNextPage endCursor } totalCount } }`
		productsReq := map[string]interface{}{
			"query":     productsQuery,
			"variables": map[string]interface{}{"first": 100},
		}
		productsResp, productsErr = doGraphQLRequest(backendURL, productsReq, sd.Subdomain)
		if productsErr != nil {
			log.Printf("[ERROR] products GraphQL request failed: %v", productsErr)
		} else {
			log.Printf("[DEBUG] products GraphQL request succeeded")
		}
	}()

	wg.Wait()

	if shopErr != nil {
		return nil, fmt.Errorf("failed to fetch shop: %w", shopErr)
	}
	if productsErr != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", productsErr)
	}

	return map[string]interface{}{
		"shop":     shopResp,
		"products": productsResp,
	}, nil
}

// Helper to do a GraphQL POST and parse JSON response
func doGraphQLRequest(url string, reqBody map[string]interface{}, subdomain string) (map[string]interface{}, error) {
	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if subdomain != "" {
		req.Header.Set("X-Shop-Subdomain", subdomain)
	}

	// Inject tracing headers and request id when available (best-effort)
	// Note: this package doesn't depend on observability to avoid circular imports; use httpClient as shared client.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result struct {
		Data   map[string]interface{} `json:"data"`
		Errors interface{}            `json:"errors"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("GraphQL HTTP %d, body: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("failed to decode GraphQL response: %w", err)
	}
	if result.Errors != nil {
		fmt.Printf("GraphQL error: %v, HTTP %d, body: %s\n", result.Errors, resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("GraphQL error: %v", result.Errors)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("GraphQL HTTP error: %d, body: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("GraphQL HTTP error: %d", resp.StatusCode)
	}
	return result.Data, nil
}

func (sd *StoreDeployer) copyTemplateAssets(manifest *TemplateManifest) error {
	log.Printf("Copying %d template assets for store %s", len(manifest.Assets), sd.Subdomain)

	templatePath := fmt.Sprintf("%s/%s", sd.TemplateName, sd.Version)
	storePath := sd.Subdomain

	// Copy each asset from template to store location
	for _, asset := range manifest.Assets {
		sourceKey := fmt.Sprintf("%s/%s", templatePath, asset.Path)
		destKey := fmt.Sprintf("%s/%s", storePath, asset.Path)

		// Copy object from templates bucket to stores bucket
		_, err := s3Client.CopyObject(rootCtx, &s3.CopyObjectInput{
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

// Update uploadStoreData to accept map[string]interface{} and write raw data files
func (sd *StoreDeployer) uploadStoreData(storeData map[string]interface{}) error {
	log.Printf("Uploading store data for %s", sd.Subdomain)

	// Write raw data files as expected by the frontend
	dataFiles := map[string]interface{}{
		"shop.json":     storeData["shop"],
		"products.json": storeData["products"],
		// settings.json and metadata.json can be left as before or empty
		"settings.json": map[string]interface{}{},
		"metadata.json": map[string]interface{}{
			"template_name":    sd.TemplateName,
			"template_version": sd.Version,
			"deployed_at":      time.Now().UTC().Format(time.RFC3339),
			"last_updated":     time.Now().UTC().Format(time.RFC3339),
		},
	}

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

	_, err = s3Client.PutObject(rootCtx, &s3.PutObjectInput{
		Bucket:       aws.String(storesBucketName),
		Key:          aws.String(key),
		Body:         strings.NewReader(string(jsonData)),
		ContentType:  aws.String("application/json"),
		CacheControl: aws.String("no-cache, no-store, must-revalidate"), // Data files should not be cached
	})
	if err != nil {
		log.Printf("[ERROR] Failed to upload %s to R2: %v", key, err)
	} else {
		log.Printf("[DEBUG] Successfully uploaded %s to R2", key)
	}
	return err
}

// updateSelectiveData updates only specific data files based on the data type
func (sd *StoreDeployer) updateSelectiveData(dataType string) error {
	log.Printf("Updating selective data for %s (type: %s)", sd.Subdomain, dataType)

	// Fetch store data based on the specific type
	var dataToUpdate interface{}
	var filename string

	switch dataType {
	case "shop":
		// Fetch only shop data
		shopData, err := sd.fetchShopDataOnly()
		if err != nil {
			return fmt.Errorf("failed to fetch shop data: %v", err)
		}
		dataToUpdate = shopData
		filename = "shop.json"

	case "products":
		// Fetch only products data
		productsData, err := sd.fetchProductsDataOnly()
		if err != nil {
			return fmt.Errorf("failed to fetch products data: %v", err)
		}
		dataToUpdate = productsData
		filename = "products.json"

	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}

	// Upload the specific data file
	if err := sd.uploadDataFile(filename, dataToUpdate); err != nil {
		return fmt.Errorf("failed to upload %s: %v", filename, err)
	}

	// Update metadata with last updated timestamp
	metadata := map[string]interface{}{
		"template_name":    sd.TemplateName,
		"template_version": sd.Version,
		"last_updated":     time.Now().UTC().Format(time.RFC3339),
		"last_update_type": dataType,
	}
	if err := sd.uploadDataFile("metadata.json", metadata); err != nil {
		log.Printf("Warning: failed to update metadata: %v", err)
	}

	log.Printf("Successfully updated %s data", dataType)
	return nil
}

// fetchShopDataOnly fetches only shop information
func (sd *StoreDeployer) fetchShopDataOnly() (interface{}, error) {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8000"
	}
	backendURL = strings.TrimRight(backendURL, "/") + "/query"

	log.Printf("[DEBUG] About to send shop GraphQL request to %s (fetchShopDataOnly)", backendURL)

	shopQuery := map[string]interface{}{
		"query": `
			query GetShop {
				shop {
					id
					title
					defaultDomain
					contactPhone
					contactEmail
					address {
						address
					}
					whatsAppNumber
					whatsAppLink
					facebookLink
					instagramLink
					images {
						siteLogo {
							url
							altText
						}
						siteLogoDark {
							url
							altText
						}
						favicon {
							url
							altText
						}
						banner {
							url
							altText
						}
						bannerDark {
							url
							altText
						}
						coverImage {
							url
							altText
						}
						coverImageDark {
							url
							altText
						}
					}
					currencyCode
					about
					seoDescription
					seoKeywords
					seoTitle
				}
			}
		`,
		"variables": map[string]interface{}{},
	}

	result, err := doGraphQLRequest(backendURL, shopQuery, sd.Subdomain)
	if err != nil {
		log.Printf("[ERROR] shop GraphQL request failed (fetchShopDataOnly): %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] shop GraphQL request succeeded (fetchShopDataOnly)")

	return result["shop"], nil
}

func (sd *StoreDeployer) fetchProductsDataOnly() (interface{}, error) {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8000"
	}
	backendURL = strings.TrimRight(backendURL, "/") + "/query"

	log.Printf("[DEBUG] About to send products GraphQL request to %s (fetchProductsDataOnly)", backendURL)

	productsQuery := map[string]interface{}{
		"query": `
			query GetProducts($first: Int) {
				products(first: $first) {
					edges {
						node {
							id
							productId
							slug
							title
							description
							attributes {
								title
								value
							}
							defaultVariant {
								id
								variationId
								price
								availableQuantity
								description
								isDefault
								attributes {
									title
									value
								}
								stockStatus
							}
							variants {
								id
								variationId
								price
								availableQuantity
								description
								isDefault
								attributes {
									title
									value
								}
								stockStatus
							}
							images {
								url
								altText
							}
							updatedAt
							createdAt
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
					totalCount
				}
			}
		`,
		"variables": map[string]interface{}{
			"first": 100,
		},
	}

	result, err := doGraphQLRequest(backendURL, productsQuery, sd.Subdomain)
	if err != nil {
		log.Printf("[ERROR] products GraphQL request failed (fetchProductsDataOnly): %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] products GraphQL request succeeded (fetchProductsDataOnly)")

	return result, nil
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
		ShopID   string `json:"shop_id"`
		DataType string `json:"data_type,omitempty"` // "shop", "products", "all", or empty (defaults to "all")
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShopID == "" {
		http.Error(w, "Missing required field: shop_id", http.StatusBadRequest)
		return
	}

	// Default to updating all data if not specified
	if req.DataType == "" {
		req.DataType = "all"
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

	// Handle selective data updates
	var updatedFiles []string
	if req.DataType == "all" {
		// Fetch all store data and update all files
		storeData, err := deployer.fetchStoreData()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch store data: %v", err), http.StatusInternalServerError)
			return
		}

		// Upload all data files
		if err := deployer.uploadStoreData(storeData); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update store data: %v", err), http.StatusInternalServerError)
			return
		}
		updatedFiles = []string{"shop.json", "products.json", "settings.json", "metadata.json"}
	} else {
		// Selective update based on data type
		if err := deployer.updateSelectiveData(req.DataType); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update %s data: %v", req.DataType, err), http.StatusInternalServerError)
			return
		}

		switch req.DataType {
		case "shop":
			updatedFiles = []string{"shop.json"}
		case "products":
			updatedFiles = []string{"products.json"}
		default:
			http.Error(w, fmt.Sprintf("Invalid data_type: %s. Valid values are 'shop', 'products', or 'all'", req.DataType), http.StatusBadRequest)
			return
		}
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":        "success",
		"message":       fmt.Sprintf("Store data updated successfully (%s)", req.DataType),
		"subdomain":     subdomain,
		"data_type":     req.DataType,
		"updated_files": updatedFiles,
		"updated_at":    time.Now().UTC().Format(time.RFC3339),
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

func cleanupStoreHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Starting cleanup for shop %s (subdomain: %s)", req.ShopID, subdomain)

	// Perform the cleanup
	err := cleanupStoreFiles(subdomain, req.ShopID)
	if err != nil {
		log.Printf("Cleanup failed for shop %s: %v", req.ShopID, err)
		// Don't return error - log it but continue
		// This ensures that database deletion isn't blocked by R2 issues
		writeJSONResponse(w, map[string]interface{}{
			"status":     "partial_success",
			"message":    fmt.Sprintf("Cleanup completed with warnings: %v", err),
			"subdomain":  subdomain,
			"shop_id":    req.ShopID,
			"cleaned_at": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	writeJSONResponse(w, map[string]interface{}{
		"status":     "success",
		"message":    "Store files cleaned up successfully",
		"subdomain":  subdomain,
		"shop_id":    req.ShopID,
		"cleaned_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// cleanupStoreFiles removes all R2 files associated with a store
func cleanupStoreFiles(subdomain, shopID string) error {
	var errors []string

	// 1. Delete deployed store files (stored under subdomain/ in stores bucket)
	if subdomain != "" {
		log.Printf("Cleaning up store files for subdomain: %s", subdomain)
		if err := deleteDirectoryInBucket(storesBucketName, subdomain+"/"); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete store files for subdomain %s: %v", subdomain, err))
		}
	}

	// 2. Delete shop images (stored under shops/{shopID}/images/ in images bucket)
	shopImagesPrefix := fmt.Sprintf("shops/%s/images/", shopID)
	log.Printf("Cleaning up shop images for shop: %s", shopID)
	if err := deleteDirectoryInBucket(imagesBucketName, shopImagesPrefix); err != nil {
		errors = append(errors, fmt.Sprintf("failed to delete shop images: %v", err))
	}

	// 3. Delete product images for this shop (stored under products/shop_{shopId}/ in images bucket)
	log.Printf("Cleaning up product images for shop: %s", shopID)
	if err := cleanupProductImagesByShop(shopID); err != nil {
		errors = append(errors, fmt.Sprintf("failed to delete product images: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errors, "; "))
	}

	log.Printf("Successfully cleaned up R2 files for shop %s (subdomain: %s)", shopID, subdomain)
	return nil
}

// deleteDirectoryInBucket deletes all objects under a given prefix in a specific bucket
func deleteDirectoryInBucket(bucketName, prefix string) error {
	// List all objects with the given prefix in the specified bucket
	objects, err := listObjectsInBucket(bucketName, prefix)
	if err != nil {
		return fmt.Errorf("failed to list objects with prefix %s in bucket %s: %v", prefix, bucketName, err)
	}

	if len(objects) == 0 {
		log.Printf("No objects found with prefix: %s in bucket: %s", prefix, bucketName)
		return nil
	}

	// Delete objects in batches (max 1000 per batch)
	const batchSize = 1000
	for i := 0; i < len(objects); i += batchSize {
		end := i + batchSize
		if end > len(objects) {
			end = len(objects)
		}

		batch := objects[i:end]
		if err := deleteObjectsBatchInBucket(bucketName, batch); err != nil {
			return fmt.Errorf("failed to delete batch in bucket %s: %v", bucketName, err)
		}

		log.Printf("Deleted %d objects from R2 bucket %s", len(batch), bucketName)
	}

	log.Printf("Successfully deleted %d objects with prefix: %s from bucket: %s", len(objects), prefix, bucketName)
	return nil
}

// deleteDirectoryInBucket deletes all objects under a given prefix in a specific bucket
func listObjectsInBucket(bucketName, prefix string) ([]string, error) {
	var objects []string
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(s3Client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(rootCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in bucket %s: %v", bucketName, err)
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				objects = append(objects, *obj.Key)
			}
		}
	}

	return objects, nil
}

// deleteObjectsBatchInBucket deletes a batch of objects from a specific bucket
func deleteObjectsBatchInBucket(bucketName string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	var objectsToDelete []types.ObjectIdentifier
	for _, key := range keys {
		objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	_, err := s3Client.DeleteObjects(rootCtx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: objectsToDelete,
		},
	})

	return err
}

// cleanupProductImagesByShop removes product images that belong to a specific shop
// Uses the new shop-aware pattern: products/shop_{shopId}/
func cleanupProductImagesByShop(shopID string) error {
	log.Printf("Starting product image cleanup for shop %s", shopID)

	// Delete shop-aware product images: products/shop_{shopId}/
	shopProductsPrefix := fmt.Sprintf("products/shop_%s/", shopID)
	err := deleteDirectoryInBucket(imagesBucketName, shopProductsPrefix)
	if err != nil {
		log.Printf("Failed to delete product images for shop %s: %v", shopID, err)
		return fmt.Errorf("failed to delete product images for shop %s: %v", shopID, err)
	}

	log.Printf("Successfully cleaned up product images for shop %s", shopID)
	return nil
}

// Helper functions

func getCurrentDeploymentInfo(subdomain string) (map[string]interface{}, error) {
	metadataKey := fmt.Sprintf("%s/data/metadata.json", subdomain)

	resp, err := s3Client.GetObject(rootCtx, &s3.GetObjectInput{
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

	_, err := s3Client.HeadObject(rootCtx, &s3.HeadObjectInput{
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

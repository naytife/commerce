package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var redisClient *redis.Client
var s3Client *s3.Client
var ctx = context.Background()

const redisQueue = "build-queue"
const buildDir = "./built_sites"
const templatesDir = "./templates"
const bucketName = "naytife-shops-static"

type BuildJob struct {
	SiteName     string `json:"site_name"`
	TemplateName string `json:"template_name"`
}

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default Redis address
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Read Cloudflare R2 credentials from environment variables
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	endpoint := os.Getenv("R2_ENDPOINT")
	region := os.Getenv("R2_BUCKET_REGION")

	if accessKey == "" || secretKey == "" || endpoint == "" {
		log.Fatal("Missing R2 credentials in environment variables or .env file")
	}

	// Load default AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create S3 client with BaseEndpoint
	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // Required for Cloudflare R2
	})

	fmt.Println("Cloudflare R2 client initialized successfully")
}

func main() {
	for {
		msg, err := redisClient.BLPop(ctx, 0*time.Second, redisQueue).Result()
		if err != nil {
			log.Println("Error fetching from Redis queue:", err)
			continue
		}

		if len(msg) < 2 {
			continue
		}

		var job BuildJob
		if err := json.Unmarshal([]byte(msg[1]), &job); err != nil {
			log.Println("Invalid JSON payload:", err)
			continue
		}

		log.Printf("Processing build job: Site=%s, Template=%s", job.SiteName, job.TemplateName)

		if err := processJob(job); err != nil {
			log.Println("Failed to process job:", err)
		}
	}
}

func processJob(job BuildJob) error {
	templatePath := filepath.Join(templatesDir, job.TemplateName)
	outputPath := filepath.Join(buildDir, job.SiteName)
	backendURL := fmt.Sprintf("http://%s.localhost:8080/api/query", job.SiteName)

	// Build the Svelte site
	if err := buildSvelteSite(templatePath, outputPath, backendURL); err != nil {
		return fmt.Errorf("failed to build site: %v", err)
	}

	// Upload to Cloudflare R2
	if err := uploadToR2(outputPath, job.SiteName); err != nil {
		return fmt.Errorf("failed to upload site: %v", err)
	}

	// Cleanup the build directory after successful upload
	if err := cleanupBuild(outputPath); err != nil {
		return fmt.Errorf("failed to cleanup build directory: %v", err)
	}

	log.Printf("Successfully built, uploaded, and cleaned up site: %s", job.SiteName)
	return nil
}

func buildSvelteSite(templatePath, outputPath, backendURL string) error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = templatePath
	cmd.Env = append(os.Environ(), fmt.Sprintf("VITE_API_URL=%s", backendURL))

	log.Println("Building site at:", templatePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Build error: %s", string(output))
		return err
	}

	// Check if the output directory already exists
	if _, err := os.Stat(outputPath); err == nil {
		// Remove the existing directory
		log.Printf("Removing existing directory: %s", outputPath)
		if err := os.RemoveAll(outputPath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// Move built files to output path
	if err := os.Rename(filepath.Join(templatePath, "build"), outputPath); err != nil {
		return fmt.Errorf("failed to move built files: %v", err)
	}

	log.Printf("Successfully built site at: %s", outputPath)
	return nil
}

func uploadToR2(directory, siteName string) error {
	// Step 1: Get the list of existing files in R2
	existingFiles, err := listFilesInR2(siteName)
	if err != nil {
		return fmt.Errorf("failed to list existing files: %v", err)
	}

	// Step 2: Prepare a set of new files from the local build
	newFiles := make(map[string]struct{})

	err = filepath.Walk(directory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(directory, filePath)
		key := fmt.Sprintf("%s/%s", siteName, strings.ReplaceAll(relPath, "\\", "/"))
		newFiles[key] = struct{}{} // Store the new file path

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get MIME type and remove charset=utf-8 if present
		contentType := mime.TypeByExtension(path.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		} else {
			contentType = strings.Split(contentType, ";")[0]
		}

		// Upload the file
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			return fmt.Errorf("failed to upload %s: %v", key, err)
		}

		log.Printf("Uploaded %s with Content-Type: %s", key, contentType)
		return nil
	})

	if err != nil {
		return err
	}

	// Step 3: Delete old files that are no longer in the new build
	for _, file := range existingFiles {
		if _, exists := newFiles[file]; !exists {
			log.Printf("Deleting stale file: %s", file)
			_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(file),
			})
			if err != nil {
				log.Printf("Failed to delete %s: %v", file, err)
			}
		}
	}

	return nil
}

// listFilesInR2 retrieves the list of existing files in the site's folder on R2
func listFilesInR2(siteName string) ([]string, error) {
	var files []string
	prefix := fmt.Sprintf("%s/", siteName) // Only list files under the site's folder

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	for {
		result, err := s3Client.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Contents {
			files = append(files, *item.Key)
		}

		if result.IsTruncated != nil && !*result.IsTruncated {
			break
		}

		input.ContinuationToken = result.NextContinuationToken
	}

	return files, nil
}

// cleanupBuild removes the build directory after successful upload
func cleanupBuild(directory string) error {
	if err := os.RemoveAll(directory); err != nil {
		return fmt.Errorf("failed to remove build directory %s: %v", directory, err)
	}
	log.Printf("Cleaned up build directory: %s", directory)
	return nil
}

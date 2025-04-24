/**
 * Uploads a file to Cloudflare R2 storage
 * @param file The file to upload
 * @param productId The product ID to organize files
 * @returns Object containing the uploaded file URL and filename
 */
export async function uploadToR2(file: File, productId: string | number): Promise<{ url: string; filename: string }> {
  try {
    // Generate a unique filename to avoid collisions
    const uniquePrefix = Date.now().toString();
    const filename = `products/${productId}/${uniquePrefix}-${file.name.replace(/\s+/g, '-')}`;
    
    // Create FormData for the file upload
    const formData = new FormData();
    formData.append('file', file);
    formData.append('filename', filename);
    
    // Send the file to your backend API endpoint that handles the R2 upload
    const response = await fetch('/api/upload-to-r2', {
      method: 'POST',
      body: formData
    });
    
    if (!response.ok) {
      const error = await response.text();
      throw new Error(`Failed to upload image: ${error}`);
    }
    
    const data = await response.json();
    return {
      url: data.url,
      filename: filename
    };
  } catch (error) {
    console.error('Error uploading to R2:', error);
    throw error;
  }
}

/**
 * Deletes a file from Cloudflare R2 storage
 * @param filename The filename to delete
 * @returns Success message
 */
export async function deleteFromR2(filename: string): Promise<{ success: boolean; message: string }> {
  try {
    const response = await fetch('/api/delete-from-r2', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ filename })
    });
    
    if (!response.ok) {
      const error = await response.text();
      throw new Error(`Failed to delete image: ${error}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error deleting from R2:', error);
    throw error;
  }
} 
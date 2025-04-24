import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';

// Configure the S3 client for Cloudflare R2
const S3 = new S3Client({
  region: 'auto',
  endpoint: `https://${env.CLOUDFLARE_ACCOUNT_ID}.r2.cloudflarestorage.com`,
  credentials: {
    accessKeyId: env.CLOUDFLARE_ACCESS_KEY_ID,
    secretAccessKey: env.CLOUDFLARE_SECRET_ACCESS_KEY,
  },
});

export async function POST({ request }) {
  try {
    const formData = await request.formData();
    const file = formData.get('file');
    const filename = formData.get('filename');

    if (!file || !(file instanceof File)) {
      return json({ error: 'No file provided' }, { status: 400 });
    }

    if (!filename || typeof filename !== 'string') {
      return json({ error: 'No filename provided' }, { status: 400 });
    }

    // Sanitize the filename to ensure it's safe for S3 storage
    const sanitizedFilename = filename.replace(/\s+/g, '-').replace(/[^a-zA-Z0-9\-\/\.\_]/g, '');

    // Read the file as an ArrayBuffer
    const arrayBuffer = await file.arrayBuffer();
    const buffer = Buffer.from(arrayBuffer);

    // Determine content type based on file extension
    const contentType = file.type || 'application/octet-stream';

    // Upload the file to R2
    const command = new PutObjectCommand({
      Bucket: env.CLOUDFLARE_R2_BUCKET,
      Key: sanitizedFilename,
      Body: buffer,
      ContentType: contentType,
      // Make the file publicly accessible (if your bucket is configured for public access)
      ACL: 'public-read',
    });

    await S3.send(command);

    // Return the URL to the uploaded file
    const url = `${env.CLOUDFLARE_PUBLIC_URL}/${sanitizedFilename}`;
    
    return json({
      url,
      filename: sanitizedFilename,
      success: true
    });
  } catch (error) {
    console.error('Error uploading to R2:', error);
    return json({ 
      error: error instanceof Error ? error.message : 'Unknown error',
      success: false 
    }, { status: 500 });
  }
} 
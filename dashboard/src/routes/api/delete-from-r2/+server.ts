import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { S3Client, DeleteObjectCommand } from '@aws-sdk/client-s3';

// Configure the S3 client for Cloudflare R2
const S3 = new S3Client({
  region: 'auto',
  endpoint: `https://${env.CLOUDFLARE_ACCOUNT_ID}.r2.cloudflarestorage.com`,
  credentials: {
    accessKeyId: env.CLOUDFLARE_ACCESS_KEY_ID,
    secretAccessKey: env.CLOUDFLARE_SECRET_ACCESS_KEY,
  },
});

export async function DELETE({ request }) {
  try {
    const body = await request.json();
    const { filename } = body;

    if (!filename) {
      return json({ error: 'No filename provided' }, { status: 400 });
    }

    // Delete the file from R2
    const command = new DeleteObjectCommand({
      Bucket: env.CLOUDFLARE_R2_BUCKET,
      Key: filename,
    });

    await S3.send(command);

    return json({
      success: true,
      message: `File ${filename} deleted successfully`
    });
  } catch (error) {
    console.error('Error deleting from R2:', error);
    return json({ 
      success: false, 
      error: error instanceof Error ? error.message : 'Unknown error'
    }, { status: 500 });
  }
} 
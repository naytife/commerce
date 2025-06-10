# Naytife Dashboard - Admin Interface

SvelteKit-based admin dashboard for managing the Naytife Commerce Platform.

## Features

- 📊 Store analytics and reporting
- 🛍️ Product catalog management
- 📦 Order processing and fulfillment
- 👥 Customer management
- 💳 Payment and billing oversight
- ⚙️ Platform configuration
- 🔐 OAuth2 authentication

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Open in browser
npm run dev -- --open
```

## Building

```bash
# Create production build
npm run build

# Preview production build
npm run preview
```

## Testing

```bash
# Run unit tests
npm run test

# Run e2e tests with Playwright
npm run test:e2e
```

## Integration

The dashboard integrates with:
- **Backend API**: REST endpoints for data management
- **Authentication**: OAuth2 via Ory Hydra
- **File Storage**: Cloudflare R2 for asset management

For complete platform setup, see the [main deployment guide](../DEPLOYMENT_GUIDE.md).
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://kit.svelte.dev/docs/adapters) for your target environment.
# ashia-dashboard

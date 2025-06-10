# Naytife Template 1 - Customer Storefront

A SvelteKit-based e-commerce storefront template for the Naytife Commerce Platform.

## Features

- 🛍️ Product catalog with search and filtering
- 🛒 Shopping cart functionality  
- 💳 Stripe payment integration
- 👤 Customer authentication via OAuth2
- 📱 Responsive design with Tailwind CSS
- ⚡ Fast loading with SvelteKit optimizations

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

## Integration with Naytife Platform

This template integrates with the Naytife backend services:

- **Authentication**: OAuth2 via Ory Hydra
- **API**: REST endpoints for products, orders, users
- **Payments**: Stripe payment processing
- **Assets**: Cloudflare R2 for image storage

For complete platform setup, see the [main deployment guide](../../DEPLOYMENT_GUIDE.md).

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.

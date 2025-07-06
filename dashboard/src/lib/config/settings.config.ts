export const settingsStructure = [
  {
    id: 'store',
    label: 'Store',
    children: [
      { id: 'general', label: 'Store Info', path: '/[shop]/settings/store/general', description: 'Set your store name, contact info, and basic preferences' },
      { id: 'branding', label: 'Branding', path: '/[shop]/settings/store/branding', description: 'Upload logo, cover image, and customize theme colors' },
      { id: 'domain', label: 'Domain Settings', path: '/[shop]/settings/store/domain', description: 'Connect a custom domain or use yourstore.platform.com' },
      { id: 'seo', label: 'SEO & Visibility', path: '/[shop]/settings/store/seo', description: 'Add meta tags, sitemap config, and control search engine indexing' },
    ]
  },
  {
    id: 'commerce',
    label: 'Commerce',
    children: [
      { id: 'checkout', label: 'Checkout Settings', path: '/[shop]/settings/commerce/checkout', description: 'Configure cart behavior, guest checkout, terms and conditions' },
      { id: 'payment-methods', label: 'Payment Methods', path: '/[shop]/settings/commerce/payment-methods', description: 'Enable Stripe, Paystack, bank transfer, or custom gateways' },
      { id: 'shipping', label: 'Shipping', path: '/[shop]/settings/commerce/shipping', description: 'Define where and how you deliver products' },
      { id: 'tax', label: 'Tax Settings', path: '/[shop]/settings/commerce/tax', description: 'Set VAT, GST, or regional tax rules' },
    ]
  },
  {
    id: 'integrations',
    label: 'Integrations',
    children: [
      { id: 'social-media', label: 'Connect Socials', path: '/[shop]/settings/integrations/social-media', description: 'Link Instagram, Facebook, TikTok for product sync or sharing' },
      { id: 'analytics', label: 'Analytics Tools', path: '/[shop]/settings/integrations/analytics', description: 'Add GA4, Meta Pixel, Hotjar, etc.' },
      { id: 'custom-apps', label: 'Developer Settings', path: '/[shop]/settings/integrations/custom-apps', description: 'API keys, webhook configuration, custom apps' },
    ]
  },
  {
    id: 'billing',
    label: 'Billing',
    children: [
      { id: 'invoices', label: 'Invoice History', path: '/[shop]/settings/billing/invoices', description: 'View and download past payment records' },
      { id: 'plan', label: 'Subscription Plan', path: '/[shop]/settings/billing/plan', description: 'Manage or upgrade your subscription' },
      { id: 'usage', label: 'Usage Overview', path: '/[shop]/settings/billing/usage', description: 'Track API, media, or other usage limits' },
    ]
  }
]; 
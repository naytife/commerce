<script lang="ts">
  import { Store, BarChart3, CreditCard, Globe, Shield, Zap } from 'lucide-svelte';
  import { CheckCircle, ArrowRight } from 'lucide-svelte';
  // Features array is now local
  const features = [
    {
      icon: Store,
      title: "Multi-Vendor Platform",
      description: "Enable multiple vendors to sell on your platform with automated commission management",
      benefits: ["Vendor onboarding", "Commission tracking", "Payout automation", "Vendor analytics"]
    },
    {
      icon: BarChart3,
      title: "Advanced Analytics",
      description: "Real-time insights into sales, customers, and market trends with actionable recommendations",
      benefits: ["Revenue tracking", "Customer insights", "Market trends", "Performance optimization"]
    },
    {
      icon: CreditCard,
      title: "Global Payments",
      description: "Accept payments worldwide with 100+ payment methods and instant settlements",
      benefits: ["Multiple payment gateways", "Currency conversion", "Fraud protection", "Instant payouts"]
    },
    {
      icon: Shield,
      title: "Enterprise Security",
      description: "Bank-level security with PCI compliance and advanced fraud protection",
      benefits: ["SSL encryption", "PCI compliance", "2FA authentication", "Regular security audits"]
    },
    {
      icon: Zap,
      title: "Lightning Performance",
      description: "Optimized for speed with CDN, caching, and modern web technologies",
      benefits: ["Global CDN", "Smart caching", "Image optimization", "99.9% uptime"]
    },
    {
      icon: Globe,
      title: "Global Reach",
      description: "Multi-language, multi-currency support with local payment methods",
      benefits: ["50+ languages", "150+ currencies", "Local payments", "Regional hosting"]
    }
  ];
</script>

<section class="relative z-10 container mx-auto px-6 py-20">
  <div class="text-center mb-16">
    <h2 class="text-3xl md:text-4xl font-bold mb-4 text-foreground">
      Why Nigerian Entrepreneurs Choose Naytife
    </h2>
    <p class="text-lg text-muted-foreground max-w-3xl mx-auto">
      Built specifically for the Nigerian market with features that drive real results for local businesses
    </p>
  </div>
  <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
    {#each features as feature}
      <div class="card-interactive group">
        <div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
          <svelte:component this={feature.icon} class="w-6 h-6 text-white" />
        </div>
        <h3 class="text-xl font-semibold mb-3 text-foreground">{feature.title}</h3>
        <p class="text-muted-foreground mb-4 leading-relaxed">
          {feature.description}
        </p>
        <ul class="space-y-2 mb-4">
          {#each feature.benefits as benefit}
            <li class="flex items-center gap-2 text-sm text-muted-foreground">
              <CheckCircle class="w-3 h-3 text-success flex-shrink-0" />
              <span>{benefit}</span>
            </li>
          {/each}
        </ul>
        <div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
          <span>Learn more</span>
          <ArrowRight class="w-4 h-4" />
        </div>
      </div>
    {/each}
  </div>
</section> 
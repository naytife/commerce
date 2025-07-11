<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Badge } from '$lib/components/ui/badge';
  import { Clock, Package, CheckCircle, Eye, Crown, Infinity, ArrowRight, Mail, Rocket, Shield, Phone } from 'lucide-svelte';
  let activePricingPeriod: 'monthly' | 'yearly' = 'monthly';
  function togglePricingPeriod() {
    activePricingPeriod = activePricingPeriod === 'monthly' ? 'yearly' : 'monthly';
  }
</script>

<section id="pricing" class="relative z-10 py-20 bg-linear-to-br from-surface via-background to-surface-elevated">
  <div class="container mx-auto px-6">
    <div class="text-center mb-16">
      <h2 class="text-3xl md:text-5xl font-bold mb-4 text-foreground">
        Choose Your Growth Plan
      </h2>
      <p class="text-lg text-muted-foreground max-w-3xl mx-auto mb-8">
        Start free forever, scale when ready. No hidden fees, no setup costs, no contracts.
      </p>
      <div class="flex items-center justify-center gap-4 mb-12">
        <span class="text-sm font-medium" class:text-muted-foreground={activePricingPeriod !== 'monthly'}>Monthly</span>
        <button 
          class="relative w-16 h-8 bg-muted rounded-full transition-colors"
          class:bg-primary={activePricingPeriod === 'yearly'}
          on:click={togglePricingPeriod}
        >
          <div 
            class="absolute w-6 h-6 bg-white rounded-full top-1 transition-transform shadow-md"
            class:translate-x-1={activePricingPeriod === 'monthly'}
            class:translate-x-9={activePricingPeriod === 'yearly'}
          ></div>
        </button>
        <span class="text-sm font-medium" class:text-muted-foreground={activePricingPeriod !== 'yearly'}>Yearly</span>
        <Badge variant="secondary" class="bg-success/10 text-success border-success/20 ml-2">Save 20%</Badge>
      </div>
      <div class="mb-8">
        <div class="inline-flex items-center gap-2 bg-linear-to-r from-primary/10 to-accent/10 border border-primary/20 rounded-full px-6 py-2">
          <Clock class="w-4 h-4 text-primary" />
          <span class="text-sm font-medium text-foreground">
            Limited Time: Get Premium for 3 months free with yearly plan
          </span>
        </div>
      </div>
    </div>
    <div class="grid gap-8 lg:grid-cols-3 max-w-7xl mx-auto">
      <!-- Free Plan -->
      <div class="card-elevated relative group hover:scale-105 transition-all duration-300">
        <div class="absolute inset-0 bg-linear-to-br from-muted/5 to-transparent rounded-2xl"></div>
        <div class="relative p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 bg-linear-to-br from-muted to-border rounded-2xl flex items-center justify-center">
              <Package class="w-6 h-6 text-muted-foreground" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-foreground">Starter</h3>
              <p class="text-sm text-muted-foreground">Perfect for testing ideas</p>
            </div>
          </div>
          <div class="mb-8">
            <div class="flex items-baseline gap-2 mb-2">
              <span class="text-4xl font-bold text-foreground">Free</span>
              <span class="text-muted-foreground">forever</span>
            </div>
            <p class="text-sm text-success font-medium">No credit card required</p>
          </div>
          <ul class="space-y-4 mb-8">
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-muted-foreground">Single product category</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-muted-foreground">1 product image per item</span>
            </li>
            <li class="flex items-start gap-3">
              <Eye class="w-5 h-5 text-primary mt-0.5 shrink-0" />
              <span class="text-sm text-foreground font-medium">Stay active with 50+ visitors every 14 days</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-muted-foreground">Basic analytics dashboard</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-muted-foreground">Mobile-optimized storefront</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-muted-foreground">SSL security included</span>
            </li>
          </ul>
          <Button href="/login" class="w-full glass border-border/50 hover:border-primary/50 transition-colors">
            Start Free Forever
          </Button>
        </div>
      </div>
      <!-- Premium Plan -->
      <div class="card-elevated relative group hover:scale-105 transition-all duration-300 border-2 border-primary/20">
        <div class="absolute -top-4 left-1/2 transform -translate-x-1/2">
          <Badge class="bg-primary text-primary-foreground px-6 py-2 text-sm font-semibold shadow-brand">
            Most Popular
          </Badge>
        </div>
        <div class="absolute inset-0 bg-linear-to-br from-primary/5 to-accent/5 rounded-2xl"></div>
        <div class="relative p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 bg-linear-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
              <Crown class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-foreground">Premium</h3>
              <p class="text-sm text-muted-foreground">For growing businesses</p>
            </div>
          </div>
          <div class="mb-8">
            <div class="flex items-baseline gap-2 mb-2">
              <span class="text-4xl font-bold text-foreground">
                
6{activePricingPeriod === 'monthly' ? '20,000' : '16,000'}
              </span>
              <span class="text-muted-foreground">
                /{activePricingPeriod === 'monthly' ? 'month' : 'month'}
              </span>
            </div>
            {#if activePricingPeriod === 'yearly'}
              <p class="text-sm text-success font-medium">
648,000 saved annually</p>
            {:else}
              <p class="text-sm text-muted-foreground">Billed monthly</p>
            {/if}
          </div>
          <ul class="space-y-4 mb-8">
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground font-medium">Customer accounts & profiles</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Up to 10 images per product</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">5 product categories</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Product collections & bundles</span>
            </li>
            <li class="flex items-start gap-3">
              <Infinity class="w-5 h-5 text-primary mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Unlimited monthly visitors</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Advanced analytics & reports</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Email marketing integration</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Priority customer support</span>
            </li>
          </ul>
          <Button href="/login" class="w-full btn-gradient shadow-brand">
            Upgrade to Premium
            <ArrowRight class="w-4 h-4 ml-2" />
          </Button>
        </div>
      </div>
      <!-- Enterprise Plan -->
      <div class="card-elevated relative group hover:scale-105 transition-all duration-300">
        <div class="absolute inset-0 bg-linear-to-br from-secondary/5 to-transparent rounded-2xl"></div>
        <div class="relative p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 bg-linear-to-br from-secondary to-muted rounded-2xl flex items-center justify-center">
              <Rocket class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-foreground">Enterprise</h3>
              <p class="text-sm text-muted-foreground">For large organizations</p>
            </div>
          </div>
          <div class="mb-8">
            <div class="flex items-baseline gap-2 mb-2">
              <span class="text-4xl font-bold text-foreground">Custom</span>
            </div>
            <p class="text-sm text-muted-foreground">Tailored to your needs</p>
          </div>
          <ul class="space-y-4 mb-8">
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Multi-vendor marketplace</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Unlimited everything</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Custom integrations</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Dedicated account manager</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">24/7 priority support</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">Custom domain & branding</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">API access & webhooks</span>
            </li>
            <li class="flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-success mt-0.5 shrink-0" />
              <span class="text-sm text-foreground">SLA guarantee</span>
            </li>
          </ul>
          <Button variant="outline" class="w-full glass border-border/50 hover:border-secondary/50 transition-colors group">
            <Mail class="w-4 h-4 mr-2" />
            Contact Sales
            <ArrowRight class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
          </Button>
        </div>
      </div>
    </div>
    <div class="text-center mt-16">
      <p class="text-muted-foreground mb-6">
        All plans include SSL security, mobile optimization, and basic SEO tools
      </p>
      <div class="flex justify-center gap-8 text-sm text-muted-foreground mb-8">
        <div class="flex items-center gap-2">
          <Shield class="w-4 h-4 text-success" />
          <span>Bank-level Security</span>
        </div>
        <div class="flex items-center gap-2">
          <CheckCircle class="w-4 h-4 text-success" />
          <span>99.9% Uptime SLA</span>
        </div>
        <div class="flex items-center gap-2">
          <Phone class="w-4 h-4 text-success" />
          <span>Nigerian Support Team</span>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-primary/20 bg-primary/5 max-w-2xl mx-auto">
        <div class="flex items-center gap-3 mb-4">
          <Shield class="w-4 h-4 text-success" />
          <span>PCI DSS Compliant</span>
        </div>
        <div class="flex items-center gap-3 mb-4">
          <CheckCircle class="w-4 h-4 text-success" />
          <span>24/7 Support</span>
        </div>
        <div class="flex items-center gap-3 mb-4">
          <Phone class="w-4 h-4 text-success" />
          <span>Call us: +234 800 NAYTIFE</span>
        </div>
      </div>
    </div>
  </div>
</section> 
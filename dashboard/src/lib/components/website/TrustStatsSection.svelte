<script lang="ts">
  import { Shield, CheckCircle, Globe } from 'lucide-svelte';
</script>

<section class="relative z-10 py-20">
  <div class="container mx-auto px-6">
    <div class="glass rounded-3xl p-12 border border-border/50">
      <div class="text-center mb-12">
        <h3 class="text-2xl font-bold text-foreground mb-4">
          Trusted by 10,000+ Nigerian Businesses
        </h3>
        <p class="text-muted-foreground">
          Join the fastest-growing e-commerce platform in Nigeria
        </p>
      </div>
      <div class="grid gap-8 md:grid-cols-4 text-center">
        <div>
          <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-2">
            10,000+
          </div>
          <div class="text-muted-foreground font-medium">Active Stores</div>
          <div class="text-xs text-success mt-1">+150% this year</div>
        </div>
        <div>
          <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-accent to-secondary bg-clip-text text-transparent mb-2">
            ₦50M+
          </div>
          <div class="text-muted-foreground font-medium">Monthly Revenue</div>
          <div class="text-xs text-success mt-1">Processed securely</div>
        </div>
        <div>
          <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-secondary to-primary bg-clip-text text-transparent mb-2">
            99.9%
          </div>
          <div class="text-muted-foreground font-medium">Uptime SLA</div>
          <div class="text-xs text-success mt-1">Guaranteed</div>
        </div>
        <div>
          <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent mb-2">
            4.9/5
          </div>
          <div class="text-muted-foreground font-medium">Customer Rating</div>
          <div class="text-xs text-success mt-1">5,000+ reviews</div>
        </div>
      </div>
      <div class="flex justify-center gap-8 mt-12 pt-8 border-t border-border/50">
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <Shield class="w-5 h-5 text-success" />
          <span>SSL Secured</span>
        </div>
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <CheckCircle class="w-5 h-5 text-success" />
          <span>PCI Compliant</span>
        </div>
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <Globe class="w-5 h-5 text-success" />
          <span>Nigeria Hosted</span>
        </div>
      </div>
    </div>
  </div>
</section> 
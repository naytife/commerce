<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Store, Mail, Phone, CheckCircle } from 'lucide-svelte';
</script>

<footer class="relative z-10 border-t border-border/50 glass">
  <div class="container mx-auto px-6 py-16">
    <div class="grid gap-8 md:grid-cols-4">
      <div class="md:col-span-2">
        <div class="flex items-center gap-3 mb-6">
          <div class="w-10 h-10 bg-linear-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
            <Store class="w-5 h-5 text-white" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-foreground">Naytife Commerce</h3>
            <p class="text-sm text-muted-foreground">Empowering Nigerian businesses</p>
          </div>
        </div>
        <p class="text-muted-foreground mb-6 max-w-md">
          The leading e-commerce platform designed specifically for Nigerian entrepreneurs. 
          Build, grow, and scale your online business with confidence.
        </p>
        <div class="flex gap-4">
          <Button variant="outline" size="sm" class="glass">
            <Mail class="w-4 h-4 mr-2" />
            support@naytife.com
          </Button>
          <Button variant="outline" size="sm" class="glass">
            <Phone class="w-4 h-4 mr-2" />
            +234 800 NAYTIFE
          </Button>
        </div>
      </div>
      <div>
        <h4 class="font-semibold text-foreground mb-4">Platform</h4>
        <ul class="space-y-3 text-sm text-muted-foreground">
          <li><a href="#pricing" class="hover:text-primary transition-colors">Pricing</a></li>
          <li><a href="/login" class="hover:text-primary transition-colors">Get Started</a></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Features</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Templates</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">API Documentation</button></li>
        </ul>
      </div>
      <div>
        <h4 class="font-semibold text-foreground mb-4">Support</h4>
        <ul class="space-y-3 text-sm text-muted-foreground">
          <li><button type="button" class="hover:text-primary transition-colors text-left">Help Center</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Contact Sales</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Live Chat</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Community</button></li>
          <li><button type="button" class="hover:text-primary transition-colors text-left">Status Page</button></li>
        </ul>
      </div>
    </div>
    <div class="border-t border-border/50 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
      <p class="text-sm text-muted-foreground">
        &copy; 2025 Naytife Commerce. Proudly building the future of Nigerian e-commerce.
      </p>
      <div class="flex gap-6 text-sm text-muted-foreground">
        <button type="button" class="hover:text-primary transition-colors">Privacy Policy</button>
        <button type="button" class="hover:text-primary transition-colors">Terms of Service</button>
        <button type="button" class="hover:text-primary transition-colors">Cookie Policy</button>
      </div>
    </div>
  </div>
</footer> 
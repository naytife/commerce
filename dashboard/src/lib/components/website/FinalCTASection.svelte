<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Users, Shield, ArrowRight, CheckCircle } from 'lucide-svelte';
  export let visitorsThisWeek: number;
  export let revenueToday: number;
</script>

<section class="relative z-10 py-20">
  <div class="container mx-auto px-6 text-center">
    <div class="max-w-4xl mx-auto">
      <div class="flex flex-col items-center gap-4 mb-8">
        <div class="inline-flex items-center gap-2 bg-gradient-to-r from-primary/10 to-accent/10 border border-primary/20 rounded-full px-6 py-3">
          <Users class="w-4 h-4 text-primary" />
          <span class="text-sm font-medium text-foreground">
            Join 150+ businesses who started this week
          </span>
        </div>
        <div class="flex items-center gap-6 text-xs text-muted-foreground">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-success rounded-full animate-pulse"></div>
            <span>{visitorsThisWeek} stores created in the last hour</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-primary rounded-full animate-pulse delay-500"></div>
            <span>₦{(revenueToday).toLocaleString()}+ processed today</span>
          </div>
        </div>
      </div>
      <h2 class="text-3xl md:text-5xl font-bold mb-6 text-foreground">
        Your Success Story Starts Today
      </h2>
      <p class="text-lg text-muted-foreground mb-8 max-w-3xl mx-auto leading-relaxed">
        Don't let another day pass watching competitors succeed online. Join thousands of Nigerian entrepreneurs 
        who've already transformed their businesses with Naytife Commerce - completely free to start.
      </p>
      <div class="glass rounded-2xl p-6 border border-success/20 bg-success/5 max-w-2xl mx-auto mb-10">
        <div class="flex items-center gap-3 mb-4">
          <Shield class="w-6 h-6 text-success" />
          <h4 class="font-semibold text-foreground">100% Risk-Free Promise</h4>
        </div>
        <p class="text-sm text-muted-foreground leading-relaxed">
          Start completely free - no payment required. If you're not completely satisfied with Premium, 
          cancel anytime within the first 30 days for a full refund. Your success is our guarantee.
        </p>
      </div>
      <div class="flex flex-col items-center gap-6">
        <Button 
          href="/login" 
          size="lg" 
          class="btn-gradient text-xl px-16 py-6 rounded-2xl shadow-brand group relative overflow-hidden"
        >
          <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
          <span class="relative flex items-center">
            Start Your FREE Store Now
            <ArrowRight class="w-6 h-6 ml-3 group-hover:translate-x-1 transition-transform" />
          </span>
        </Button>
        <div class="text-center space-y-2">
          <div class="flex items-center justify-center gap-6 text-sm text-muted-foreground">
            <div class="flex items-center gap-2">
              <CheckCircle class="w-4 h-4 text-success" />
              <span>Free forever plan</span>
            </div>
            <div class="flex items-center gap-2">
              <CheckCircle class="w-4 h-4 text-success" />
              <span>Setup in 5 minutes</span>
            </div>
            <div class="flex items-center gap-2">
              <CheckCircle class="w-4 h-4 text-success" />
              <span>No credit card needed</span>
            </div>
          </div>
          <p class="text-xs text-muted-foreground">
            Stay active with 50+ visitors every 14 days • Upgrade anytime • No hidden fees
          </p>
        </div>
      </div>
      <div class="mt-16 p-8 glass rounded-2xl border border-border/50">
        <p class="text-sm text-muted-foreground mb-4">
          Trusted by businesses across Nigeria
        </p>
        <div class="flex justify-center items-center gap-8 text-xs text-muted-foreground">
          <span>Lagos • Abuja • Port Harcourt • Kano • Ibadan • Enugu</span>
        </div>
      </div>
    </div>
  </div>
</section> 
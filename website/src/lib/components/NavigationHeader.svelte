<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Store, Sun, Moon, ArrowRight, Phone } from 'lucide-svelte';
  export let mode: string;
  export let toggleMode: () => void;
</script>
<header class="relative z-50 glass border-b border-border/50">
  <div class="container mx-auto px-6 py-4">
    <div class="flex items-center justify-between">
      <!-- Logo -->
      <div class="flex items-center gap-3">
        <div class="relative">
          <div class="w-10 h-10 bg-linear-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
            <Store class="w-5 h-5 text-white" />
          </div>
          <div class="absolute -top-1 -right-1 w-4 h-4 bg-linear-to-br from-accent to-secondary rounded-full animate-pulse"></div>
        </div>
        <div>
          <h1 class="text-xl font-bold bg-linear-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
            Naytife Commerce
          </h1>
          <p class="text-xs text-muted-foreground">Multi-vendor Platform</p>
        </div>
      </div>
      <!-- Navigation Actions -->
      <div class="flex items-center gap-4">
        <!-- Theme Toggle -->
        <Button 
          variant="outline" 
          size="icon" 
          on:click={toggleMode}
          class="glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-300 shadow-glass rounded-xl"
        >
          <div class="relative overflow-hidden">
            {#if mode === 'dark'}
              <Sun class="h-4 w-4 transition-all duration-300 rotate-0 scale-100 text-foreground" />
            {:else}
              <Moon class="h-4 w-4 transition-all duration-300 rotate-0 scale-100 text-foreground" />
            {/if}
          </div>
          <span class="sr-only">Toggle Theme</span>
        </Button>
        <Button variant="ghost" size="sm" class="hidden sm:inline-flex" on:click={() => document.getElementById('pricing')?.scrollIntoView({behavior: 'smooth'})}>
          Pricing
        </Button>
        <Button variant="ghost" size="sm" class="hidden sm:inline-flex">
          <Phone class="w-4 h-4 mr-2" />
          (234) 800-NAYTIFE
        </Button>
        <Button href="/login" size="sm" class="btn-gradient rounded-xl shadow-brand">
          Get Started FREE
          <ArrowRight class="w-4 h-4 ml-2" />
        </Button>
      </div>
    </div>
  </div>
</header> 
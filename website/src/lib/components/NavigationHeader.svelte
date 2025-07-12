<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Sun, Moon, ArrowRight, Phone } from 'lucide-svelte';
    import type { SystemModeValue } from 'mode-watcher';
  export let mode: SystemModeValue;
  export let toggleMode: () => void;
</script>
<header class="relative z-50 glass border-b border-border/50">
  <div class="container mx-auto px-6 py-4">
    <div class="flex items-center justify-between">
      <!-- Logo -->
      <a href="/" class="flex items-center gap-3 group transition-all duration-300 hover:scale-105">
        <div class="relative">
          {#if mode === 'dark'}
            <img src="/images/naytifelogo-dark.svg" alt="Naytife Logo" class="h-8 w-auto object-contain filter drop-shadow-sm" />
          {:else}
            <img src="/images/naytifelogo.svg" alt="Naytife Logo" class="h-8 w-auto object-contain filter drop-shadow-sm" />
          {/if}
        </div>
        <div class="hidden lg:block">
          <p class="text-xs text-muted-foreground/80 font-medium tracking-wide">Multi-vendor Platform</p>
        </div>
      </a>
      <!-- Navigation Actions -->
      <div class="flex items-center gap-4">
        <!-- Theme Toggle -->
        <Button 
          variant="outline" 
          size="icon" 
          onclick={toggleMode}
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
        <Button variant="ghost" size="sm" class="hidden sm:inline-flex" onclick={() => document.getElementById('pricing')?.scrollIntoView({behavior: 'smooth'})}>
          Pricing
        </Button>
        <Button variant="ghost" size="sm" class="hidden sm:inline-flex">
          <Phone class="w-4 h-4 mr-2" />
          (234) 800-NAYTIFE
        </Button>
        <Button href="https://dashboard.naytife.com/register" size="sm" class="btn-gradient rounded-xl shadow-brand">
          Get Started FREE
          <ArrowRight class="w-4 h-4 ml-2" />
        </Button>
      </div>
    </div>
  </div>
</header> 
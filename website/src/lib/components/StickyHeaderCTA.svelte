<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Sun, Moon, ArrowRight } from 'lucide-svelte';
    import type { SystemModeValue } from 'mode-watcher';
  export let showStickyHeader: boolean;
  export let mode: SystemModeValue;
  export let toggleMode: () => void;
</script>

{#if showStickyHeader}
  <div class="fixed top-0 left-0 right-0 z-[100] glass border-b border-border/50 backdrop-blur-xl animate-slide-in">
    <div class="container mx-auto px-6 py-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="relative">
            {#if mode === 'dark'}
              <img src="/images/naytifelogo-dark.svg" alt="Naytife Logo" class="h-6 w-auto object-contain filter drop-shadow-sm" />
            {:else}
              <img src="/images/naytifelogo.svg" alt="Naytife Logo" class="h-6 w-auto object-contain filter drop-shadow-sm" />
            {/if}
          </div>
          <div class="hidden sm:block">
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <div class="w-1.5 h-1.5 bg-success rounded-full animate-pulse"></div>
              <span>Free forever • No credit card required</span>
            </div>
          </div>
        </div>
        <div class="flex items-center gap-4">
          <Button 
            variant="outline" 
            size="icon" 
            onclick={toggleMode}
            class="h-8 w-8 glass border-border/50 hover:bg-surface-elevated transition-all"
          >
            {#if mode === 'dark'}
              <Sun class="h-3 w-3 text-foreground" />
            {:else}
              <Moon class="h-3 w-3 text-foreground" />
            {/if}
          </Button>
          <Button href="https://dashboard.naytife.com/register" size="sm" class="btn-gradient rounded-xl shadow-brand text-sm">
            Start FREE Now
            <ArrowRight class="w-3 h-3 ml-2" />
          </Button>
        </div>
      </div>
    </div>
  </div>
{/if} 
<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Rocket, ArrowRight, CheckCircle } from 'lucide-svelte';
  export let showExitIntent: boolean;
  export let close: () => void;
</script>

{#if showExitIntent}
  <div class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in">
    <div class="glass rounded-3xl p-8 border border-primary/20 max-w-lg mx-6 animate-scale-in shadow-2xl">
      <button 
        class="absolute top-4 right-4 w-8 h-8 rounded-full bg-background/50 hover:bg-background flex items-center justify-center transition-colors"
        on:click={close}
      >
        <span class="text-muted-foreground text-lg">&times;</span>
      </button>
      <div class="text-center">
        <div class="w-16 h-16 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-brand">
          <Rocket class="w-8 h-8 text-white" />
        </div>
        <h3 class="text-2xl font-bold text-foreground mb-4">
          Wait! Before You Go...
        </h3>
        <p class="text-muted-foreground mb-6 leading-relaxed">
          Join 150+ Nigerian entrepreneurs who started their stores this week. 
          <strong class="text-primary">Completely free</strong> - no credit card needed!
        </p>
        <div class="space-y-4">
          <Button href="/login" class="w-full btn-gradient text-lg py-4 rounded-xl shadow-brand group">
            <span class="flex items-center justify-center">
              Start My FREE Store Now
              <ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
            </span>
          </Button>
          <div class="flex items-center justify-center gap-4 text-xs text-muted-foreground">
            <div class="flex items-center gap-1">
              <CheckCircle class="w-3 h-3 text-success" />
              <span>Free forever</span>
            </div>
            <div class="flex items-center gap-1">
              <CheckCircle class="w-3 h-3 text-success" />
              <span>5-min setup</span>
            </div>
            <div class="flex items-center gap-1">
              <CheckCircle class="w-3 h-3 text-success" />
              <span>No contracts</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if} 
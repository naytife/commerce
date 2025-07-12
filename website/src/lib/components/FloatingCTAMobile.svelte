<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { ArrowRight } from 'lucide-svelte';
</script>

<div class="fixed bottom-6 left-6 right-6 z-50 md:hidden">
  <Button 
    href="https://dashboard.naytife.com/register" 
    class="w-full btn-gradient text-lg py-4 rounded-2xl shadow-xl group"
  >
    Start FREE Forever
    <ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
  </Button>
</div> 
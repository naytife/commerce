<script lang="ts">
  export let message: string = 'Coming Soon';
</script>

<div class="relative w-full h-full">
  <slot />
  <div class="absolute inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm select-none pointer-events-auto cursor-not-allowed">
    <div class="flex flex-col items-center">
      <span class="text-3xl md:text-4xl font-extrabold text-white drop-shadow-lg animate-pulse tracking-wide">
        {message}
      </span>
      <span class="mt-2 text-base text-white/80 font-medium">This feature is not available yet.</span>
    </div>
  </div>
</div>

<style>
  /* Ensures overlay blocks all pointer events below it */
  .pointer-events-auto {
    pointer-events: auto;
  }
</style> 
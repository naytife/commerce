<script lang="ts">
  import { Shield, CheckCircle, Globe } from 'lucide-svelte';
  import { onMount } from 'svelte';

  let currentSlide = 0;
  let carouselContainer: HTMLDivElement;
  
  const stats = [
    {
      value: '10,000+',
      label: 'Active Stores',
      subtitle: '+150% this year',
      gradient: 'from-primary to-accent'
    },
    {
      value: '₦50M+',
      label: 'Monthly Revenue',
      subtitle: 'Processed securely',
      gradient: 'from-accent to-secondary'
    },
    {
      value: '99.9%',
      label: 'Uptime SLA',
      subtitle: 'Guaranteed',
      gradient: 'from-secondary to-primary'
    },
    {
      value: '4.9/5',
      label: 'Customer Rating',
      subtitle: '5,000+ reviews',
      gradient: 'from-primary to-secondary'
    }
  ];

  onMount(() => {
    const interval = setInterval(() => {
      currentSlide = (currentSlide + 1) % stats.length;
      if (carouselContainer) {
        const slideWidth = carouselContainer.scrollWidth / stats.length;
        carouselContainer.scrollTo({
          left: slideWidth * currentSlide,
          behavior: 'smooth'
        });
      }
    }, 3000); // Change slide every 3 seconds

    return () => clearInterval(interval);
  });
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
      <div class="overflow-x-auto md:overflow-x-visible scrollbar-hide" bind:this={carouselContainer}>
        <div class="flex gap-8 md:grid md:grid-cols-4 text-center min-w-max md:min-w-0">
          {#each stats as stat}
            <div class="flex-shrink-0 md:flex-shrink px-4 md:px-0 w-64 md:w-auto">
              <div class="text-4xl md:text-5xl font-bold bg-gradient-to-r {stat.gradient} bg-clip-text text-transparent mb-2">
                {stat.value}
              </div>
              <div class="text-muted-foreground font-medium">{stat.label}</div>
              <div class="text-xs text-success mt-1">{stat.subtitle}</div>
            </div>
          {/each}
        </div>
      </div>
      
      <!-- Mobile carousel indicators -->
      <div class="flex justify-center gap-2 mt-6 md:hidden">
        {#each stats as _, index}
          <button
            class="w-2 h-2 rounded-full transition-all duration-300 {index === currentSlide ? 'bg-primary' : 'bg-muted-foreground/30'}"
            aria-label="Go to slide {index + 1}"
            on:click={() => {
              currentSlide = index;
              if (carouselContainer) {
                const slideWidth = carouselContainer.scrollWidth / stats.length;
                carouselContainer.scrollTo({
                  left: slideWidth * index,
                  behavior: 'smooth'
                });
              }
            }}
          ></button>
        {/each}
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
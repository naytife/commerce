<script lang="ts">
  import { Star, CheckCircle, Users } from 'lucide-svelte';
  import { onMount } from 'svelte';
  // Testimonials array is local
  const testimonials = [
    {
      name: "Adebayo Ogundimu",
      role: "Founder, Lagos Fashion Hub",
      content: "Naytife transformed our local fashion store into a thriving multi-vendor marketplace. We've grown from 5 to 150+ vendors in just 6 months!",
      image: "/images/testimonial-1.jpg",
      rating: 5,
      revenue: "₦2.5M+ monthly"
    },
    {
      name: "Fatima Al-Hassan",
      role: "CEO, Abuja Craft Market",
      content: "The analytics insights helped us increase our conversion rate by 340%. The vendor management tools are absolutely game-changing.",
      image: "/images/testimonial-2.jpg",
      rating: 5,
      revenue: "₦1.8M+ monthly"
    },
    {
      name: "Chidi Okafor",
      role: "Director, Port Harcourt Tech Stores",
      content: "Moving from WooCommerce to Naytife was the best decision we made. The performance improvement alone increased our sales by 60%.",
      image: "/images/testimonial-3.jpg",
      rating: 5,
      revenue: "₦4.2M+ monthly"
    }
  ];
  let activeTestimonial = 0;
  let interval: any;
  onMount(() => {
    interval = setInterval(() => {
      activeTestimonial = (activeTestimonial + 1) % testimonials.length;
    }, 5000);
    return () => clearInterval(interval);
  });
  function setActiveTestimonial(idx: number) {
    activeTestimonial = idx;
  }
</script>

<section class="relative z-10 py-20 bg-linear-to-r from-primary/5 via-accent/5 to-secondary/5">
  <div class="container mx-auto px-6">
    <div class="text-center mb-16">
      <h2 class="text-3xl md:text-4xl font-bold mb-4 text-foreground">
        Real Results from Real Businesses
      </h2>
      <p class="text-lg text-muted-foreground max-w-2xl mx-auto">
        See how Nigerian entrepreneurs are building successful online businesses with Naytife
      </p>
    </div>
    <div class="grid gap-8 lg:grid-cols-3">
      {#each testimonials as testimonial, index}
        <div 
          class="card-elevated group transition-all duration-500"
          class:ring-2={activeTestimonial === index}
          class:ring-primary={activeTestimonial === index}
        >
          <div class="p-8">
            <div class="flex gap-1 mb-4">
              {#each Array(testimonial.rating) as _}
                <Star class="w-4 h-4 fill-yellow-400 text-yellow-400" />
              {/each}
            </div>
            <blockquote class="text-muted-foreground mb-6 italic leading-relaxed">
              "{testimonial.content}"
            </blockquote>
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 bg-linear-to-br from-primary to-accent rounded-full flex items-center justify-center">
                <span class="text-white font-bold text-lg">
                  {testimonial.name.charAt(0)}
                </span>
              </div>
              <div>
                <div class="font-semibold text-foreground">{testimonial.name}</div>
                <div class="text-sm text-muted-foreground">{testimonial.role}</div>
                <div class="text-sm font-medium text-success">{testimonial.revenue}</div>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
    <div class="flex justify-center gap-2 mt-8">
      {#each testimonials as _, index}
        <button 
          class="w-3 h-3 rounded-full transition-all"
          class:bg-primary={activeTestimonial === index}
          class:bg-muted={activeTestimonial !== index}
          on:click={() => setActiveTestimonial(index)}
        ></button>
      {/each}
    </div>
  </div>
</section> 
<script>
	import { onMount } from 'svelte';
	let isVisible = false;

	onMount(() => {
		isVisible = true;
	});
</script>

<section class="relative min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted overflow-hidden flex flex-col items-center justify-center">
	<!-- Animated SVG blobs background -->
	<svg class="absolute -top-32 -left-32 w-[480px] h-[480px] z-0 animate-pulse opacity-60" viewBox="0 0 600 600" fill="none" xmlns="http://www.w3.org/2000/svg">
		<path d="M300,520Q220,540,140,480Q60,420,80,320Q100,220,180,140Q260,60,360,100Q460,140,500,220Q540,300,480,380Q420,460,300,520Z" fill="#e0e7ff" class="dark:fill-primary/10"/>
	</svg>
	<svg class="absolute top-0 right-0 w-[320px] h-[320px] z-0 opacity-40" viewBox="0 0 600 600" fill="none" xmlns="http://www.w3.org/2000/svg">
		<path d="M300,520Q220,540,140,480Q60,420,80,320Q100,220,180,140Q260,60,360,100Q460,140,500,220Q540,300,480,380Q420,460,300,520Z" fill="#fef3c7" class="dark:fill-accent/10"/>
	</svg>

	<!-- Decorative SVG highlights -->
	<svg class="absolute left-1/2 top-24 -translate-x-1/2 z-10" width="540" height="40" viewBox="0 0 540 40" fill="none" xmlns="http://www.w3.org/2000/svg">
		<rect x="0" y="20" width="280" height="16" rx="8" fill="#E9D5FF" fill-opacity="0.7" class="dark:fill-primary/10"/>
		<rect x="200" y="0" width="340" height="16" rx="8" fill="#FDE68A" fill-opacity="0.7" class="dark:fill-accent/10"/>
	</svg>

	<div class="relative z-20 container mx-auto px-4 py-12 lg:py-20 flex flex-col items-center text-center gap-6 max-w-7xl transition-all duration-1000 {isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}">
		<!-- Top headline and highlights -->
		<h1 class="text-4xl md:text-6xl font-extrabold text-foreground leading-tight relative mb-4 md:mb-6 transition-all duration-700 delay-300 {isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
			<span class="relative z-20 block">
				<span class="inline-block">
					Scale <span class="inline-block relative"><span class="z-10 relative">your online store</span><span class="absolute left-0 right-0 bottom-0 h-3 bg-primary/20 rounded-full -z-10 transition-all"></span></span> and
				</span>
			</span>
			<span class="relative z-20 block">
				<span class="inline-block">
					multiply your sales with <span class="italic font-serif text-primary bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">Naytife</span>
				</span>
			</span>
		</h1>

		<!-- Two-column layout -->
		<div class="w-full flex flex-col md:flex-row items-center justify-center gap-12 md:gap-20 text-left transition-all duration-700 delay-500 {isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
			<!-- Left: Large vertical text, subheading, CTA -->
			<div class="flex-1 flex flex-col items-center md:items-start justify-center gap-6 max-w-xl">
				<div class="text-3xl md:text-5xl lg:text-6xl font-light leading-tight space-y-2">
					<div class="tracking-tight transition-all duration-500 delay-700 {isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-4'} text-foreground">Sell</div>
					<div class="tracking-tight transition-all duration-500 delay-900 {isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-4'} text-foreground">Manage</div>
					<div class="transition-all duration-500 delay-1100 {isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-4'}">
						<span class="text-primary font-semibold">Grow</span><span class="text-4xl md:text-5xl lg:text-6xl font-light text-muted-foreground animate-ping">|</span>
					</div>
				</div>
				<div class="text-lg md:text-xl font-medium text-muted-foreground mt-2 transition-all duration-500 delay-1300 {isVisible ? 'opacity-100' : 'opacity-0'}">Nigeria's Leading E-commerce Platform</div>
				<div class="flex gap-4 mt-4 transition-all duration-500 delay-1500 {isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'}">
					<a href="https://dashboard.naytife.com/register" class="inline-block px-8 py-4 text-lg font-semibold btn-gradient rounded-lg shadow-brand hover:scale-105 active:scale-95">Start Selling Free</a>
					<a href="#pricing" class="inline-block px-8 py-4 text-lg font-semibold text-primary bg-background/80 backdrop-blur border border-primary/20 rounded-lg shadow hover:bg-surface-elevated focus:ring-4 focus:ring-primary/20 transition-all hover:scale-105 active:scale-95 dark:bg-surface-elevated/80 dark:text-primary dark:border-primary/40">View Plans</a>
				</div>
			</div>
			<!-- Right: Image -->
			<div class="flex-1 flex justify-center items-center w-full">
				<div class="relative w-full max-w-2xl rounded-3xl overflow-hidden shadow-2xl border-4 border-background aspect-video bg-background/70 backdrop-blur-lg transition-all duration-700 delay-1000 {isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'} dark:border-surface-elevated dark:bg-surface-elevated/70">
					<img 
						src="https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?auto=format&fit=crop&w=1200&q=80" 
						alt="Nigerian entrepreneur managing online store with multiple payment methods and inventory tracking on computer screen" 
						class="w-full h-full object-cover object-center"
					/>
				</div>
			</div>
		</div>
	</div>
</section>



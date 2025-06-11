<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { mode, toggleMode } from 'mode-watcher';
	import { 
		ArrowRight, 
		Store, 
		BarChart3, 
		CreditCard, 
		Globe, 
		Shield, 
		Zap,
		Users,
		TrendingUp,
		CheckCircle,
		Star,
		Clock,
		Image,
		Package,
		ShoppingCart,
		Eye,
		Trash2,
		Crown,
		Infinity,
		Mail,
		Phone,
		ChevronDown,
		Play,
		Award,
		Target,
		Rocket,
		Moon,
		Sun
	} from 'lucide-svelte';

	let activeTestimonial = 0;
	let activePricingPeriod = 'monthly';
	let currentSlide = 0;
	let showStickyHeader = false;
	let showExitIntent = false;
	let visitorsThisWeek = 147; // Dynamic number for social proof
	let revenueToday = 850000; // Dynamic revenue counter

	// Show sticky header on scroll
	function handleScroll() {
		showStickyHeader = window.scrollY > 800;
	}

	// Exit intent detection
	function handleMouseLeave(e: MouseEvent) {
		if (e.clientY <= 0 && !showExitIntent) {
			showExitIntent = true;
		}
	}

	// Simulate live visitor counter
	function updateLiveStats() {
		visitorsThisWeek += Math.floor(Math.random() * 3);
		revenueToday += Math.floor(Math.random() * 50000);
	}

	const demoSlides = [
		{
			title: "Beautiful Storefronts",
			description: "Create stunning, mobile-optimized stores that convert visitors into customers",
			feature: "storefront"
		},
		{
			title: "Powerful Analytics",
			description: "Track sales, monitor performance, and make data-driven decisions",
			feature: "analytics"
		},
		{
			title: "Vendor Management",
			description: "Effortlessly manage multiple vendors and automate commission payouts",
			feature: "vendors"
		}
	];

	// Auto-rotate testimonials
	setInterval(() => {
		activeTestimonial = (activeTestimonial + 1) % testimonials.length;
	}, 5000);

	// Auto-rotate demo slides
	setInterval(() => {
		currentSlide = (currentSlide + 1) % demoSlides.length;
	}, 4000);

	// Handle scroll for sticky header
	if (typeof window !== 'undefined') {
		window.addEventListener('scroll', handleScroll);
		document.addEventListener('mouseleave', handleMouseLeave);
	}

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

	const features = [
		{
			icon: Store,
			title: "Multi-Vendor Platform",
			description: "Enable multiple vendors to sell on your platform with automated commission management",
			benefits: ["Vendor onboarding", "Commission tracking", "Payout automation", "Vendor analytics"]
		},
		{
			icon: BarChart3,
			title: "Advanced Analytics",
			description: "Real-time insights into sales, customers, and market trends with actionable recommendations",
			benefits: ["Revenue tracking", "Customer insights", "Market trends", "Performance optimization"]
		},
		{
			icon: CreditCard,
			title: "Global Payments",
			description: "Accept payments worldwide with 100+ payment methods and instant settlements",
			benefits: ["Multiple payment gateways", "Currency conversion", "Fraud protection", "Instant payouts"]
		},
		{
			icon: Shield,
			title: "Enterprise Security",
			description: "Bank-level security with PCI compliance and advanced fraud protection",
			benefits: ["SSL encryption", "PCI compliance", "2FA authentication", "Regular security audits"]
		},
		{
			icon: Zap,
			title: "Lightning Performance",
			description: "Optimized for speed with CDN, caching, and modern web technologies",
			benefits: ["Global CDN", "Smart caching", "Image optimization", "99.9% uptime"]
		},
		{
			icon: Globe,
			title: "Global Reach",
			description: "Multi-language, multi-currency support with local payment methods",
			benefits: ["50+ languages", "150+ currencies", "Local payments", "Regional hosting"]
		}
	];

	// Auto-rotate testimonials
	setInterval(() => {
		activeTestimonial = (activeTestimonial + 1) % testimonials.length;
	}, 5000);
</script>

<div class="min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
	<!-- Sticky Header CTA -->
	{#if showStickyHeader}
		<div class="fixed top-0 left-0 right-0 z-[100] glass border-b border-border/50 backdrop-blur-xl animate-slide-in">
			<div class="container mx-auto px-6 py-3">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-primary to-accent rounded-xl flex items-center justify-center">
							<Store class="w-4 h-4 text-white" />
						</div>
						<div>
							<span class="font-semibold text-foreground text-sm">Naytife Commerce</span>
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
							on:click={toggleMode}
							class="h-8 w-8 glass border-border/50 hover:bg-surface-elevated transition-all"
						>
							{#if $mode === 'dark'}
								<Sun class="h-3 w-3 text-foreground" />
							{:else}
								<Moon class="h-3 w-3 text-foreground" />
							{/if}
						</Button>
						
						<Button href="/login" size="sm" class="btn-gradient rounded-xl shadow-brand text-sm">
							Start FREE Now
							<ArrowRight class="w-3 h-3 ml-2" />
						</Button>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Ambient background effects -->
	<div class="absolute inset-0 overflow-hidden">
		<!-- Floating orbs -->
		<div class="absolute top-1/4 left-1/4 w-96 h-96 bg-gradient-to-r from-primary/20 to-accent/20 rounded-full blur-3xl animate-pulse"></div>
		<div class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-gradient-to-r from-accent/20 to-secondary/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
		<div class="absolute top-1/2 right-1/3 w-64 h-64 bg-gradient-to-r from-secondary/20 to-primary/20 rounded-full blur-3xl animate-pulse delay-2000"></div>
	</div>

	<!-- Grid pattern overlay -->
	<div class="absolute inset-0 opacity-[0.02] dark:opacity-[0.05]" style="background-image: radial-gradient(circle at 1px 1px, currentColor 1px, transparent 0); background-size: 40px 40px;"></div>

	<!-- Navigation Header -->
	<header class="relative z-50 glass border-b border-border/50">
		<div class="container mx-auto px-6 py-4">
			<div class="flex items-center justify-between">
				<!-- Logo -->
				<div class="flex items-center gap-3">
					<div class="relative">
						<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
							<Store class="w-5 h-5 text-white" />
						</div>
						<div class="absolute -top-1 -right-1 w-4 h-4 bg-gradient-to-br from-accent to-secondary rounded-full animate-pulse"></div>
					</div>
					<div>
						<h1 class="text-xl font-bold bg-gradient-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
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
							{#if $mode === 'dark'}
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

	<!-- Hero Section -->
	<section class="relative z-10 container mx-auto px-6 pt-16 pb-12 text-center">
		<!-- Trust Indicators -->
		<div class="flex justify-center mb-8 animate-fade-in">
			<div class="flex flex-col sm:flex-row items-center gap-4">
				<Badge variant="secondary" class="glass px-4 py-2 text-sm font-medium border-border/50">
					<Award class="w-4 h-4 mr-2 text-success" />
					#1 Nigerian E-commerce Platform
				</Badge>
				<Badge variant="secondary" class="glass px-3 py-2 text-xs font-medium border-border/50">
					<Target class="w-3 h-3 mr-2 text-primary" />
					94% Customer Success Rate
				</Badge>
			</div>
		</div>

		<!-- Main Heading -->
		<h1 class="text-4xl md:text-6xl lg:text-7xl font-bold mb-6 animate-fade-in delay-100">
			<span class="block text-foreground">Launch Your</span>
			<span class="block bg-gradient-to-r from-primary via-accent to-secondary bg-clip-text text-transparent">
				Online Store in Minutes
			</span>
		</h1>

		<!-- Benefit-focused Subtitle -->
		<p class="text-lg md:text-xl text-muted-foreground max-w-4xl mx-auto mb-8 leading-relaxed animate-fade-in delay-200">
			Join 10,000+ Nigerian entrepreneurs making <strong class="text-success">₦50M+ monthly</strong> with our 
			zero-setup-fee platform. Start selling today with our <strong class="text-primary">freemium plan</strong> - 
			no technical skills required.
		</p>

		<!-- Urgency & Value Proposition -->
		<div class="flex justify-center mb-10 animate-fade-in delay-300">
			<div class="glass rounded-2xl p-4 border border-primary/20 bg-primary/5">
				<div class="flex items-center gap-3 text-sm">
					<Rocket class="w-5 h-5 text-primary" />
					<span class="text-foreground font-medium">
						<strong class="text-primary">Always Free:</strong> Get your store live in under 10 minutes + Free domain setup
					</span>
				</div>
			</div>
		</div>

		<!-- Primary CTA with Social Proof -->
		<div class="flex flex-col items-center gap-6 mb-12 animate-fade-in delay-400">
			<Button href="/login" size="lg" class="btn-gradient text-xl px-12 py-5 rounded-2xl shadow-brand min-w-[280px] group relative overflow-hidden">
				<div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
				<span class="relative flex items-center">
					Start FREE Forever
					<ArrowRight class="w-6 h-6 ml-3 group-hover:translate-x-1 transition-transform" />
				</span>
			</Button>
			
			<div class="flex items-center gap-4 text-sm text-muted-foreground">
				<div class="flex items-center gap-2">
					<CheckCircle class="w-4 h-4 text-success" />
					<span>No Credit Card Required</span>
				</div>
				<div class="flex items-center gap-2">
					<CheckCircle class="w-4 h-4 text-success" />
					<span>Setup in 5 Minutes</span>
				</div>
				<div class="flex items-center gap-2">
					<CheckCircle class="w-4 h-4 text-success" />
					<span>Always Free Plan</span>
				</div>
			</div>
		</div>

		<!-- Customer Testimonial Preview -->
		<div class="max-w-2xl mx-auto animate-fade-in delay-500">
			<div class="glass rounded-2xl p-6 border border-border/50">
				<div class="flex items-center gap-4 mb-4">
					<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center">
						<span class="text-white font-bold text-lg">A</span>
					</div>
					<div class="text-left">
						<div class="font-semibold text-foreground">Adebayo Ogundimu</div>
						<div class="text-sm text-muted-foreground">Lagos Fashion Hub</div>
					</div>
					<div class="ml-auto flex gap-1">
						{#each Array(5) as _}
							<Star class="w-4 h-4 fill-yellow-400 text-yellow-400" />
						{/each}
					</div>
				</div>
				<p class="text-muted-foreground italic">
					"Started completely free, grew to 150+ vendors in 6 months. Now making ₦2.5M monthly - best decision ever!"
				</p>
			</div>
		</div>
	</section>

	<!-- Hero Visual Section -->
	<section class="relative z-10 py-20 overflow-hidden">
		<div class="container mx-auto px-6 max-w-none">
			<!-- Platform Preview -->
			<div class="max-w-7xl mx-auto relative">
				<div class="text-center mb-16">
					<h3 class="text-3xl md:text-4xl font-bold text-foreground mb-6">
						See Why 10,000+ Stores Choose Naytife
					</h3>
					<p class="text-lg text-muted-foreground max-w-2xl mx-auto mb-8">
						Watch our platform in action - from setup to sale in minutes, not months
					</p>
					
					<!-- Slide Indicators -->
					<div class="flex justify-center gap-3 mb-8">
						{#each demoSlides as _, i}
							<button 
								class="w-3 h-3 rounded-full transition-all duration-300"
								class:bg-primary={i === currentSlide}
								class:bg-muted={i !== currentSlide}
								class:scale-125={i === currentSlide}
								on:click={() => currentSlide = i}
							></button>
						{/each}
					</div>
				</div>

				<!-- Interactive Demo Slideshow -->
				<div class="relative overflow-hidden">
					<!-- Slide Container -->
					<div 
						class="flex transition-transform duration-700 ease-in-out"
						style="transform: translateX(-{currentSlide * 100}%)"
					>
						{#each demoSlides as slide, slideIndex}
							<div class="w-full flex-shrink-0">
								<div class="glass rounded-3xl p-8 border border-border/50 shadow-2xl bg-gradient-to-br from-background/50 to-surface-elevated/50">
									<!-- Browser Chrome -->
									<div class="flex items-center gap-2 mb-6 pb-4 border-b border-border/50">
										<div class="flex gap-2">
											<div class="w-3 h-3 bg-red-500 rounded-full animate-pulse"></div>
											<div class="w-3 h-3 bg-yellow-500 rounded-full animate-pulse delay-100"></div>
											<div class="w-3 h-3 bg-green-500 rounded-full animate-pulse delay-200"></div>
										</div>
										<div class="flex-1 mx-4">
											<div class="bg-surface-muted rounded-lg px-4 py-2 text-sm text-muted-foreground flex items-center gap-2">
												<Shield class="w-3 h-3 text-success" />
												https://{slide.feature}.naytife.com
											</div>
										</div>
										<Badge variant="secondary" class="bg-success/10 text-success border-success/20 text-xs">
											Live Demo
										</Badge>
									</div>

									<!-- Slide Content -->
									{#if slide.feature === 'storefront'}
										<!-- Storefront Demo -->
										<div class="space-y-6">
											<!-- Store Header -->
											<div class="border border-border/30 rounded-2xl p-6 bg-gradient-to-br from-background to-surface-elevated">
												<div class="flex items-center justify-between mb-6">
													<div class="flex items-center gap-4">
														<div class="w-16 h-16 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-lg">
															<Store class="w-8 h-8 text-white" />
														</div>
														<div>
															<h4 class="text-xl font-bold text-foreground">Lagos Fashion Hub</h4>
															<p class="text-sm text-muted-foreground">Premium Nigerian Fashion • 150+ Products</p>
														</div>
													</div>
													<div class="flex items-center gap-3">
														<Badge variant="secondary" class="bg-success/10 text-success border-success/20">
															<div class="w-2 h-2 bg-success rounded-full mr-2 animate-pulse"></div>
															Live Store
														</Badge>
														<Button size="sm" class="btn-gradient shadow-lg">
															<ShoppingCart class="w-4 h-4 mr-2" />
															Shop Now
														</Button>
													</div>
												</div>

												<!-- Featured Products -->
												<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
													{#each Array(4) as _, i}
														<div class="bg-surface rounded-xl p-4 border border-border/30 hover:border-primary/30 transition-colors group">
															<div class="aspect-square bg-gradient-to-br from-primary/10 to-accent/10 rounded-lg mb-3 flex items-center justify-center group-hover:scale-105 transition-transform">
																{#if i === 0}<Image class="w-8 h-8 text-primary" />{/if}
																{#if i === 1}<Package class="w-8 h-8 text-accent" />{/if}
																{#if i === 2}<ShoppingCart class="w-8 h-8 text-success" />{/if}
																{#if i === 3}<Crown class="w-8 h-8 text-secondary" />{/if}
															</div>
															<h5 class="font-medium text-foreground text-sm mb-1">Premium Item {i + 1}</h5>
															<p class="text-primary font-bold text-sm">₦{(25000 + i * 10000).toLocaleString()}</p>
															<div class="flex items-center gap-1 mt-2">
																{#each Array(5) as _}
																	<Star class="w-3 h-3 fill-yellow-400 text-yellow-400" />
																{/each}
																<span class="text-xs text-muted-foreground ml-1">(4.9)</span>
															</div>
														</div>
													{/each}
												</div>
											</div>
										</div>
									{:else if slide.feature === 'analytics'}
										<!-- Analytics Demo -->
										<div class="space-y-6">
											<div class="border border-border/30 rounded-2xl p-6 bg-gradient-to-br from-background to-surface-elevated">
												<div class="flex items-center justify-between mb-6">
													<h4 class="text-xl font-bold text-foreground">Sales Dashboard</h4>
													<div class="flex items-center gap-3">
														<Badge variant="secondary" class="bg-primary/10 text-primary border-primary/20">
															<div class="w-2 h-2 bg-primary rounded-full mr-2 animate-pulse"></div>
															Real-time Data
														</Badge>
														<Button size="sm" variant="outline">
															<BarChart3 class="w-4 h-4 mr-2" />
															View Report
														</Button>
													</div>
												</div>
												
												<!-- Revenue Metrics -->
												<div class="grid grid-cols-2 md:grid-cols-4 gap-6 mb-6">
													<div class="text-center p-4 bg-primary/5 rounded-xl border border-primary/20">
														<div class="text-3xl font-bold text-primary mb-1 animate-pulse">₦4.2M</div>
														<div class="text-sm text-muted-foreground">This Month</div>
														<div class="text-xs text-success mt-1">↗ +24.5%</div>
													</div>
													<div class="text-center p-4 bg-accent/5 rounded-xl border border-accent/20">
														<div class="text-3xl font-bold text-accent mb-1">2,847</div>
														<div class="text-sm text-muted-foreground">Orders</div>
														<div class="text-xs text-success mt-1">↗ +18.2%</div>
													</div>
													<div class="text-center p-4 bg-success/5 rounded-xl border border-success/20">
														<div class="text-3xl font-bold text-success mb-1">89.4%</div>
														<div class="text-sm text-muted-foreground">Conversion</div>
														<div class="text-xs text-success mt-1">↗ +5.1%</div>
													</div>
													<div class="text-center p-4 bg-secondary/5 rounded-xl border border-secondary/20">
														<div class="text-3xl font-bold text-secondary mb-1">15.2K</div>
														<div class="text-sm text-muted-foreground">Visitors</div>
														<div class="text-xs text-success mt-1">↗ +12.8%</div>
													</div>
												</div>

												<!-- Chart Visualization -->
												<div class="bg-surface/50 rounded-xl p-6 border border-border/30">
													<div class="flex items-center justify-between mb-4">
														<h5 class="font-semibold text-foreground">Revenue Trend</h5>
														<Badge variant="secondary" class="text-xs">Last 30 Days</Badge>
													</div>
													<div class="h-32 bg-gradient-to-t from-primary/20 to-transparent rounded-lg relative overflow-hidden">
														<div class="absolute inset-0 flex items-end justify-around pb-2">
															{#each Array(7) as _, i}
																<div 
																	class="bg-primary rounded-t w-8 transition-all duration-1000 delay-{i * 100}"
																	style="height: {20 + Math.random() * 60}%"
																></div>
															{/each}
														</div>
													</div>
												</div>
											</div>
										</div>
									{:else if slide.feature === 'vendors'}
										<!-- Vendor Management Demo -->
										<div class="space-y-6">
											<div class="border border-border/30 rounded-2xl p-6 bg-gradient-to-br from-background to-surface-elevated">
												<div class="flex items-center justify-between mb-6">
													<h4 class="text-xl font-bold text-foreground">Vendor Management</h4>
													<div class="flex items-center gap-3">
														<Badge variant="secondary" class="bg-accent/10 text-accent border-accent/20">
															<Users class="w-3 h-3 mr-2" />
															150+ Vendors
														</Badge>
														<Button size="sm" class="btn-gradient">
															<Users class="w-4 h-4 mr-2" />
															Add Vendor
														</Button>
													</div>
												</div>

												<!-- Vendor Grid -->
												<div class="space-y-4">
													{#each Array(3) as _, i}
														<div class="flex items-center justify-between p-4 bg-surface rounded-xl border border-border/30 hover:border-primary/30 transition-colors">
															<div class="flex items-center gap-4">
																<div class="w-12 h-12 bg-gradient-to-br from-accent to-secondary rounded-xl flex items-center justify-center">
																	<Store class="w-6 h-6 text-white" />
																</div>
																<div>
																	<h6 class="font-semibold text-foreground">
																		{#if i === 0}Fashion Hub Lagos{/if}
																		{#if i === 1}Tech Store Abuja{/if}
																		{#if i === 2}Craft Market PH{/if}
																	</h6>
																	<p class="text-sm text-muted-foreground">
																		{#if i === 0}₦850K revenue • 12.5% commission{/if}
																		{#if i === 1}₦420K revenue • 10% commission{/if}
																		{#if i === 2}₦380K revenue • 15% commission{/if}
																	</p>
																</div>
															</div>
															<div class="flex items-center gap-3">
																<Badge variant="secondary" class="bg-success/10 text-success border-success/20 text-xs">
																	Active
																</Badge>
																<div class="text-right">
																	<div class="text-sm font-bold text-foreground">
																		{#if i === 0}₦106K{/if}
																		{#if i === 1}₦42K{/if}
																		{#if i === 2}₦57K{/if}
																	</div>
																	<div class="text-xs text-muted-foreground">Commission</div>
																</div>
															</div>
														</div>
													{/each}
												</div>

												<!-- Commission Summary -->
												<div class="mt-6 p-4 bg-primary/5 rounded-xl border border-primary/20">
													<div class="flex items-center justify-between">
														<div>
															<h6 class="font-semibold text-foreground">Total Commission This Month</h6>
															<p class="text-sm text-muted-foreground">From 150+ active vendors</p>
														</div>
														<div class="text-right">
															<div class="text-2xl font-bold text-primary">₦205K</div>
															<div class="text-xs text-success">↗ +28.4%</div>
														</div>
													</div>
												</div>
											</div>
										</div>
									{/if}

									<!-- Slide Info -->
									<div class="mt-8 text-center">
										<h4 class="text-2xl font-bold text-foreground mb-2">{slide.title}</h4>
										<p class="text-muted-foreground max-w-2xl mx-auto">{slide.description}</p>
									</div>
								</div>
							</div>
						{/each}
					</div>

					<!-- Navigation Arrows -->
					<button 
						class="absolute left-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-background/90 hover:bg-background border border-border/50 rounded-full flex items-center justify-center transition-all duration-200 hover:scale-110 z-10"
						on:click={() => currentSlide = currentSlide === 0 ? demoSlides.length - 1 : currentSlide - 1}
					>
						<ArrowRight class="w-5 h-5 rotate-180 text-foreground" />
					</button>
					<button 
						class="absolute right-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-background/90 hover:bg-background border border-border/50 rounded-full flex items-center justify-center transition-all duration-200 hover:scale-110 z-10"
						on:click={() => currentSlide = (currentSlide + 1) % demoSlides.length}
					>
						<ArrowRight class="w-5 h-5 text-foreground" />
					</button>

					<!-- Floating Feature Callouts - Positioned outside container -->
				</div>

				<!-- Feature Callouts positioned relative to the outer container -->
				<div class="absolute inset-0 pointer-events-none">
					<div class="absolute left-4 top-16 hidden xl:block animate-float-up z-20 pointer-events-auto">
						<div class="glass rounded-xl p-4 border border-primary/20 bg-primary/5 w-60 shadow-lg backdrop-blur-sm">
							<div class="flex items-center gap-3 mb-2">
								<Zap class="w-5 h-5 text-primary" />
								<span class="font-semibold text-foreground text-sm">Lightning Fast</span>
							</div>
							<p class="text-xs text-muted-foreground">Sub-2 second load times with global CDN</p>
						</div>
					</div>

					<div class="absolute right-4 top-32 hidden xl:block animate-float-down z-20 pointer-events-auto">
						<div class="glass rounded-xl p-4 border border-accent/20 bg-accent/5 w-60 shadow-lg backdrop-blur-sm">
							<div class="flex items-center gap-3 mb-2">
								<Shield class="w-5 h-5 text-accent" />
								<span class="font-semibold text-foreground text-sm">Bank-Level Security</span>
							</div>
							<p class="text-xs text-muted-foreground">SSL encryption + PCI compliance included</p>
						</div>
					</div>

					<div class="absolute left-4 bottom-16 hidden xl:block animate-float-up delay-1000 z-20 pointer-events-auto">
						<div class="glass rounded-xl p-4 border border-success/20 bg-success/5 w-60 shadow-lg backdrop-blur-sm">
							<div class="flex items-center gap-3 mb-2">
								<Users class="w-5 h-5 text-success" />
								<span class="font-semibold text-foreground text-sm">Multi-Vendor Ready</span>
							</div>
							<p class="text-xs text-muted-foreground">Scale to 1000+ vendors seamlessly</p>
						</div>
					</div>
				</div>

				<!-- Demo CTAs -->
				<div class="text-center mt-16 space-y-4">
					<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
						<Button size="lg" class="btn-gradient shadow-lg group">
							<Play class="w-5 h-5 mr-2" />
							Watch Full Demo (2 min)
							<ArrowRight class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
						</Button>
						<Button variant="outline" size="lg" class="glass border-border/50 hover:border-primary/50 transition-colors">
							Schedule Live Demo
						</Button>
					</div>
					<p class="text-sm text-muted-foreground">
						See how easy it is to get started • No technical skills required
					</p>
				</div>

				<!-- Mobile Feature Highlights - Show when floating callouts are hidden -->
				<div class="block xl:hidden mt-12">
					<div class="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
						<div class="glass rounded-xl p-6 border border-primary/20 bg-primary/5 text-center">
							<div class="w-12 h-12 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-4">
								<Zap class="w-6 h-6 text-primary" />
							</div>
							<h4 class="font-semibold text-foreground mb-2">Lightning Fast</h4>
							<p class="text-sm text-muted-foreground">Sub-2 second load times with global CDN</p>
						</div>
						<div class="glass rounded-xl p-6 border border-accent/20 bg-accent/5 text-center">
							<div class="w-12 h-12 bg-accent/20 rounded-full flex items-center justify-center mx-auto mb-4">
								<Shield class="w-6 h-6 text-accent" />
							</div>
							<h4 class="font-semibold text-foreground mb-2">Bank-Level Security</h4>
							<p class="text-sm text-muted-foreground">SSL encryption + PCI compliance included</p>
						</div>
						<div class="glass rounded-xl p-6 border border-success/20 bg-success/5 text-center">
							<div class="w-12 h-12 bg-success/20 rounded-full flex items-center justify-center mx-auto mb-4">
								<Users class="w-6 h-6 text-success" />
							</div>
							<h4 class="font-semibold text-foreground mb-2">Multi-Vendor Ready</h4>
							<p class="text-sm text-muted-foreground">Scale to 1000+ vendors seamlessly</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Pricing Section -->
	<section id="pricing" class="relative z-10 py-20 bg-gradient-to-br from-surface via-background to-surface-elevated">
		<div class="container mx-auto px-6">
			<div class="text-center mb-16">
				<h2 class="text-3xl md:text-5xl font-bold mb-4 text-foreground">
					Choose Your Growth Plan
				</h2>
				<p class="text-lg text-muted-foreground max-w-3xl mx-auto mb-8">
					Start free forever, scale when ready. No hidden fees, no setup costs, no contracts.
				</p>
				
				<!-- Pricing Toggle -->
				<div class="flex items-center justify-center gap-4 mb-12">
					<span class="text-sm font-medium" class:text-muted-foreground={activePricingPeriod !== 'monthly'}>Monthly</span>
					<button 
						class="relative w-16 h-8 bg-muted rounded-full transition-colors"
						class:bg-primary={activePricingPeriod === 'yearly'}
						on:click={() => activePricingPeriod = activePricingPeriod === 'monthly' ? 'yearly' : 'monthly'}
					>
						<div 
							class="absolute w-6 h-6 bg-white rounded-full top-1 transition-transform shadow-md"
							class:translate-x-1={activePricingPeriod === 'monthly'}
							class:translate-x-9={activePricingPeriod === 'yearly'}
						></div>
					</button>
					<span class="text-sm font-medium" class:text-muted-foreground={activePricingPeriod !== 'yearly'}>Yearly</span>
					<Badge variant="secondary" class="bg-success/10 text-success border-success/20 ml-2">Save 20%</Badge>
				</div>
				
				<!-- Limited Time Offer -->
				<div class="mb-8">
					<div class="inline-flex items-center gap-2 bg-gradient-to-r from-primary/10 to-accent/10 border border-primary/20 rounded-full px-6 py-2">
						<Clock class="w-4 h-4 text-primary" />
						<span class="text-sm font-medium text-foreground">
							Limited Time: Get Premium for 3 months free with yearly plan
						</span>
					</div>
				</div>
			</div>

			<div class="grid gap-8 lg:grid-cols-3 max-w-7xl mx-auto">
				<!-- Free Plan -->
				<div class="card-elevated relative group hover:scale-105 transition-all duration-300">
					<div class="absolute inset-0 bg-gradient-to-br from-muted/5 to-transparent rounded-2xl"></div>
					<div class="relative p-8">
						<div class="flex items-center gap-3 mb-6">
							<div class="w-12 h-12 bg-gradient-to-br from-muted to-border rounded-2xl flex items-center justify-center">
								<Package class="w-6 h-6 text-muted-foreground" />
							</div>
							<div>
								<h3 class="text-xl font-bold text-foreground">Starter</h3>
								<p class="text-sm text-muted-foreground">Perfect for testing ideas</p>
							</div>
						</div>

						<div class="mb-8">
							<div class="flex items-baseline gap-2 mb-2">
								<span class="text-4xl font-bold text-foreground">Free</span>
								<span class="text-muted-foreground">forever</span>
							</div>
							<p class="text-sm text-success font-medium">No credit card required</p>
						</div>

						<ul class="space-y-4 mb-8">
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-muted-foreground">Single product category</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-muted-foreground">1 product image per item</span>
							</li>
							<li class="flex items-start gap-3">
								<Eye class="w-5 h-5 text-primary mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground font-medium">Stay active with 50+ visitors every 14 days</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-muted-foreground">Basic analytics dashboard</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-muted-foreground">Mobile-optimized storefront</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-muted-foreground">SSL security included</span>
							</li>
						</ul>

						<Button href="/login" class="w-full glass border-border/50 hover:border-primary/50 transition-colors">
							Start Free Forever
						</Button>
					</div>
				</div>

				<!-- Premium Plan -->
				<div class="card-elevated relative group hover:scale-105 transition-all duration-300 border-2 border-primary/20">
					<!-- Popular Badge -->
					<div class="absolute -top-4 left-1/2 transform -translate-x-1/2">
						<Badge class="bg-primary text-primary-foreground px-6 py-2 text-sm font-semibold shadow-brand">
							Most Popular
						</Badge>
					</div>
					
					<div class="absolute inset-0 bg-gradient-to-br from-primary/5 to-accent/5 rounded-2xl"></div>
					<div class="relative p-8">
						<div class="flex items-center gap-3 mb-6">
							<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
								<Crown class="w-6 h-6 text-white" />
							</div>
							<div>
								<h3 class="text-xl font-bold text-foreground">Premium</h3>
								<p class="text-sm text-muted-foreground">For growing businesses</p>
							</div>
						</div>

						<div class="mb-8">
							<div class="flex items-baseline gap-2 mb-2">
								<span class="text-4xl font-bold text-foreground">
									₦{activePricingPeriod === 'monthly' ? '20,000' : '16,000'}
								</span>
								<span class="text-muted-foreground">
									/{activePricingPeriod === 'monthly' ? 'month' : 'month'}
								</span>
							</div>
							{#if activePricingPeriod === 'yearly'}
								<p class="text-sm text-success font-medium">₦48,000 saved annually</p>
							{:else}
								<p class="text-sm text-muted-foreground">Billed monthly</p>
							{/if}
						</div>

						<ul class="space-y-4 mb-8">
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground font-medium">Customer accounts & profiles</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Up to 10 images per product</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">5 product categories</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Product collections & bundles</span>
							</li>
							<li class="flex items-start gap-3">
								<Infinity class="w-5 h-5 text-primary mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Unlimited monthly visitors</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Advanced analytics & reports</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Email marketing integration</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Priority customer support</span>
							</li>
						</ul>

						<Button href="/login" class="w-full btn-gradient shadow-brand">
							Upgrade to Premium
							<ArrowRight class="w-4 h-4 ml-2" />
						</Button>
					</div>
				</div>

				<!-- Enterprise Plan -->
				<div class="card-elevated relative group hover:scale-105 transition-all duration-300">
					<div class="absolute inset-0 bg-gradient-to-br from-secondary/5 to-transparent rounded-2xl"></div>
					<div class="relative p-8">
						<div class="flex items-center gap-3 mb-6">
							<div class="w-12 h-12 bg-gradient-to-br from-secondary to-muted rounded-2xl flex items-center justify-center">
								<Rocket class="w-6 h-6 text-white" />
							</div>
							<div>
								<h3 class="text-xl font-bold text-foreground">Enterprise</h3>
								<p class="text-sm text-muted-foreground">For large organizations</p>
							</div>
						</div>

						<div class="mb-8">
							<div class="flex items-baseline gap-2 mb-2">
								<span class="text-4xl font-bold text-foreground">Custom</span>
							</div>
							<p class="text-sm text-muted-foreground">Tailored to your needs</p>
						</div>

						<ul class="space-y-4 mb-8">
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Multi-vendor marketplace</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Unlimited everything</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Custom integrations</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Dedicated account manager</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">24/7 priority support</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">Custom domain & branding</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">API access & webhooks</span>
							</li>
							<li class="flex items-start gap-3">
								<CheckCircle class="w-5 h-5 text-success mt-0.5 flex-shrink-0" />
								<span class="text-sm text-foreground">SLA guarantee</span>
							</li>
						</ul>

						<Button variant="outline" class="w-full glass border-border/50 hover:border-secondary/50 transition-colors group">
							<Mail class="w-4 h-4 mr-2" />
							Contact Sales
							<ArrowRight class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
						</Button>
					</div>
				</div>
			</div>

			<!-- Pricing Footer -->
			<div class="text-center mt-16">
				<p class="text-muted-foreground mb-6">
					All plans include SSL security, mobile optimization, and basic SEO tools
				</p>
				<div class="flex justify-center gap-8 text-sm text-muted-foreground mb-8">
					<div class="flex items-center gap-2">
						<Shield class="w-4 h-4 text-success" />
						<span>Bank-level Security</span>
					</div>
					<div class="flex items-center gap-2">
						<CheckCircle class="w-4 h-4 text-success" />
						<span>99.9% Uptime SLA</span>
					</div>
					<div class="flex items-center gap-2">
						<Phone class="w-4 h-4 text-success" />
						<span>Nigerian Support Team</span>
					</div>
				</div>
				
				<!-- Freemium Explanation -->
				<div class="glass rounded-2xl p-6 border border-primary/20 bg-primary/5 max-w-2xl mx-auto">
					<div class="flex items-center gap-3 mb-4">
						<Eye class="w-5 h-5 text-primary" />
						<h4 class="font-semibold text-foreground">How our Free Forever plan works</h4>
					</div>
					<p class="text-sm text-muted-foreground leading-relaxed">
						Your store stays active as long as it receives at least <strong class="text-foreground">50 unique visitors every 14 days</strong>. 
						If visitor count drops below this threshold, your store is temporarily paused until traffic picks up again. 
						No data is lost, and reactivation is automatic when you hit 50+ visitors.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section class="relative z-10 container mx-auto px-6 py-20">
		<div class="text-center mb-16">
			<h2 class="text-3xl md:text-4xl font-bold mb-4 text-foreground">
				Why Nigerian Entrepreneurs Choose Naytife
			</h2>
			<p class="text-lg text-muted-foreground max-w-3xl mx-auto">
				Built specifically for the Nigerian market with features that drive real results for local businesses
			</p>
		</div>

		<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
			{#each features as feature}
				<div class="card-interactive group">
					<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
						<svelte:component this={feature.icon} class="w-6 h-6 text-white" />
					</div>
					<h3 class="text-xl font-semibold mb-3 text-foreground">{feature.title}</h3>
					<p class="text-muted-foreground mb-4 leading-relaxed">
						{feature.description}
					</p>
					
					<!-- Feature Benefits -->
					<ul class="space-y-2 mb-4">
						{#each feature.benefits as benefit}
							<li class="flex items-center gap-2 text-sm text-muted-foreground">
								<CheckCircle class="w-3 h-3 text-success flex-shrink-0" />
								<span>{benefit}</span>
							</li>
						{/each}
					</ul>
					
					<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
						<span>Learn more</span>
						<ArrowRight class="w-4 h-4" />
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- Success Stories Section -->
	<section class="relative z-10 py-20 bg-gradient-to-r from-primary/5 via-accent/5 to-secondary/5">
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
							<!-- Rating -->
							<div class="flex gap-1 mb-4">
								{#each Array(testimonial.rating) as _}
									<Star class="w-4 h-4 fill-yellow-400 text-yellow-400" />
								{/each}
							</div>

							<!-- Testimonial Content -->
							<blockquote class="text-muted-foreground mb-6 italic leading-relaxed">
								"{testimonial.content}"
							</blockquote>

							<!-- Author Info -->
							<div class="flex items-center gap-4">
								<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center">
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

			<!-- Testimonial Navigation -->
			<div class="flex justify-center gap-2 mt-8">
				{#each testimonials as _, index}
					<button 
						class="w-3 h-3 rounded-full transition-all"
						class:bg-primary={activeTestimonial === index}
						class:bg-muted={activeTestimonial !== index}
						on:click={() => activeTestimonial = index}
					></button>
				{/each}
			</div>
		</div>
	</section>

	<!-- Trust & Stats Section -->
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
				
				<div class="grid gap-8 md:grid-cols-4 text-center">
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-2">
							10,000+
						</div>
						<div class="text-muted-foreground font-medium">Active Stores</div>
						<div class="text-xs text-success mt-1">+150% this year</div>
					</div>
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-accent to-secondary bg-clip-text text-transparent mb-2">
							₦50M+
						</div>
						<div class="text-muted-foreground font-medium">Monthly Revenue</div>
						<div class="text-xs text-success mt-1">Processed securely</div>
					</div>
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-secondary to-primary bg-clip-text text-transparent mb-2">
							99.9%
						</div>
						<div class="text-muted-foreground font-medium">Uptime SLA</div>
						<div class="text-xs text-success mt-1">Guaranteed</div>
					</div>
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent mb-2">
							4.9/5
						</div>
						<div class="text-muted-foreground font-medium">Customer Rating</div>
						<div class="text-xs text-success mt-1">5,000+ reviews</div>
					</div>
				</div>

				<!-- Trust Badges -->
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

	<!-- FAQ Section -->
	<section class="relative z-10 py-20 bg-gradient-to-br from-surface-elevated via-background to-surface">
		<div class="container mx-auto px-6">
			<div class="text-center mb-16">
				<h2 class="text-3xl md:text-4xl font-bold mb-4 text-foreground">
					Frequently Asked Questions
				</h2>
				<p class="text-lg text-muted-foreground max-w-2xl mx-auto">
					Everything you need to know about getting started with Naytife Commerce
				</p>
			</div>

			<div class="max-w-4xl mx-auto grid gap-6 md:grid-cols-2">
				<div class="card-elevated">
					<div class="p-6">
						<h4 class="font-semibold text-foreground mb-3">How quickly can I start selling?</h4>
						<p class="text-muted-foreground text-sm">
							You can have your store live and accepting payments in under 10 minutes. Our setup wizard guides you through every step.
						</p>
					</div>
				</div>
				
				<div class="card-elevated">
					<div class="p-6">
						<h4 class="font-semibold text-foreground mb-3">Is the free plan really free forever?</h4>
						<p class="text-muted-foreground text-sm">
							Yes! As long as your store gets at least 50 visitors every 14 days, it stays active forever. No hidden fees or time limits.
						</p>
					</div>
				</div>
				
				<div class="card-elevated">
					<div class="p-6">
						<h4 class="font-semibold text-foreground mb-3">What happens if I get less than 50 visitors?</h4>
						<p class="text-muted-foreground text-sm">
							Your store will be temporarily paused. Once you get 50+ visitors again within 14 days, it automatically reactivates.
						</p>
					</div>
				</div>
				
				<div class="card-elevated">
					<div class="p-6">
						<h4 class="font-semibold text-foreground mb-3">Can I upgrade or downgrade anytime?</h4>
						<p class="text-muted-foreground text-sm">
							Absolutely! Change your plan anytime with immediate effect. No contracts, no penalties, complete flexibility.
						</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Final CTA Section -->
	<section class="relative z-10 py-20">
		<div class="container mx-auto px-6 text-center">
			<div class="max-w-4xl mx-auto">
				<!-- Urgency Element -->
				<div class="flex flex-col items-center gap-4 mb-8">
					<div class="inline-flex items-center gap-2 bg-gradient-to-r from-primary/10 to-accent/10 border border-primary/20 rounded-full px-6 py-3">
						<Users class="w-4 h-4 text-primary" />
						<span class="text-sm font-medium text-foreground">
							Join 150+ businesses who started this week
						</span>
					</div>
					
					<!-- Live Activity Indicator -->
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

				<!-- Risk-Free Guarantee -->
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

				<!-- CTA with Risk Reversal -->
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

				<!-- Social Proof Footer -->
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

	<!-- Footer -->
	<footer class="relative z-10 border-t border-border/50 glass">
		<div class="container mx-auto px-6 py-16">
			<div class="grid gap-8 md:grid-cols-4">
				<!-- Company Info -->
				<div class="md:col-span-2">
					<div class="flex items-center gap-3 mb-6">
						<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
							<Store class="w-5 h-5 text-white" />
						</div>
						<div>
							<h3 class="text-xl font-bold text-foreground">Naytife Commerce</h3>
							<p class="text-sm text-muted-foreground">Empowering Nigerian businesses</p>
						</div>
					</div>
					<p class="text-muted-foreground mb-6 max-w-md">
						The leading e-commerce platform designed specifically for Nigerian entrepreneurs. 
						Build, grow, and scale your online business with confidence.
					</p>
					<div class="flex gap-4">
						<Button variant="outline" size="sm" class="glass">
							<Mail class="w-4 h-4 mr-2" />
							support@naytife.com
						</Button>
						<Button variant="outline" size="sm" class="glass">
							<Phone class="w-4 h-4 mr-2" />
							+234 800 NAYTIFE
						</Button>
					</div>
				</div>

				<!-- Quick Links -->
				<div>
					<h4 class="font-semibold text-foreground mb-4">Platform</h4>
					<ul class="space-y-3 text-sm text-muted-foreground">
						<li><a href="#pricing" class="hover:text-primary transition-colors">Pricing</a></li>
						<li><a href="/login" class="hover:text-primary transition-colors">Get Started</a></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Features</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Templates</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">API Documentation</button></li>
					</ul>
				</div>

				<!-- Support -->
				<div>
					<h4 class="font-semibold text-foreground mb-4">Support</h4>
					<ul class="space-y-3 text-sm text-muted-foreground">
						<li><button type="button" class="hover:text-primary transition-colors text-left">Help Center</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Contact Sales</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Live Chat</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Community</button></li>
						<li><button type="button" class="hover:text-primary transition-colors text-left">Status Page</button></li>
					</ul>
				</div>
			</div>

			<!-- Footer Bottom -->
			<div class="border-t border-border/50 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
				<p class="text-sm text-muted-foreground">
					&copy; 2025 Naytife Commerce. Proudly building the future of Nigerian e-commerce.
				</p>
				<div class="flex gap-6 text-sm text-muted-foreground">
					<button type="button" class="hover:text-primary transition-colors">Privacy Policy</button>
					<button type="button" class="hover:text-primary transition-colors">Terms of Service</button>
					<button type="button" class="hover:text-primary transition-colors">Cookie Policy</button>
				</div>
			</div>
		</div>
	</footer>

	<!-- Floating CTA for Mobile -->
	<div class="fixed bottom-6 left-6 right-6 z-50 md:hidden">
		<Button 
			href="/login" 
			class="w-full btn-gradient text-lg py-4 rounded-2xl shadow-xl group"
		>
			Start FREE Forever
			<ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
		</Button>
	</div>

	<!-- Exit Intent Modal -->
	{#if showExitIntent}
		<div class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in">
			<div class="glass rounded-3xl p-8 border border-primary/20 max-w-lg mx-6 animate-scale-in shadow-2xl">
				<button 
					class="absolute top-4 right-4 w-8 h-8 rounded-full bg-background/50 hover:bg-background flex items-center justify-center transition-colors"
					on:click={() => showExitIntent = false}
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
</div>

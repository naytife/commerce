<script>
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
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
		CheckCircle
	} from 'lucide-svelte';
</script>

<div class="min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
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
					<Button variant="ghost" size="sm" class="hidden sm:inline-flex">
						Features
					</Button>
					<Button variant="ghost" size="sm" class="hidden sm:inline-flex">
						Pricing
					</Button>
					<Button href="/login" size="sm" class="btn-gradient rounded-xl shadow-brand">
						Get Started
						<ArrowRight class="w-4 h-4 ml-2" />
					</Button>
				</div>
			</div>
		</div>
	</header>

	<!-- Hero Section -->
	<section class="relative z-10 container mx-auto px-6 pt-20 pb-16 text-center">
		<!-- Status Badge -->
		<div class="flex justify-center mb-8 animate-fade-in">
			<Badge variant="secondary" class="glass px-4 py-2 text-sm font-medium border-border/50">
				<div class="w-2 h-2 bg-success rounded-full mr-2 animate-pulse"></div>
				Now supporting multi-vendor stores
			</Badge>
		</div>

		<!-- Main Heading -->
		<h1 class="text-4xl md:text-6xl lg:text-7xl font-bold mb-6 animate-fade-in delay-200">
			<span class="block text-foreground">Scale Your</span>
			<span class="block bg-gradient-to-r from-primary via-accent to-secondary bg-clip-text text-transparent">
				E-commerce Empire
			</span>
		</h1>

		<!-- Subtitle -->
		<p class="text-lg md:text-xl text-muted-foreground max-w-3xl mx-auto mb-10 leading-relaxed animate-fade-in delay-300">
			The next-generation platform for multi-vendor marketplaces. Build, scale, and monetize 
			your online business with advanced analytics, seamless payments, and powerful automation.
		</p>

		<!-- CTA Buttons -->
		<div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16 animate-fade-in delay-500">
			<Button href="/login" size="lg" class="btn-gradient text-lg px-8 py-4 rounded-2xl shadow-brand min-w-[200px] group">
				Start Building
				<ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
			</Button>
			<Button variant="outline" size="lg" class="glass text-lg px-8 py-4 rounded-2xl border-border/50 min-w-[200px]">
				Watch Demo
			</Button>
		</div>

		<!-- Social Proof -->
		<div class="flex flex-col sm:flex-row items-center justify-center gap-8 text-sm text-muted-foreground animate-fade-in delay-700">
			<div class="flex items-center gap-2">
				<CheckCircle class="w-4 h-4 text-success" />
				<span>Free 14-day trial</span>
			</div>
			<div class="flex items-center gap-2">
				<CheckCircle class="w-4 h-4 text-success" />
				<span>No credit card required</span>
			</div>
			<div class="flex items-center gap-2">
				<CheckCircle class="w-4 h-4 text-success" />
				<span>Cancel anytime</span>
			</div>
		</div>
	</section>

	<!-- Features Grid -->
	<section class="relative z-10 container mx-auto px-6 py-20">
		<div class="text-center mb-16">
			<h2 class="text-3xl md:text-4xl font-bold mb-4 text-foreground">
				Everything you need to succeed
			</h2>
			<p class="text-lg text-muted-foreground max-w-2xl mx-auto">
				Powerful features designed to help entrepreneurs and businesses thrive in the digital marketplace
			</p>
		</div>

		<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
			<!-- Feature 1: Multi-vendor -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<Store class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Multi-Vendor Platform</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Enable multiple vendors to sell on your platform. Manage commissions, payouts, and vendor relationships effortlessly.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>

			<!-- Feature 2: Analytics -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-accent to-secondary rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<BarChart3 class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Advanced Analytics</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Get deep insights into sales performance, customer behavior, and market trends with real-time dashboards.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>

			<!-- Feature 3: Payments -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-secondary to-primary rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<CreditCard class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Global Payments</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Accept payments worldwide with support for 100+ payment methods and automatic currency conversion.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>

			<!-- Feature 4: Global Reach -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-primary to-secondary rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<Globe class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Global Infrastructure</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Fast, reliable hosting with CDN, multi-region deployment, and 99.9% uptime guarantee for your store.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>

			<!-- Feature 5: Security -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-accent to-primary rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<Shield class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Enterprise Security</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Bank-level security with SSL encryption, PCI compliance, and advanced fraud protection for your business.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>

			<!-- Feature 6: Performance -->
			<div class="card-interactive group">
				<div class="w-12 h-12 bg-gradient-to-br from-secondary to-accent rounded-2xl flex items-center justify-center mb-6 shadow-brand group-hover:scale-110 transition-transform">
					<Zap class="w-6 h-6 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-3 text-foreground">Lightning Fast</h3>
				<p class="text-muted-foreground mb-4 leading-relaxed">
					Optimized for speed with advanced caching, image optimization, and modern web technologies.
				</p>
				<div class="flex items-center text-primary font-medium group-hover:gap-3 gap-2 transition-all">
					<span>Learn more</span>
					<ArrowRight class="w-4 h-4" />
				</div>
			</div>
		</div>
	</section>

	<!-- Stats Section -->
	<section class="relative z-10 py-20">
		<div class="container mx-auto px-6">
			<div class="glass rounded-3xl p-12 border border-border/50">
				<div class="grid gap-8 md:grid-cols-3 text-center">
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-2">
							10k+
						</div>
						<div class="text-muted-foreground font-medium">Active Stores</div>
					</div>
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-accent to-secondary bg-clip-text text-transparent mb-2">
							$50M+
						</div>
						<div class="text-muted-foreground font-medium">Revenue Processed</div>
					</div>
					<div>
						<div class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-secondary to-primary bg-clip-text text-transparent mb-2">
							99.9%
						</div>
						<div class="text-muted-foreground font-medium">Uptime Guarantee</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="relative z-10 py-20">
		<div class="container mx-auto px-6 text-center">
			<div class="max-w-4xl mx-auto">
				<h2 class="text-3xl md:text-4xl font-bold mb-6 text-foreground">
					Ready to transform your business?
				</h2>
				<p class="text-lg text-muted-foreground mb-10 max-w-2xl mx-auto leading-relaxed">
					Join thousands of successful entrepreneurs who've grown their businesses with Naytife Commerce. 
					Start your free trial today and experience the difference.
				</p>
				<Button 
					href="/login" 
					size="lg" 
					class="btn-gradient text-lg px-12 py-4 rounded-2xl shadow-brand group"
				>
					Start Your Free Trial
					<ArrowRight class="w-5 h-5 ml-3 group-hover:translate-x-1 transition-transform" />
				</Button>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="relative z-10 border-t border-border/50 glass">
		<div class="container mx-auto px-6 py-12">
			<div class="text-center text-muted-foreground">
				<p>&copy; 2025 Naytife Commerce. Building the future of e-commerce.</p>
			</div>
		</div>
	</footer>
</div>

<script lang="ts">
	import { mode, toggleMode } from 'mode-watcher';
	import Hero from '$lib/components/website/Hero.svelte';
	import StickyHeaderCTA from '$lib/components/website/StickyHeaderCTA.svelte';
	import AmbientBackground from '$lib/components/website/AmbientBackground.svelte';
	import NavigationHeader from '$lib/components/website/NavigationHeader.svelte';
	import PricingSection from '$lib/components/website/PricingSection.svelte';
	import FeaturesSection from '$lib/components/website/FeaturesSection.svelte';
	import TestimonialsSection from '$lib/components/website/TestimonialsSection.svelte';
	import TrustStatsSection from '$lib/components/website/TrustStatsSection.svelte';
	import FAQSection from '$lib/components/website/FAQSection.svelte';
	import FinalCTASection from '$lib/components/website/FinalCTASection.svelte';
	import FooterSection from '$lib/components/website/FooterSection.svelte';
	import FloatingCTAMobile from '$lib/components/website/FloatingCTAMobile.svelte';
	import ExitIntentModal from '$lib/components/website/ExitIntentModal.svelte';
	import { Store, BarChart3, CreditCard, Globe, Shield, Zap } from 'lucide-svelte';

	let showStickyHeader = false;
	let showExitIntent = false;
	let visitorsThisWeek = 147;
	let revenueToday = 850000;

	function handleScroll() {
		showStickyHeader = window.scrollY > 800;
	}
	function handleMouseLeave(e: MouseEvent) {
		if (e.clientY <= 0 && !showExitIntent) {
			showExitIntent = true;
		}
	}
	function updateLiveStats() {
		visitorsThisWeek += Math.floor(Math.random() * 3);
		revenueToday += Math.floor(Math.random() * 50000);
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

	if (typeof window !== 'undefined') {
		window.addEventListener('scroll', handleScroll);
		document.addEventListener('mouseleave', handleMouseLeave);
	}

	function closeExitIntent() {
		showExitIntent = false;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
	<StickyHeaderCTA {showStickyHeader} mode={$mode ?? 'light'} {toggleMode} />
	<AmbientBackground />
	<NavigationHeader mode={$mode ?? 'light'} {toggleMode} />
	<Hero />
	<PricingSection />
	<FeaturesSection />
	<TestimonialsSection />
	<TrustStatsSection />
	<FAQSection />
	<FinalCTASection {visitorsThisWeek} {revenueToday} />
	<FooterSection />
	<FloatingCTAMobile />
	<ExitIntentModal {showExitIntent} close={closeExitIntent} />
</div>

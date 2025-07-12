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

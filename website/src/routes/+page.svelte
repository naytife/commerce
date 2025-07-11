<svelte:head>
	<title>Naytife - Nigeria's #1 E-commerce Platform for Online Vendors</title>
	<meta name="description" content="Build your dream online store in minutes with Naytife. Trusted by 10,000+ Nigerian merchants. Multi-vendor platform, advanced analytics, secure payments, and 24/7 support. Start your free trial today!" />
	<meta name="keywords" content="e-commerce platform Nigeria, online store builder, multi-vendor marketplace, Nigerian e-commerce, online business platform, sell online Nigeria" />
	<meta property="og:title" content="Naytife - Nigeria's #1 E-commerce Platform for Online Vendors" />
	<meta property="og:description" content="Join 10,000+ successful Nigerian merchants who've built thriving online businesses with Naytife. Start your free trial today!" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://naytife.com" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="Naytife - Nigeria's #1 E-commerce Platform" />
	<meta name="twitter:description" content="Build your dream online store in minutes. Trusted by 10,000+ Nigerian merchants." />
</svelte:head>

<script lang="ts">
  import { mode, toggleMode } from '$lib/mode-watcher';
  import Hero from '$lib/components/Hero.svelte';
  import StickyHeaderCTA from '$lib/components/StickyHeaderCTA.svelte';
  import AmbientBackground from '$lib/components/AmbientBackground.svelte';
  import NavigationHeader from '$lib/components/NavigationHeader.svelte';
  import PricingSection from '$lib/components/PricingSection.svelte';
  import FeaturesSection from '$lib/components/FeaturesSection.svelte';
  import TestimonialsSection from '$lib/components/TestimonialsSection.svelte';
  import TrustStatsSection from '$lib/components/TrustStatsSection.svelte';
  import FAQSection from '$lib/components/FAQSection.svelte';
  import FinalCTASection from '$lib/components/FinalCTASection.svelte';
  import FooterSection from '$lib/components/FooterSection.svelte';
  import FloatingCTAMobile from '$lib/components/FloatingCTAMobile.svelte';
  import ExitIntentModal from '$lib/components/ExitIntentModal.svelte';

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

<div class="min-h-screen bg-linear-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
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
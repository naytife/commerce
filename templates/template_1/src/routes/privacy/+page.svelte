<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Shield, Lock, Eye } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Privacy Policy</title>
	<meta name="description" content="Learn how we collect, use, and protect your personal information" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-8">
		<h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">Privacy Policy</h1>
		<p class="text-lg text-gray-600 dark:text-gray-400">
			Last updated: {new Date().toLocaleDateString()}
		</p>
	</div>

	<div class="space-y-8">
		<!-- Introduction -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Shield class="h-5 w-5" />
					Our Commitment to Your Privacy
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					We are committed to protecting your privacy and ensuring you have a positive experience on our website. 
					This Privacy Policy outlines how we collect, use, and safeguard your information when you visit our 
					website or make a purchase from us.
				</p>
			</Card.Content>
		</Card.Root>

		<!-- Information We Collect -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Eye class="h-5 w-5" />
					Information We Collect
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Personal Information</h4>
					<p class="text-gray-600 dark:text-gray-400">
						When you make a purchase or create an account, we collect personal information such as your name, 
						email address, phone number, shipping address, and payment information.
					</p>
				</div>
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Usage Information</h4>
					<p class="text-gray-600 dark:text-gray-400">
						We automatically collect information about how you use our website, including your IP address, 
						browser type, pages visited, and the time spent on our site.
					</p>
				</div>
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Cookies and Tracking Technologies</h4>
					<p class="text-gray-600 dark:text-gray-400">
						We use cookies and similar technologies to enhance your experience, analyze site usage, 
						and assist in our marketing efforts.
					</p>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- How We Use Your Information -->
		<Card.Root>
			<Card.Header>
				<Card.Title>How We Use Your Information</Card.Title>
			</Card.Header>
			<Card.Content>
				<ul class="space-y-2 text-gray-600 dark:text-gray-400">
					<li>• Process and fulfill your orders</li>
					<li>• Communicate with you about your purchases</li>
					<li>• Provide customer support</li>
					<li>• Send you marketing communications (with your consent)</li>
					<li>• Improve our website and services</li>
					<li>• Prevent fraud and ensure security</li>
				</ul>
			</Card.Content>
		</Card.Root>

		<!-- Data Protection -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Lock class="h-5 w-5" />
					How We Protect Your Information
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					We implement appropriate security measures to protect your personal information against unauthorized 
					access, alteration, disclosure, or destruction. We use industry-standard encryption for sensitive 
					data transmission and storage.
				</p>
			</Card.Content>
		</Card.Root>

		<!-- Your Rights -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Your Rights</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
					You have the right to:
				</p>
				<ul class="space-y-2 text-gray-600 dark:text-gray-400">
					<li>• Access your personal information</li>
					<li>• Correct inaccurate information</li>
					<li>• Request deletion of your information</li>
					<li>• Opt-out of marketing communications</li>
					<li>• Data portability</li>
				</ul>
			</Card.Content>
		</Card.Root>

		<!-- Contact -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Contact Us</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					If you have questions about this Privacy Policy or wish to exercise your rights, 
					please contact us at privacy@example.com or use our contact form.
				</p>
			</Card.Content>
		</Card.Root>
	</div>
</div>

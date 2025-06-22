<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { FileText, Users, Gavel } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Terms of Service</title>
	<meta name="description" content="Terms and conditions for using our website and services" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-8">
		<h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">Terms of Service</h1>
		<p class="text-lg text-gray-600 dark:text-gray-400">
			Last updated: {new Date().toLocaleDateString()}
		</p>
	</div>

	<div class="space-y-8">
		<!-- Introduction -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<FileText class="h-5 w-5" />
					Agreement to Terms
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					By accessing and using this website, you accept and agree to be bound by the terms and provision 
					of this agreement. These Terms of Service govern your use of our website and services.
				</p>
			</Card.Content>
		</Card.Root>

		<!-- Use License -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Gavel class="h-5 w-5" />
					Use License
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					Permission is granted to temporarily download one copy of the materials on our website for 
					personal, non-commercial transitory viewing only. This is the grant of a license, not a 
					transfer of title, and under this license you may not:
				</p>
				<ul class="space-y-2 text-gray-600 dark:text-gray-400 pl-4">
					<li>• Modify or copy the materials</li>
					<li>• Use the materials for any commercial purpose or for any public display</li>
					<li>• Attempt to reverse engineer any software contained on the website</li>
					<li>• Remove any copyright or other proprietary notations from the materials</li>
				</ul>
			</Card.Content>
		</Card.Root>

		<!-- User Accounts -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Users class="h-5 w-5" />
					User Accounts
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Account Responsibility</h4>
					<p class="text-gray-600 dark:text-gray-400">
						You are responsible for maintaining the confidentiality of your account and password and 
						for restricting access to your computer or device.
					</p>
				</div>
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Accurate Information</h4>
					<p class="text-gray-600 dark:text-gray-400">
						You agree to provide accurate, current, and complete information during the registration 
						process and to update such information to keep it accurate, current, and complete.
					</p>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Orders and Payments -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Orders and Payments</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Order Acceptance</h4>
					<p class="text-gray-600 dark:text-gray-400">
						We reserve the right to refuse or cancel any order for any reason at any time. 
						You may be asked to provide additional information before we accept your order.
					</p>
				</div>
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Pricing</h4>
					<p class="text-gray-600 dark:text-gray-400">
						All prices are subject to change without notice. We reserve the right to modify 
						or discontinue products at any time.
					</p>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Prohibited Uses -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Prohibited Uses</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
					You may not use our service:
				</p>
				<ul class="space-y-2 text-gray-600 dark:text-gray-400 pl-4">
					<li>• For any unlawful purpose or to solicit others to perform acts</li>
					<li>• To violate any international, federal, provincial, or state regulations, rules, laws, or local ordinances</li>
					<li>• To infringe upon or violate our intellectual property rights or the intellectual property rights of others</li>
					<li>• To harass, abuse, insult, harm, defame, slander, disparage, intimidate, or discriminate</li>
					<li>• To submit false or misleading information</li>
				</ul>
			</Card.Content>
		</Card.Root>

		<!-- Disclaimer -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Disclaimer</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					The materials on our website are provided on an 'as is' basis. We make no warranties, 
					expressed or implied, and hereby disclaim and negate all other warranties including without 
					limitation, implied warranties or conditions of merchantability, fitness for a particular 
					purpose, or non-infringement of intellectual property or other violation of rights.
				</p>
			</Card.Content>
		</Card.Root>

		<!-- Contact -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Contact Information</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-gray-600 dark:text-gray-400 leading-relaxed">
					If you have any questions about these Terms of Service, please contact us at 
					legal@example.com or through our contact form.
				</p>
			</Card.Content>
		</Card.Root>
	</div>
</div>

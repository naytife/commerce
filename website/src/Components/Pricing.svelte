<script>
	const freeFeatures = [
		'Unlimited products',
		'24/7 support',
		'Free SSL certificate',
		'Fraud analysis'
	];
	const premiumFeatures = [
		'Advanced report builder',
		'Dedicated account support',
		'Custom reporting',
		'Lower transaction fees'
	];
</script>

<div class="bg-black text-white w-max lg:py-24 py-20">
	<!-- <div class="mx-auto w-max px-6 lg:px-8"> -->
	<div class="text-center">
		<p class="mt-6 text-lg leading-8 text-gray-400">Pricing</p>
		<h2 class="text-4xl font-bold tracking-tight text-gray-100 sm:text-4xl">
			Simple no-tricks pricing
		</h2>
		<p class="mt-6 text-lg leading-8 text-gray-400">
			No additional fees or hidden costs. Get started today
		</p>
	</div>
	<div
		class="relative bg-inherit overflow-hidden pl-[2px] pr-[2px] mx-auto mt-16 max-w-2xl ring-1 ring-green-950/40 sm:mt-20 lg:mx-0 lg:grid lg:grid-cols-2 lg:max-w-none"
	>
		<!-- Content Area for Free Plan -->
		<div class="relative pl-[1px] overflow-hidden">
			<div
				class="absolute inset-0 bg-[linear-gradient(45deg,green_20%,transparent_50%)] bg-[length:200%_200%] animate-moveGradientleft z-0"
			></div>
			<div class="relative h-full bg-black lg:flex-auto lg:grid lg:grid-col p-8 sm:p-10">
				<div>
					<h3 class="text-2xl font-bold tracking-tight text-gray-100">Start</h3>
					<p class="text-sm leading-7 text-gray-400">For developers and early stage startups</p>
				</div>
				<p class="mt-6 flex items-baseline justify-start gap-x-2">
					<span class="text-5xl font-bold tracking-tight text-gray-100">Free</span>
				</p>
				<div class="mt-10 flex items-center gap-x-4">
					<h4 class="flex-none text-sm font-semibold leading-6 text-green-600">Free, Forever:</h4>
				</div>
				<div>
					<ul
						role="list"
						class="mt-8 flex flex-col gap-4 text-sm leading-6 text-gray-600 sm:grid-cols-2 sm:gap-6"
					>
						<!-- List items for Free Features -->
						{#each freeFeatures as feature}
							<li class="flex gap-x-3">
								<svg
									class="h-6 w-5 flex-none text-green-600"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
								>
									<path
										fill-rule="evenodd"
										d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
										clip-rule="evenodd"
									/>
								</svg>
								{feature}
							</li>
						{/each}
					</ul>
					<a
						href="#"
						class="mt-10 block w-full rounded-md bg-transparent border-2 border-green-600 px-3 py-2 text-center text-sm font-medium text-green-600 shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
						>Get started</a
					>
				</div>
			</div>
		</div>
		<!-- Content Area for Premium Plan -->
		<div class="relative pr-[1px] overflow-hidden">
			<!-- Gradient Background -->

			<div
				class="absolute bg-black inset-0 bg-[linear-gradient(45deg,green_20%,transparent_50%)] bg-[length:200%_200%] animate-moveGradientright z-0"
			></div>
			<div
				class="relative h-full -mt-2 bg-black lg:mt-0 lg:flex-shrink-0 p-8 sm:p-10 lg:grid lg:grid-col lg:flex-auto"
			>
				<div>
					<h3 class="text-2xl font-bold text-gray-100">Scale</h3>
					<p class="text-sm leading-7 text-gray-400">
						For fast growing teams who need additional support
					</p>
				</div>
				<p class="mt-6 flex flex-col items-baseline gap-x-2">
					<span class="text-5xl font-bold tracking-tight flex gap-3 items-end text-gray-100"
						>$599 <span class="text-sm leading-7 font-normal text-gray-400">/month</span></span
					>
					<span class="text-sm leading-7 text-gray-400"
						>Required for customers spending $50,000+ per month</span
					>
				</p>
				<div class="mt-10 flex items-center gap-x-4">
					<h4 class="flex-none text-sm font-semibold leading-6 text-green-600">
						Everything in Start, Plus:
					</h4>
				</div>
				<div>
					<ul
						role="list"
						class="mt-8 flex flex-col gap-4 text-sm leading-6 text-gray-600 sm:grid-cols-2 sm:gap-6"
					>
						<!-- List items for Premium Features -->
						{#each premiumFeatures as feature}
							<li class="flex gap-x-3">
								<svg
									class="h-6 w-5 flex-none text-green-600"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
								>
									<path
										fill-rule="evenodd"
										d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
										clip-rule="evenodd"
									/>
								</svg>
								{feature}
							</li>
						{/each}
					</ul>
					<a
						href="#"
						class="mt-10 block w-full rounded-md bg-green-600 px-3 py-2 text-center text-sm font-medium text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
						>Get started</a
					>
				</div>
			</div>
		</div>
	</div>
	<!-- </div> -->
</div>

<style>
	@keyframes moveGradientleft {
		0% {
			background-position: 0% 0%;
		}
		100% {
			background-position: -100% -100%; /* Move gradient diagonally */
		}
	}

	@keyframes moveGradientright {
		0% {
			background-position: 0% 0%;
		}
		100% {
			background-position: 100% 100%; /* Move gradient diagonally */
		}
	}

	.card::before {
		content: '';
		position: absolute;
		top: 0;
		left: -2px;
		right: 0;
		bottom: 0;
		background: linear-gradient(45deg, green 20%, transparent 45%);
		background-size: 200% 200%; /* Make gradient larger to animate */
		z-index: -1;
		animation: moveGradient 4s linear infinite;
	}
</style>

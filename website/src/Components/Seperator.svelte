<script>
	import separator from '$lib/images/seperatorAsset.svg';
	import separatorreverse from '$lib/images/seperatorAssetreverse.svg';
</script>

<div class="relative overflow-x-auto">
	<div
		class="lg:section-separator relative lg:gap-4 justify-evenly flex items-center min-w-full lg:my-32 my-20 md:mb-[64px] md:mt-[156px] is-visible"
	>
		<!-- Left gradient line -->
		<div class="flex w-1/2 h-[1px] bg-gradient-to-r from-transparent to-green-900/60"></div>

		<!-- Separator images with gradients -->
		<div class="!w-1/2 relative">
			<div class="bg-gradient-to-l from-transparent to-gray-200 absolute left-0 w-1/2 h-full"></div>
			<img alt="separator" class="w-1/2" decoding="async" loading="lazy" src={separator} />
		</div>
		<div class="!w-1/2 relative">
			<div class="bg-gradient-to-l from-gray-100/10 to-gray-200 absolute left-0 w-1/2 h-full"></div>
			<img alt="separator" class="w-1/2" decoding="async" loading="lazy" src={separator} />
		</div>

		<!-- Border divs (hidden on smaller screens) -->
		<div class="hidden w-1/2 border-gray-300 border lg:flex"></div>

		<!-- Center separator image -->
		<div class="w-1/3 lg:w-1/5 justify-start flex">
			<img
				alt="separator"
				class="w-2.5 h-3 rotate-90 object-cover"
				decoding="async"
				loading="lazy"
				src={separator}
			/>
		</div>

		<!-- Border divs (hidden on smaller screens) -->
		<div class="hidden w-1/2 border-gray-300 border lg:flex"></div>

		<!-- Separator images with gradients on the right -->
		<div class="!w-1/2 relative">
			<div
				class="bg-gradient-to-r from-transparent via-gray-200 to-transparent absolute right-0 w-full h-full"
			></div>
			<img alt="separator" class="w-1/2" decoding="async" loading="lazy" src={separatorreverse} />
		</div>
		<div class="!w-1/2 relative">
			<div
				class="bg-gradient-to-r from-transparent via-gray-200 to-transparent absolute right-0 w-full h-full"
			></div>
			<img alt="separator" class="w-1/2" decoding="async" loading="lazy" src={separatorreverse} />
		</div>

		<!-- Right gradient line -->
		<div class="flex w-1/2 h-[1px] bg-gradient-to-l from-transparent to-green-900/60"></div>
	</div>
</div>

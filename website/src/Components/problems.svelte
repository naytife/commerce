<!-- src/routes/problems.svelte -->
<script>
	// If you have any JavaScript logic, include it here
</script>

<div
	class="flex justify-center overflow-x-hidden pb-16 md:pb-48 pt-5 md:pt-8 problem-items max-w-[1512px] mx-auto"
>
	<div class="w-max">
		<div class="w-[720px] md:w-[1512px] flex flex-wrap gap-4 md:gap-6">
			<div
				style="--animation-delay: -1.7580303333432776s"
				class="w-168 bg-white-64 problem-item rounded-xl px-5 pb-4 pt-8 relative"
			>
				<div
					style="
              --background-color-hex: #DDFBFE;
              --background-color-p3: color(display-p3 0.8876 0.9798 0.9924);
            "
					class="mx-auto mb-8 w-18 h-11 rounded-full problem-item-icon"
					on:animationiteration={() => {
						// Handle animation iteration if needed
					}}
				></div>
				<div class="w-8 h-8 top-10 left-17 absolute flex items-center justify-center">
					<div class="problem-item-icon-img relative"></div>
				</div>
				<div class="text-center text-sm font-medium problem-item-label">
					How do I reduce data transfer costs?
				</div>
			</div>

			<!-- Repeat .problem-item divs as needed -->
		</div>
	</div>
</div>

<style>
	/* Scoped styles for this component */

	/* Transferring your custom styles */
	.problem-items {
		-webkit-mask-image: linear-gradient(
			to right,
			transparent 0%,
			#d9d9d9 20%,
			#d9d9d9 80%,
			transparent 100%
		);
		mask-image: linear-gradient(
			to right,
			transparent 0%,
			#d9d9d9 20%,
			#d9d9d9 80%,
			transparent 100%
		);
	}

	@media (max-width: 767px) {
		.problem-items {
			-webkit-mask-image: linear-gradient(
				to right,
				transparent 0%,
				#d9d9d9 10%,
				#d9d9d9 90%,
				transparent 100%
			);
			mask-image: linear-gradient(
				to right,
				transparent 0%,
				#d9d9d9 10%,
				#d9d9d9 90%,
				transparent 100%
			);
		}

		.problem-item:nth-child(n + 9) {
			display: none;
			animation: unset;
		}
	}

	.problem-item {
		animation: problem-item-box 7s infinite backwards;
		transform: translateZ(0);
		box-shadow:
			0 0 0 1px #0e3f7e0f,
			0 1px 1px -0.5px #2a334608,
			0 2px 2px -1px #2a33460a,
			0 3px 3px -1.5px #2a33460a,
			0 5px 5px -2.5px #2a334608,
			0 10px 10px -5px #2a334608,
			0 24px 24px -8px #2a334608;
		animation-delay: var(--animation-delay);
		background-color: var(--white-24);
	}

	@keyframes problem-item-box {
		17.8%,
		57.1% {
			background-color: var(--white-24);
		}

		25%,
		50% {
			background-color: var(--white-72);
		}
	}

	.problem-item-label {
		animation: problem-item-label 7s infinite backwards;
		opacity: 0;
		transform: translateY(10px);
		animation-delay: calc(var(--animation-delay) + 0.1s);
	}

	@keyframes problem-item-label {
		0%,
		17.8% {
			opacity: 0;
			transform: translateY(10px);
		}

		25%,
		78.5% {
			opacity: 1;
			transform: translateY(0);
			filter: blur(0);
		}

		85.6% {
			filter: blur(10px) grayscale(1);
			opacity: 0;
			transform: translateY(-20px);
		}
	}

	.problem-item-icon {
		animation: problem-item-icon 7s infinite backwards;
		opacity: 0;
		transform: translateY(20px);
		filter: blur(10px);
		animation-delay: var(--animation-delay);
		background-color: var(--background-color-p3, var(--background-color-hex));
	}

	@keyframes problem-item-icon {
		0%,
		14.2% {
			opacity: 0;
			transform: translateY(20px);
			filter: blur(10px) grayscale(1);
		}

		21.4%,
		78.5% {
			opacity: 1;
			transform: translateY(0);
			filter: blur(0);
		}

		82% {
			filter: blur(5px) grayscale(1);
		}

		to {
			opacity: 0;
			transform: translateY(-100px) scale(2);
			filter: blur(20px) grayscale(1);
		}
	}

	.problem-item-icon-img {
		animation: problem-item-icon-img 7s infinite backwards;
		animation-delay: var(--animation-delay);
		height: 0px;
		width: 0px;
		overflow: hidden;
	}

	@keyframes problem-item-icon-img {
		0%,
		17.8%,
		83% {
			height: 0px;
			width: 0px;
			transform: translateY(0);
		}

		32.1%,
		78.5% {
			height: 32px;
			width: 32px;
			opacity: 1;
			transform: translateY(0);
		}

		82% {
			opacity: 0;
			height: 32px;
			width: 32px;
			transform: translateY(-10px);
		}
	}
</style>

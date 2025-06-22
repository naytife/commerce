<script lang="ts">
	import ChevronDown from "lucide-svelte/icons/chevron-down";
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: HTMLAttributes<HTMLButtonElement> & { 
		ref?: HTMLButtonElement | null;
	} = $props();
</script>

<button
	bind:this={ref}
	class={cn("flex cursor-default items-center justify-center py-1", className)}
	{...restProps}
>
	<ChevronDown class="size-4" />
</button>

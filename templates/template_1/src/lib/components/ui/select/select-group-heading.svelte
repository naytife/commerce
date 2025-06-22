<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: HTMLAttributes<HTMLDivElement> & { 
		ref?: HTMLDivElement | null;
		children?: any;
	} = $props();
</script>

<div
	bind:this={ref}
	class={cn("py-1.5 pl-8 pr-2 text-sm font-semibold", className)}
	{...restProps}
>
	{@render children?.()}
</div>

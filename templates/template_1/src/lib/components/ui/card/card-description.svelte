<script lang="ts">
	import type { HTMLAttributes } from "svelte/elements";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: HTMLAttributes<HTMLParagraphElement> & { 
		ref?: HTMLParagraphElement | null;
		children?: any;
	} = $props();
</script>

<p bind:this={ref} class={cn("text-muted-foreground text-sm", className)} {...restProps}>
	{@render children?.()}
</p>

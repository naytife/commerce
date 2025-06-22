<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		htmlFor,
		children,
		...restProps
	}: HTMLAttributes<HTMLLabelElement> & { 
		ref?: HTMLLabelElement | null;
		htmlFor?: string;
		children?: any;
	} = $props();
</script>

<label
	bind:this={ref}
	for={htmlFor}
	class={cn(
		"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
		className
	)}
	{...restProps}
>
	{@render children?.()}
</label>

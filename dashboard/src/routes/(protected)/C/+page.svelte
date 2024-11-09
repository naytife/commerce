<script lang="ts">
	import File from 'lucide-svelte/icons/file';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import CirclePlus from 'lucide-svelte/icons/circle-plus';
	import EllipsisVertical from 'lucide-svelte/icons/ellipsis-vertical';
	import Upload from 'lucide-svelte/icons/upload';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { page } from '$app/stores';
	import type { PageData } from './$types';
	import { Header } from '$lib/components/ui/calendar';
	import CategoryForm from './category-form.svelte';

	export let data: PageData;

	$: ({ Categories } = data);

	let isVisible = false;
	let isView = true;
	let selectedCategory = null; // Store the selected category

	function toggleVisibility() {
		isVisible = !isVisible;
		isView = false;
	}

	function toggleView(category) {
		isVisible = true;
		isView = true;
		selectedCategory = category; // Set the selected category
	}
</script>

<div>
	<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
		<Tabs.Root value="all">
			<div class="mx-auto grid max-w-7xl flex-1 auto-rows-max gap-4">
				<div class="ml-auto flex items-center gap-2">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger asChild let:builder>
							<Button builders={[builder]} variant="outline" size="sm" class="h-7 gap-1">
								<ListFilter class="h-3.5 w-3.5" />
								<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Filter </span>
							</Button>
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end">
							<DropdownMenu.Label>Filter by</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.CheckboxItem checked>Active</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem>Title</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem>ID</DropdownMenu.CheckboxItem>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
					<Button size="sm" variant="outline" class="h-7 gap-1">
						<File class="h-3.5 w-3.5" />
						<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Export </span>
					</Button>
					<Button size="sm" class="h-7 gap-1" on:click={toggleVisibility}>
						<CirclePlus class="h-3.5 w-3.5" />
						<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Add Category </span>
					</Button>
				</div>
				<div
					class="grid gap-2 lg:gap-5"
					style="grid-template-columns: {isVisible ? '60% 35%' : '1fr'};"
				>
					<!-- panel 1 -->
					<div class="grid auto-rows-max items-start gap-4">
						<Card.Root>
							<Card.Header></Card.Header>
							<Card.Content>
								{#if $Categories?.data?.categories?.edges}
									<ul>
										{#each $Categories.data.categories.edges as { node: category } (category.id)}
											<li>
												<details open>
													<summary>
														<Button
															type="button"
															variant="outline"
															size="sm"
															class="mt-2 gap-2"
															on:click={() => toggleView(category)}
														>
															{category.title}
															<DropdownMenu.Root>
																<DropdownMenu.Trigger asChild let:builder>
																	<Button
																		aria-haspopup="true"
																		size="sm"
																		variant="ghost"
																		builders={[builder]}
																		on:click={(event) => event.stopPropagation()}
																	>
																		<EllipsisVertical class="h-4 w-4" />
																		<span class="sr-only">Toggle menu</span>
																	</Button>
																</DropdownMenu.Trigger>
																<DropdownMenu.Content align="end">
																	<DropdownMenu.Label>Actions</DropdownMenu.Label>
																	<DropdownMenu.Item>Delete</DropdownMenu.Item>
																</DropdownMenu.Content>
															</DropdownMenu.Root>
														</Button>
													</summary>

													{#if category.children?.length}
														<ul class="pl-20">
															{#each category.children as subcategory (subcategory.id)}
																<li>
																	<details open>
																		<summary>
																			<Button
																				type="button"
																				variant="outline"
																				size="sm"
																				class="mt-2 gap-2"
																				on:click={() => toggleView(subcategory)}
																			>
																				{subcategory.title}
																				<DropdownMenu.Root>
																					<DropdownMenu.Trigger asChild let:builder>
																						<Button
																							aria-haspopup="true"
																							size="sm"
																							variant="ghost"
																							builders={[builder]}
																							on:click={(event) => event.stopPropagation()}
																						>
																							<EllipsisVertical class="h-4 w-4" />
																							<span class="sr-only">Toggle menu</span>
																						</Button>
																					</DropdownMenu.Trigger>
																					<DropdownMenu.Content align="end">
																						<DropdownMenu.Label>Actions</DropdownMenu.Label>
																						<DropdownMenu.Item>Delete</DropdownMenu.Item>
																					</DropdownMenu.Content>
																				</DropdownMenu.Root>
																			</Button>
																		</summary>

																		{#if subcategory.children?.length}
																			<ul class="pl-20">
																				{#each subcategory.children as nestedSubcategory (nestedSubcategory.id)}
																					<li>
																						<details open>
																							<summary>
																								<Button
																									type="button"
																									variant="outline"
																									size="sm"
																									class="mt-2 gap-2"
																									on:click={() => toggleView(nestedSubcategory)}
																								>
																									{nestedSubcategory.title}
																									<DropdownMenu.Root>
																										<DropdownMenu.Trigger asChild let:builder>
																											<Button
																												aria-haspopup="true"
																												size="sm"
																												variant="ghost"
																												builders={[builder]}
																												on:click={(event) =>
																													event.stopPropagation()}
																											>
																												<EllipsisVertical class="h-4 w-4" />
																												<span class="sr-only">Toggle menu</span>
																											</Button>
																										</DropdownMenu.Trigger>
																										<DropdownMenu.Content align="end">
																											<DropdownMenu.Label
																												>Actions</DropdownMenu.Label
																											>
																											<DropdownMenu.Item>Delete</DropdownMenu.Item>
																										</DropdownMenu.Content>
																									</DropdownMenu.Root>
																								</Button>
																							</summary>

																							{#if nestedSubcategory.children?.length}
																								<ul class="pl-20">
																									{#each nestedSubcategory.children as deepSubcategory (deepSubcategory.id)}
																										<li>
																											<Button
																												type="button"
																												variant="outline"
																												size="sm"
																												class="mt-2 gap-2"
																												on:click={() => toggleView(deepSubcategory)}
																											>
																												{deepSubcategory.title}
																											</Button>
																										</li>
																									{/each}
																								</ul>
																							{/if}
																						</details>
																					</li>
																				{/each}
																			</ul>
																		{/if}
																	</details>
																</li>
															{/each}
														</ul>
													{/if}
												</details>
											</li>
										{/each}
									</ul>
								{:else}
									<div class="pl-10 leading-6 text-slate-600 dark:text-slate-400">
										<p class="mt-2 text-xs text-white">
											There are no categories available in this shop
										</p>
										<p>
											<button on:click={toggleVisibility} class="text-xs underline">
												add category
											</button>
										</p>
									</div>
								{/if}
							</Card.Content>
						</Card.Root>
					</div>

					<!-- panel 2 -->
					{#if isVisible}
						<div class="grid auto-rows-max items-start gap-4">
							<Card.Root class="overflow-hidden">
								<Card.Header>
									<Card.Title>Category details</Card.Title>
								</Card.Header>
								<Card.Content>
									{#if !isView}
										<CategoryForm data={data?.form || {}} />
									{:else}
										<div class="grid gap-6">
											{#if selectedCategory}
												<div class="grid gap-3">
													<Card.Title>Title</Card.Title>
													<Card.Description>{selectedCategory.title}</Card.Description>
												</div>

												<div class="grid gap-3">
													<Card.Title>Description</Card.Title>
													<Card.Description>{selectedCategory.description}</Card.Description>
												</div>

												<div class="grid gap-3">
													<Card.Title>Sub category</Card.Title>
													<Card.Description>
														{#if selectedCategory.children?.length}
															<ul>
																{#each selectedCategory.children as child}
																	<li>{child.title}</li>
																{/each}
															</ul>
														{:else}
															no subcategories
														{/if}
													</Card.Description>
												</div>
												<div class="grid gap-2">
													<img
														alt="Product"
														class="aspect-square w-full rounded-md object-cover"
														height="300"
														src="/images/placeholder.png"
														width="300"
													/>
												</div>
											{:else}
												<div>Select a category to view its details</div>
											{/if}
										</div>
									{/if}
								</Card.Content>
							</Card.Root>
						</div>
					{/if}
				</div>
			</div>
		</Tabs.Root>
	</main>
</div>

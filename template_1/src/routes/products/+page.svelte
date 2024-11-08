<script lang="ts">
	import File from 'lucide-svelte/icons/file';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import CirclePlus from 'lucide-svelte/icons/circle-plus';
	import { Button } from '$lib/components/ui/button';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Dialog from '$lib/components/ui/dialog';
	import { HomeIcon, SortAsc } from 'lucide-svelte';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Upload from 'lucide-svelte/icons/upload';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as ToggleGroup from '$lib/components/ui/toggle-group';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	let Electronics: any;
</script>

<section class="bg-gray-50 py-8 antialiased md:py-12 dark:bg-gray-900">
	<div class="mx-auto max-w-screen-xl px-4 2xl:px-0">
		<!-- Heading & Filters -->
		<div class="mb-4 items-end justify-between space-y-4 sm:flex sm:space-y-0 md:mb-8">
			<div>
				<Breadcrumb.Root class="hidden md:flex">
					<Breadcrumb.List>
						<Breadcrumb.Item>
							<Breadcrumb.Link
								class="hover:text-primary-600 inline-flex items-center text-sm font-medium text-gray-700 dark:text-gray-400 dark:hover:text-white"
								href="/admin/"
							>
								<HomeIcon class="me-2.5 h-4 w-4" /> Home</Breadcrumb.Link
							>
						</Breadcrumb.Item>
						<Breadcrumb.Separator />
						<Breadcrumb.Item
							class="hover:text-primary-600 inline-flex items-center text-sm font-medium text-gray-700 dark:text-gray-400 dark:hover:text-white"
						>
							<Breadcrumb.Link href="/admin/products">Products</Breadcrumb.Link>
						</Breadcrumb.Item>
						<Breadcrumb.Separator />
						<Breadcrumb.Item
							class="hover:text-primary-600 inline-flex items-center text-sm font-medium text-gray-700 dark:text-gray-400 dark:hover:text-white"
						>
							<Breadcrumb.Link href="/admin/categories">Electronics</Breadcrumb.Link>
						</Breadcrumb.Item>
					</Breadcrumb.List>
				</Breadcrumb.Root>
				<h2 class="mt-3 text-xl font-semibold text-gray-900 sm:text-2xl dark:text-white">
					Electronics
				</h2>
			</div>

			<div class="ml-auto flex items-center gap-2">
				<Dialog.Root>
					<Dialog.Trigger>
						<Button
							variant="outline"
							size="sm"
							class="h-7 gap-1 px-5 py-4 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-100"
						>
							<ListFilter class="h-3.5 w-3.5" />
							<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Filters </span>
						</Button>
					</Dialog.Trigger>
					<Dialog.Content>
						<Dialog.Header>
							<Dialog.Title>Filters</Dialog.Title>
							<div class="flex gap-5">
								<div>Brand</div>
								<div>Advanced filters</div>
							</div>
						</Dialog.Header>
						<hr />
						<div>
							<div class="grid gap-6">
								<div class="flex gap-3">
									<Checkbox id="terms" />
									<Label
										id="terms-label"
										for="terms"
										class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
									>
										Accessories
									</Label>
								</div>
							</div>
						</div>
						<Dialog.Footer>
							<Button type="submit">Show 50 results</Button>
							<Button type="reset">Reset</Button>
						</Dialog.Footer>
					</Dialog.Content>
				</Dialog.Root>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger asChild let:builder>
						<Button
							builders={[builder]}
							size="sm"
							variant="outline"
							class="h-7 gap-1 px-5 py-4 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-100"
						>
							<SortAsc class="h-3.5 w-3.5" />
							<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Sort </span>
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="end">
						<DropdownMenu.Label>Filter by</DropdownMenu.Label>
						<DropdownMenu.Separator />
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							checked>The most popular</DropdownMenu.CheckboxItem
						>
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>Newest</DropdownMenu.CheckboxItem
						>
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>Increasing price</DropdownMenu.CheckboxItem
						>
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>Decreasing price</DropdownMenu.CheckboxItem
						>
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>No. reviews</DropdownMenu.CheckboxItem
						>
						<DropdownMenu.CheckboxItem
							class="group  items-center rounded-md hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>Discount %</DropdownMenu.CheckboxItem
						>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</div>

		<div class="mb-4 grid gap-4 sm:grid-cols-2 md:mb-8 lg:grid-cols-3 xl:grid-cols-4">
			{#each Electronics as item}
				<div
					class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-700 dark:bg-gray-800"
				>
					<div class="h-56 w-full">
						<a href="#">
							<img
								class="mx-auto h-full dark:hidden"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
								alt=""
							/>
							<img
								class="mx-auto hidden h-full dark:block"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front-dark.svg"
								alt=""
							/>
						</a>
					</div>
					<div class="pt-6">
						<div class="mb-4 flex items-center justify-between gap-4">
							<span
								class="bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-300 me-2 rounded px-2.5 py-0.5 text-xs font-medium"
							>
								Up to 35% off
							</span>

							<div class="flex items-center justify-end gap-1">
								<button
									type="button"
									data-tooltip-target="tooltip-quick-look"
									class="rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
								>
									<span class="sr-only"> Quick look </span>
									<svg
										class="h-5 w-5"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="none"
										viewBox="0 0 24 24"
									>
										<path
											stroke="currentColor"
											stroke-width="2"
											d="M21 12c0 1.2-4.03 6-9 6s-9-4.8-9-6c0-1.2 4.03-6 9-6s9 4.8 9 6Z"
										/>
										<path
											stroke="currentColor"
											stroke-width="2"
											d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
										/>
									</svg>
								</button>
								<div
									id="tooltip-quick-look"
									role="tooltip"
									class="tooltip invisible absolute z-10 inline-block rounded-lg bg-gray-900 px-3 py-2 text-sm font-medium text-white opacity-0 shadow-sm transition-opacity duration-300 dark:bg-gray-700"
									data-popper-placement="top"
								>
									Quick look
									<div class="tooltip-arrow" data-popper-arrow=""></div>
								</div>

								<button
									type="button"
									data-tooltip-target="tooltip-add-to-favorites"
									class="rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
								>
									<span class="sr-only"> Add to Favorites </span>
									<svg
										class="h-5 w-5"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
									>
										<path
											stroke="currentColor"
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M12 6C6.5 1 1 8 5.8 13l6.2 7 6.2-7C23 8 17.5 1 12 6Z"
										/>
									</svg>
								</button>
								<div
									id="tooltip-add-to-favorites"
									role="tooltip"
									class="tooltip invisible absolute z-10 inline-block rounded-lg bg-gray-900 px-3 py-2 text-sm font-medium text-white opacity-0 shadow-sm transition-opacity duration-300 dark:bg-gray-700"
									data-popper-placement="top"
								>
									Add to favorites
									<div class="tooltip-arrow" data-popper-arrow=""></div>
								</div>
							</div>
						</div>

						<a
							href="#"
							class="text-lg font-semibold leading-tight text-gray-900 hover:underline dark:text-white"
							>Apple iMac 27", 1TB HDD, Retina 5K Display, M3 Max</a
						>

						<div class="mt-2 flex items-center gap-2">
							<div class="flex items-center">
								<svg
									class="h-4 w-4 text-yellow-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										d="M13.8 4.2a2 2 0 0 0-3.6 0L8.4 8.4l-4.6.3a2 2 0 0 0-1.1 3.5l3.5 3-1 4.4c-.5 1.7 1.4 3 2.9 2.1l3.9-2.3 3.9 2.3c1.5 1 3.4-.4 3-2.1l-1-4.4 3.4-3a2 2 0 0 0-1.1-3.5l-4.6-.3-1.8-4.2Z"
									/>
								</svg>

								<svg
									class="h-4 w-4 text-yellow-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										d="M13.8 4.2a2 2 0 0 0-3.6 0L8.4 8.4l-4.6.3a2 2 0 0 0-1.1 3.5l3.5 3-1 4.4c-.5 1.7 1.4 3 2.9 2.1l3.9-2.3 3.9 2.3c1.5 1 3.4-.4 3-2.1l-1-4.4 3.4-3a2 2 0 0 0-1.1-3.5l-4.6-.3-1.8-4.2Z"
									/>
								</svg>

								<svg
									class="h-4 w-4 text-yellow-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										d="M13.8 4.2a2 2 0 0 0-3.6 0L8.4 8.4l-4.6.3a2 2 0 0 0-1.1 3.5l3.5 3-1 4.4c-.5 1.7 1.4 3 2.9 2.1l3.9-2.3 3.9 2.3c1.5 1 3.4-.4 3-2.1l-1-4.4 3.4-3a2 2 0 0 0-1.1-3.5l-4.6-.3-1.8-4.2Z"
									/>
								</svg>

								<svg
									class="h-4 w-4 text-yellow-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										d="M13.8 4.2a2 2 0 0 0-3.6 0L8.4 8.4l-4.6.3a2 2 0 0 0-1.1 3.5l3.5 3-1 4.4c-.5 1.7 1.4 3 2.9 2.1l3.9-2.3 3.9 2.3c1.5 1 3.4-.4 3-2.1l-1-4.4 3.4-3a2 2 0 0 0-1.1-3.5l-4.6-.3-1.8-4.2Z"
									/>
								</svg>

								<svg
									class="h-4 w-4 text-yellow-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										d="M13.8 4.2a2 2 0 0 0-3.6 0L8.4 8.4l-4.6.3a2 2 0 0 0-1.1 3.5l3.5 3-1 4.4c-.5 1.7 1.4 3 2.9 2.1l3.9-2.3 3.9 2.3c1.5 1 3.4-.4 3-2.1l-1-4.4 3.4-3a2 2 0 0 0-1.1-3.5l-4.6-.3-1.8-4.2Z"
									/>
								</svg>
							</div>

							<p class="text-sm font-medium text-gray-900 dark:text-white">5.0</p>
							<p class="text-sm font-medium text-gray-500 dark:text-gray-400">(455)</p>
						</div>

						<ul class="mt-2 flex items-center gap-4">
							<li class="flex items-center gap-2">
								<svg
									class="h-4 w-4 text-gray-500 dark:text-gray-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
								>
									<path
										stroke="currentColor"
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M13 7h6l2 4m-8-4v8m0-8V6a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v9h2m8 0H9m4 0h2m4 0h2v-4m0 0h-5m3.5 5.5a2.5 2.5 0 1 1-5 0 2.5 2.5 0 0 1 5 0Zm-10 0a2.5 2.5 0 1 1-5 0 2.5 2.5 0 0 1 5 0Z"
									/>
								</svg>
								<p class="text-sm font-medium text-gray-500 dark:text-gray-400">Fast Delivery</p>
							</li>

							<li class="flex items-center gap-2">
								<svg
									class="h-4 w-4 text-gray-500 dark:text-gray-400"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
								>
									<path
										stroke="currentColor"
										stroke-linecap="round"
										stroke-width="2"
										d="M8 7V6c0-.6.4-1 1-1h11c.6 0 1 .4 1 1v7c0 .6-.4 1-1 1h-1M3 18v-7c0-.6.4-1 1-1h11c.6 0 1 .4 1 1v7c0 .6-.4 1-1 1H4a1 1 0 0 1-1-1Zm8-3.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"
									/>
								</svg>
								<p class="text-sm font-medium text-gray-500 dark:text-gray-400">Best Price</p>
							</li>
						</ul>

						<div class="mt-4 flex items-center justify-between gap-4">
							<p class="text-2xl font-extrabold leading-tight text-gray-900 dark:text-white">
								$1,699
							</p>

							<button
								type="button"
								class="bg-primary-700 hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800 inline-flex items-center rounded-lg px-5 py-2.5 text-sm font-medium text-white focus:outline-none focus:ring-4"
							>
								<svg
									class="-ms-2 me-2 h-5 w-5"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									width="24"
									height="24"
									fill="none"
									viewBox="0 0 24 24"
								>
									<path
										stroke="currentColor"
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M4 4h1.5L8 16m0 0h8m-8 0a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm8 0a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm.75-3H7.5M11 7H6.312M17 4v6m-3-3h6"
									/>
								</svg>
								Add to cart
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
		<div class="w-full text-center">
			<button
				type="button"
				class="hover:text-primary-700 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
				>Show more</button
			>
		</div>
	</div>
	<!-- Filter modal -->
</section>

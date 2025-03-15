<script lang="ts">

	import * as Card from '$lib/components/ui/card';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Info, MapPin } from 'lucide-svelte';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import {Button} from '$lib/components/ui/button';

	// import { formSchema, type FormSchema } from './schema';

	export let data: { type: string }; // Define the form data type
	const phoneNumber = '09150610759';
</script>

<section class="bg-white py-8 antialiased dark:bg-gray-900 md:py-16">
	<div class="mx-auto max-w-screen-xl px-4 2xl:px-0">
		<div class="lg: grid gap-4 md:grid-cols-[60%_40%] lg:gap-10">
			<div class="flex flex-col gap-4">
				<div class="gap-2 lg:grid lg:grid-cols-4">
					<div class="col-span-1 flex h-28 w-28 flex-col gap-4">
						<button class="rounded-md p-2 ring-[1.5px] ring-gray-400">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
							/>
						</button>
						<button class="rounded-md p-2 ring-[1.5px] ring-gray-400">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
							/>
						</button>
						<button class="rounded-md p-2 ring-[1.5px] ring-gray-400">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
							/>
						</button>
						<button class="rounded-md p-2 ring-[1.5px] ring-gray-400">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
							/>
						</button>
					</div>
					<img
						alt="Product"
						class="col-span-4 col-start-2 aspect-square w-full rounded-md object-cover"
						height="300"
						src="https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg"
						width="300"
					/>
				</div>
				<div class="">
					<Accordion.Root class="max-w-[70%] md:w-full md:max-w-full">
						<Accordion.Item value="item-1">
							<Accordion.Trigger>Product Details</Accordion.Trigger>
							<Accordion.Content>
								The product is a high-quality, durable solution designed to meet the needs of modern
								consumers. It features advanced technology and ergonomic design for optimal
								performance and comfort. Key features include a sleek interface, customizable
								settings, and compatibility with various devices. It is ideal for professionals and
								enthusiasts alike.
							</Accordion.Content>
						</Accordion.Item>
						<Accordion.Item value="item-2">
							<Accordion.Trigger>Specifications</Accordion.Trigger>
							<Accordion.Content>
								Specifications include advanced processor capabilities, extensive storage options,
								and seamless connectivity. The product is designed to handle demanding tasks
								efficiently. Additional features include high-resolution display, long-lasting
								battery life, and intuitive user interface for enhanced productivity.
							</Accordion.Content>
						</Accordion.Item>
						<Accordion.Item value="item-3">
							<Accordion.Trigger>Warranty</Accordion.Trigger>
							<Accordion.Content
								>The product comes with a comprehensive warranty that covers manufacturing defects
								and malfunctions. It provides peace of mind and assurance of product quality.
								Customers can enjoy reliable support and prompt assistance for any issues related to
								the product during the warranty period.</Accordion.Content
							>
						</Accordion.Item>
					</Accordion.Root>
				</div>
			</div>
			<div class="mt-6 sm:mt-8 lg:mt-0">
				<Card.Root>
					<Card.Content>
						<h1 class="text-xl font-semibold text-gray-900 dark:text-white sm:text-2xl">
							Apple iMac 24" All-In-One Computer, Apple M1, 8GB RAM, 256GB SSD, Mac OS, Pink
						</h1>
						<div class="mt-4 sm:flex sm:items-center sm:gap-4">
							<p
								class="flex items-center rounded-sm bg-red-600 px-2 py-1 text-xs font-medium leading-none text-white"
							>
								the last 2 products
							</p>
							<div class="mt-2 flex items-center gap-2 sm:mt-0">
								<div class="flex items-center gap-1">
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
								</div>

								<a
									href="/"
									class="text-sm font-medium leading-none text-gray-900 underline hover:no-underline dark:text-white"
								>
									345 Reviews
								</a>
							</div>
						</div>
						<p class="mt-4 flex items-center gap-2 text-sm font-medium leading-none text-blue-600">
							<MapPin class="blue-600" /> Deliver to Bonnie Green- Sacramento 23647
						</p>
						<div class="mt-6 flex items-center justify-between">
							<p class="text-2xl font-extrabold text-gray-900 dark:text-white sm:text-3xl">
								$1,249.99
							</p>
							<div class="flex items-center gap-4">
								<p class="flex items-center gap-1 text-sm font-semibold">
									Quantity <Info class="h-5 w-4" />
								</p>
								<Select.Root>
									<Select.Trigger class="w-16 p-4">
										<Select.Value placeholder="0" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item aria-selected value="0">0</Select.Item>
										<Select.Item value="1">1</Select.Item>
										<Select.Item value="2">2</Select.Item>
										<Select.Item value="3">3</Select.Item>
										<Select.Item value="4">4</Select.Item>
										<Select.Item value="5">5</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>
						</div>

						<div class="mt-6 flex flex-col sm:mt-8 sm:flex sm:items-center sm:gap-4">
							<a
								href="/"
								title=""
								class="hover:text-primary-700 flex w-full items-center justify-center rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
								role="button"
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
										d="M12.01 6.001C6.5 1 1 8 5.782 13.001L12.011 20l6.23-7C23 8 17.5 1 12.01 6.002Z"
									/>
								</svg>
								Add to favorites
							</a>

							<a
								href={`https://wa.me/${phoneNumber}`}
								target="_blank"
								rel="noopener noreferrer"
								title=""
								class="hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800 mt-4 flex w-full items-center justify-center rounded-lg bg-green-600 px-5 py-2.5 text-sm font-medium text-white focus:outline-none focus:ring-4 sm:mt-0"
								role="button"
							>
								Continue to WhatsApp
							</a>
						</div>
						<p class="mt-6 text-sm">
							Also available at competitive prices from <a
								class="text-blue-600"
								href="http://retailauhtority.com">authorized retailers</a
							>, with optional Premium delivery for expedited shipping
						</p>

						<hr class="my-6 border-gray-200 dark:border-gray-800 md:my-8" />
						<div class="flex w-full flex-col gap-2">
							<p class="text-lg font-bold">Colour</p>
							<div class="flex gap-2">
								<Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">Green</Button
								><Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">Pink</Button
								><Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600"
									>Silver</Button
								><Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">Blue</Button
								>
							</div>
						</div>
						<div class="mt-6 flex w-full flex-col gap-2">
							<p class="text-lg font-bold">SSD capacity</p>
							<div class="flex gap-2">
								<Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">256GB</Button
								>
								<Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">512GB</Button
								>
								<Button class="bg-neutral-500 hover:bg-neutral-400 focus:bg-green-600">1TB</Button>
							</div>
						</div>
						<div class="mt-6 w-full gap-4">
							<p class="text-lg font-bold">Pickup</p>
							<RadioGroup.Root bind:value={data.type} class="flex flex-col space-y-1">
								<div class="flex items-center space-x-3 space-y-0">
									<RadioGroup.Item value="all" />
									<div>
										<Label class="font-normal">Shipping - $19</Label>
										<p class="text-xs text-gray-500">Arrives Nov 17</p>
									</div>
								</div>
								<div class="flex items-center space-x-3 space-y-0">
									<RadioGroup.Item value="mentions" />
									<div>
										<Label class="font-normal">Pickup from Flowbox- $9</Label>
										<p><a href="/" class="text-xs text-blue-600">Pick a Flowbox near you</a></p>
									</div>
								</div>
								<div class="flex items-center space-x-3 space-y-0">
									<RadioGroup.Item value="" />
									<div>
										<Label class="font-normal">Pickup from our store</Label>
										<p class="text-xs text-gray-500">Not available</p>
									</div>
								</div>
								<RadioGroup.Input name="type" />
							</RadioGroup.Root>
						</div>
						<div class="mt-6 w-full gap-4 space-y-1">
							<p class="text-lg font-bold">Add extra warranty</p>
							<RadioGroup.Root bind:value={data.type} class="flex space-y-1">
								<div class="flex items-center gap-2 space-y-0">
									<RadioGroup.Item value="1year" />
									<Label class="font-semibold">1 year - $39</Label>
								</div>
								<div class="flex items-center gap-2 space-y-0">
									<RadioGroup.Item value="2year" />
									<Label class="font-semibold">2 years - $69</Label>
								</div>
								<div class="flex items-center gap-2 space-y-0">
									<RadioGroup.Item value="3year" />
									<Label class="font-semibold">3 years - $991</Label>
								</div>
								<RadioGroup.Input name="type" />
							</RadioGroup.Root>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
	</div>
</section>

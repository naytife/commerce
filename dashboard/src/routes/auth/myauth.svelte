<script lang="ts">
	import UserAuthForm from './(components)/user-auth-form.svelte';
	import { Button } from '$lib/components/ui/button';

	// import AuthenticationLight from '$lib/img/authentication-light.png?enhanced';
	// import AuthenticationDark from '$lib/img/authentication-dark.png?enhanced';
</script>

<!-- <div class="md:hidden">
	<enhanced:img src={AuthenticationLight} alt="Authentication" class="block dark:hidden" />
	<enhanced:img src={AuthenticationDark} alt="Authentication" class="hidden dark:block" />
</div> -->
<div
	class="container relative hidden h-svh flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>
	<Button
		href="/examples/authentication"
		variant="ghost"
		class="absolute right-4 top-4 md:right-8 md:top-8"
	>
		Login
	</Button>
	<div class="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
		<div
			class="absolute inset-0 bg-cover"
			style="
				background-image:
					url(https://plus.unsplash.com/premium_photo-1682432005975-4afd1fff1134?q=80&w=1770&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D);"
		/>
		<div class="relative z-20 flex items-center text-lg font-medium">
			<!-- <Command class="mr-2 h-6 w-6" /> -->
			Ashia Inc
		</div>
	</div>
	<div class="lg:p-8">
		<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">Account Login</h1>
				<p class="text-sm text-muted-foreground">
					Enter your credentials below to login to your account
				</p>
			</div>
			<UserAuthForm />
			<p class="px-8 text-center text-sm text-muted-foreground">
				By clicking continue, you agree to our
				<a href="/terms" class="underline underline-offset-4 hover:text-primary">
					Terms of Service
				</a>
				and
				<a href="/privacy" class="underline underline-offset-4 hover:text-primary">
					Privacy Policy
				</a>
				.
			</p>
		</div>
	</div>
</div>

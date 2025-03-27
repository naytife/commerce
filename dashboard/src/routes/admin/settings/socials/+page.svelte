<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
</script>

<div class="grid gap-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Social Settings</Card.Title>
			<Card.Description>Used to identify your store in the marketplace.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" class="space-y-8"  id="general-settings-form">
				<div class="flex w-full flex-col gap-3">
					<Label for="whatsapp">WhatsApp</Label>
					<Input id="whatsapp" type="text" />
					<p class="text-muted-foreground text-sm">WhatsApp link.</p>
				</div>
				<div class="flex w-full flex-col gap-3">
					<Label for="facebook">Facebook</Label>
					<Input id="facebook" type="text" />
					<p class="text-muted-foreground text-sm">Facebook link.</p>
				</div>
				<div class="flex w-full flex-col gap-3">
					<Label for="instagram">Instagram</Label>
					<Input id="instagram" type="text" />
					<p class="text-muted-foreground text-sm">Instagram link.</p>
				</div>
				<div class="flex w-full flex-col gap-3">
					<Label for="tiktok">TikTok</Label>
					<Input id="tiktok" type="text" />
					<p class="text-muted-foreground text-sm">TikTok link.</p>
				</div>
			</form>
		</Card.Content>
		<Card.Footer class="border-t px-6 py-4">
			<Button>Save</Button>
		</Card.Footer>
	</Card.Root>
</div>

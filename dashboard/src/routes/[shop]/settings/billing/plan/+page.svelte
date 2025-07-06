<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { getCurrencySymbol } from '$lib/utils/currency';
	import ComingSoon from '$lib/components/ui/coming-soon.svelte';

  // Static plan data for now
  const currentPlan = {
    name: 'Free',
    price: 0,
    nextBilling: '2024-07-01',
    status: 'Active',
    currency: 'USD'
  };

  const availablePlans = [
    {
      name: 'Free',
      price: 0,
      features: [
        'Up to 10 products',
        'Basic analytics',
        'Community support'
      ],
      currency: 'USD',
      status: currentPlan.name === 'Free' ? 'Active' : 'Available'
    },
    {
      name: 'Pro',
      price: 49.99,
      features: [
        'Unlimited products',
        'Priority support',
        'Advanced analytics',
        'Custom domain',
        '24/7 support'
      ],
      currency: 'USD',
      status: currentPlan.name === 'Pro' ? 'Active' : 'Available'
    },
    {
      name: 'Enterprise',
      price: 199.99,
      features: [
        'Everything in Pro',
        'Dedicated account manager',
        'Custom integrations',
        'SLA & priority support'
      ],
      currency: 'USD',
      status: currentPlan.name === 'Enterprise' ? 'Active' : 'Available'
    }
  ];

  function handleChangePlan(planName: string) {
    // TODO: Implement plan change logic
    alert(`Change to ${planName} plan (not implemented)`);
  }
</script>

<div class="space-y-6 max-w-2xl mx-auto">
  <h2 class="text-2xl font-bold">Subscription Plan</h2>
  <Card.Root>
    <Card.Header>
      <Card.Title>Current Plan</Card.Title>
      <Card.Description>Your current subscription plan</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="p-4 border rounded-md">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="font-medium">{currentPlan.name} Plan</h3>
            <p class="text-sm text-muted-foreground mb-1">Next billing: {currentPlan.nextBilling}</p>
          </div>
          <span class="text-sm font-semibold bg-primary/10 text-primary px-2 py-1 rounded">{currentPlan.status}</span>
        </div>
        <div class="mt-2 text-lg font-bold">{getCurrencySymbol(currentPlan.currency)}{currentPlan.price} <span class="text-sm font-normal text-muted-foreground">/month</span></div>
      </div>
    </Card.Content>
  </Card.Root>
<ComingSoon>
  <Card.Root>
    <Card.Header>
      <Card.Title>Available Plans</Card.Title>
      <Card.Description>Upgrade or downgrade your subscription</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="grid gap-4 md:grid-cols-3">
        {#each availablePlans as plan}
          <div class="p-4 border rounded-md hover:border-primary cursor-pointer flex flex-col justify-between">
            <div>
              <h4 class="font-medium">{plan.name} Plan</h4>
              <ul class="text-sm mt-2 ml-4 list-disc mb-2">
                {#each plan.features as feature}
                  <li>{feature}</li>
                {/each}
              </ul>
            </div>
            <div class="mt-auto">
              <div class="text-xl font-bold">{getCurrencySymbol(plan.currency)}{plan.price}</div>
              <div class="text-sm text-muted-foreground mb-2">/month</div>
              {#if plan.status === 'Active'}
                <Button class="w-full" variant="outline" disabled>Current Plan</Button>
              {:else}
                <Button class="w-full" variant="outline" on:click={() => handleChangePlan(plan.name)}>Select Plan</Button>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    </Card.Content>
  </Card.Root>
</ComingSoon>
</div>
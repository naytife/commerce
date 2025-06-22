<script>
  import { page } from '$app/stores';
</script>

<svelte:head>
  <title>Shipping Information - Shop</title>
  <meta name="description" content="Learn about our shipping options and delivery information" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Shipping Information</h1>

    <!-- Shipping Options -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Shipping Options</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="border border-gray-200 rounded-lg p-4">
          <h3 class="font-semibold text-gray-900 mb-2">Standard Shipping</h3>
          <p class="text-2xl font-bold text-blue-600 mb-2">Free</p>
          <p class="text-sm text-gray-600 mb-2">5-7 business days</p>
          <p class="text-xs text-gray-500">Free on orders over $50</p>
        </div>
        <div class="border border-gray-200 rounded-lg p-4">
          <h3 class="font-semibold text-gray-900 mb-2">Express Shipping</h3>
          <p class="text-2xl font-bold text-blue-600 mb-2">$9.99</p>
          <p class="text-sm text-gray-600 mb-2">2-3 business days</p>
          <p class="text-xs text-gray-500">Available for most items</p>
        </div>
        <div class="border border-gray-200 rounded-lg p-4">
          <h3 class="font-semibold text-gray-900 mb-2">Next Day Delivery</h3>
          <p class="text-2xl font-bold text-blue-600 mb-2">$19.99</p>
          <p class="text-sm text-gray-600 mb-2">1 business day</p>
          <p class="text-xs text-gray-500">Order by 2 PM</p>
        </div>
      </div>
    </div>

    <!-- Processing Time -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Processing Time</h2>
      <p class="text-gray-600 mb-4">
        Most orders are processed within 1-2 business days. Custom or personalized items may require additional processing time.
      </p>
      <ul class="list-disc list-inside text-gray-600 space-y-2">
        <li>Standard items: 1-2 business days</li>
        <li>Custom items: 3-5 business days</li>
        <li>Large items: 2-3 business days</li>
        <li>Electronics: 1-2 business days</li>
      </ul>
    </div>

    <!-- International Shipping -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">International Shipping</h2>
      <p class="text-gray-600 mb-4">
        We ship to over 100 countries worldwide. International shipping costs and delivery times vary by destination.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">Shipping Costs</h3>
          <ul class="text-gray-600 space-y-1 text-sm">
            <li>Canada: Starting at $12.99</li>
            <li>Europe: Starting at $19.99</li>
            <li>Asia: Starting at $24.99</li>
            <li>Other regions: Starting at $29.99</li>
          </ul>
        </div>
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">Delivery Times</h3>
          <ul class="text-gray-600 space-y-1 text-sm">
            <li>Canada: 7-14 business days</li>
            <li>Europe: 10-21 business days</li>
            <li>Asia: 14-28 business days</li>
            <li>Other regions: 21-35 business days</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Tracking -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Order Tracking</h2>
      <p class="text-gray-600 mb-4">
        Once your order ships, you'll receive a tracking number via email. You can track your package using our order tracking system or directly on the carrier's website.
      </p>
      <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <p class="text-blue-800 text-sm">
          <strong>Tip:</strong> Create an account to easily track all your orders in one place and receive automatic updates on delivery status.
        </p>
      </div>
    </div>

    <!-- Special Items -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Special Items</h2>
      <div class="space-y-4">
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">Large Items</h3>
          <p class="text-gray-600 text-sm">Furniture and large appliances may require special delivery arrangements. We'll contact you to schedule delivery.</p>
        </div>
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">Hazardous Materials</h3>
          <p class="text-gray-600 text-sm">Some items containing batteries or chemicals may have shipping restrictions and longer delivery times.</p>
        </div>
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">Fresh/Perishable Items</h3>
          <p class="text-gray-600 text-sm">Perishable items are shipped with expedited delivery and require someone to receive the package.</p>
        </div>
      </div>
    </div>
  </div>
</div>

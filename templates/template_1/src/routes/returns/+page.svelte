<script>
  import { page } from '$app/stores';
</script>

<svelte:head>
  <title>Returns & Exchanges - Shop</title>
  <meta name="description" content="Learn about our return policy and how to process returns" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Returns & Exchanges</h1>

    <!-- Return Policy Overview -->
    <div class="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-8">
      <h2 class="text-xl font-semibold text-blue-900 mb-4">Our Return Promise</h2>
      <p class="text-blue-800 mb-4">
        We want you to be completely satisfied with your purchase. If you're not happy with your order, you can return most items within 30 days for a full refund or exchange.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
        <div class="bg-white rounded-lg p-4">
          <div class="text-2xl mb-2">📅</div>
          <h3 class="font-semibold text-blue-900">30 Days</h3>
          <p class="text-sm text-blue-700">Return window</p>
        </div>
        <div class="bg-white rounded-lg p-4">
          <div class="text-2xl mb-2">🆓</div>
          <h3 class="font-semibold text-blue-900">Free Returns</h3>
          <p class="text-sm text-blue-700">On orders over $50</p>
        </div>
        <div class="bg-white rounded-lg p-4">
          <div class="text-2xl mb-2">💰</div>
          <h3 class="font-semibold text-blue-900">Full Refund</h3>
          <p class="text-sm text-blue-700">Original payment method</p>
        </div>
      </div>
    </div>

    <!-- How to Return -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">How to Return an Item</h2>
      <div class="space-y-6">
        <div class="flex items-start">
          <div class="flex-shrink-0 w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-semibold mr-4">1</div>
          <div>
            <h3 class="font-semibold text-gray-900 mb-2">Start Your Return</h3>
            <p class="text-gray-600">Log into your account and go to "Order History" to initiate a return, or contact our customer service team.</p>
          </div>
        </div>
        <div class="flex items-start">
          <div class="flex-shrink-0 w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-semibold mr-4">2</div>
          <div>
            <h3 class="font-semibold text-gray-900 mb-2">Print Return Label</h3>
            <p class="text-gray-600">We'll email you a prepaid return shipping label. Print it and attach it to your package.</p>
          </div>
        </div>
        <div class="flex items-start">
          <div class="flex-shrink-0 w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-semibold mr-4">3</div>
          <div>
            <h3 class="font-semibold text-gray-900 mb-2">Package & Ship</h3>
            <p class="text-gray-600">Pack the item securely in its original packaging (if possible) and drop it off at any authorized shipping location.</p>
          </div>
        </div>
        <div class="flex items-start">
          <div class="flex-shrink-0 w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-semibold mr-4">4</div>
          <div>
            <h3 class="font-semibold text-gray-900 mb-2">Get Your Refund</h3>
            <p class="text-gray-600">Once we receive and process your return, your refund will be issued within 3-5 business days.</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Return Conditions -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Return Conditions</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="font-semibold text-green-600 mb-3">✅ Returnable Items</h3>
          <ul class="text-gray-600 space-y-2 text-sm">
            <li>• Items in original condition</li>
            <li>• Unused and unworn items</li>
            <li>• Items with original tags attached</li>
            <li>• Items in original packaging</li>
            <li>• Non-personalized items</li>
          </ul>
        </div>
        <div>
          <h3 class="font-semibold text-red-600 mb-3">❌ Non-Returnable Items</h3>
          <ul class="text-gray-600 space-y-2 text-sm">
            <li>• Personalized or custom items</li>
            <li>• Perishable goods</li>
            <li>• Intimate apparel</li>
            <li>• Digital downloads</li>
            <li>• Gift cards</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Exchanges -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Exchanges</h2>
      <p class="text-gray-600 mb-4">
        Want a different size or color? We make exchanges easy! Follow the same return process and let us know what you'd like instead.
      </p>
      <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <p class="text-yellow-800 text-sm">
          <strong>Note:</strong> If there's a price difference, we'll charge or refund the difference. If the new item is out of stock, we'll issue a full refund and you can place a new order.
        </p>
      </div>
    </div>

    <!-- Damaged Items -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Damaged or Defective Items</h2>
      <p class="text-gray-600 mb-4">
        If your item arrives damaged or defective, we'll make it right immediately. Contact us within 7 days of delivery with photos of the damage.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="bg-green-50 border border-green-200 rounded-lg p-4">
          <h3 class="font-semibold text-green-900 mb-2">What We'll Do</h3>
          <ul class="text-green-800 text-sm space-y-1">
            <li>• Send a replacement immediately</li>
            <li>• Provide a full refund</li>
            <li>• Cover all return shipping costs</li>
          </ul>
        </div>
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 class="font-semibold text-blue-900 mb-2">What We Need</h3>
          <ul class="text-blue-800 text-sm space-y-1">
            <li>• Photos of the damaged item</li>
            <li>• Your order number</li>
            <li>• Description of the issue</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Contact Info -->
    <div class="bg-gray-50 rounded-lg p-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Need Help?</h2>
      <p class="text-gray-600 mb-4">
        Our customer service team is here to help with any questions about returns or exchanges.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="text-center">
          <div class="text-2xl mb-2">📧</div>
          <h3 class="font-semibold text-gray-900">Email</h3>
          <p class="text-sm text-gray-600">returns@shop.com</p>
        </div>
        <div class="text-center">
          <div class="text-2xl mb-2">📞</div>
          <h3 class="font-semibold text-gray-900">Phone</h3>
          <p class="text-sm text-gray-600">1-800-123-4567</p>
        </div>
        <div class="text-center">
          <div class="text-2xl mb-2">💬</div>
          <h3 class="font-semibold text-gray-900">Live Chat</h3>
          <p class="text-sm text-gray-600">Available 24/7</p>
        </div>
      </div>
    </div>
  </div>
</div>

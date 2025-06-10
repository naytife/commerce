import { writable, derived } from 'svelte/store';

export const currencyCode = writable<string>('USD');

export const currencySymbol = derived(currencyCode, ($code) => {
  switch ($code) {
    case 'USD':
      return '$';
    case 'NGN':
      return '₦';
    case 'EUR':
      return '€';
    case 'GBP':
      return '£';
    default:
      return $code;
  }
}); 
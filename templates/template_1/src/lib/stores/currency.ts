import { writable, derived } from 'svelte/store';

// Currency mapping - matches the dashboard currency utility
const CURRENCIES: Record<string, string> = {
  USD: '$', NGN: '₦', EUR: '€', GBP: '£', JPY: '¥', CAD: 'C$', AUD: 'A$',
  CHF: 'CHF', CNY: '¥', INR: '₹', ZAR: 'R', BRL: 'R$', MXN: '$', SGD: 'S$',
  HKD: 'HK$', SEK: 'kr', NOK: 'kr', DKK: 'kr', PLN: 'zł', CZK: 'Kč',
  HUF: 'Ft', TRY: '₺', RUB: '₽', KRW: '₩', THB: '฿', MYR: 'RM',
  IDR: 'Rp', PHP: '₱', VND: '₫', EGP: '£', MAD: 'DH', KES: 'KSh', GHS: '₵'
};

export const currencyCode = writable<string>('USD');

export const currencySymbol = derived(currencyCode, ($code) => {
  return CURRENCIES[$code?.toUpperCase()] || $code || '$';
});

/**
 * Format amount with currency symbol and proper locale formatting
 */
export function formatCurrencyWithLocale(amount: number, currencyCode: string, locale: string = 'en-US'): string {
  const symbol = CURRENCIES[currencyCode?.toUpperCase()] || currencyCode || '$';
  const formattedNumber = amount.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
  return `${symbol}${formattedNumber}`;
} 
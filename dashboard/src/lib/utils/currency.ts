/**
 * Currency utility functions for the dashboard
 * Maps currency codes to their respective symbols and provides formatting
 */

export interface CurrencyInfo {
  symbol: string;
  name: string;
  code: string;
}

export const CURRENCIES: Record<string, CurrencyInfo> = {
  USD: { symbol: '$', name: 'US Dollar', code: 'USD' },
  NGN: { symbol: '₦', name: 'Nigerian Naira', code: 'NGN' },
  EUR: { symbol: '€', name: 'Euro', code: 'EUR' },
  GBP: { symbol: '£', name: 'British Pound', code: 'GBP' },
  JPY: { symbol: '¥', name: 'Japanese Yen', code: 'JPY' },
  CAD: { symbol: 'C$', name: 'Canadian Dollar', code: 'CAD' },
  AUD: { symbol: 'A$', name: 'Australian Dollar', code: 'AUD' },
  CHF: { symbol: 'CHF', name: 'Swiss Franc', code: 'CHF' },
  CNY: { symbol: '¥', name: 'Chinese Yuan', code: 'CNY' },
  INR: { symbol: '₹', name: 'Indian Rupee', code: 'INR' },
  ZAR: { symbol: 'R', name: 'South African Rand', code: 'ZAR' },
  BRL: { symbol: 'R$', name: 'Brazilian Real', code: 'BRL' },
  MXN: { symbol: '$', name: 'Mexican Peso', code: 'MXN' },
  SGD: { symbol: 'S$', name: 'Singapore Dollar', code: 'SGD' },
  HKD: { symbol: 'HK$', name: 'Hong Kong Dollar', code: 'HKD' },
  SEK: { symbol: 'kr', name: 'Swedish Krona', code: 'SEK' },
  NOK: { symbol: 'kr', name: 'Norwegian Krone', code: 'NOK' },
  DKK: { symbol: 'kr', name: 'Danish Krone', code: 'DKK' },
  PLN: { symbol: 'zł', name: 'Polish Zloty', code: 'PLN' },
  CZK: { symbol: 'Kč', name: 'Czech Koruna', code: 'CZK' },
  HUF: { symbol: 'Ft', name: 'Hungarian Forint', code: 'HUF' },
  TRY: { symbol: '₺', name: 'Turkish Lira', code: 'TRY' },
  RUB: { symbol: '₽', name: 'Russian Ruble', code: 'RUB' },
  KRW: { symbol: '₩', name: 'South Korean Won', code: 'KRW' },
  THB: { symbol: '฿', name: 'Thai Baht', code: 'THB' },
  MYR: { symbol: 'RM', name: 'Malaysian Ringgit', code: 'MYR' },
  IDR: { symbol: 'Rp', name: 'Indonesian Rupiah', code: 'IDR' },
  PHP: { symbol: '₱', name: 'Philippine Peso', code: 'PHP' },
  VND: { symbol: '₫', name: 'Vietnamese Dong', code: 'VND' },
  EGP: { symbol: '£', name: 'Egyptian Pound', code: 'EGP' },
  MAD: { symbol: 'DH', name: 'Moroccan Dirham', code: 'MAD' },
  KES: { symbol: 'KSh', name: 'Kenyan Shilling', code: 'KES' },
  GHS: { symbol: '₵', name: 'Ghanaian Cedi', code: 'GHS' },
};

/**
 * Get currency symbol from currency code
 */
export function getCurrencySymbol(currencyCode: string): string {
  const currency = CURRENCIES[currencyCode?.toUpperCase()];
  return currency?.symbol || currencyCode || '$';
}

/**
 * Get currency info from currency code
 */
export function getCurrencyInfo(currencyCode: string): CurrencyInfo {
  const currency = CURRENCIES[currencyCode?.toUpperCase()];
  return currency || { symbol: '$', name: 'Unknown Currency', code: currencyCode || 'USD' };
}

/**
 * Format amount with currency symbol
 */
export function formatCurrency(amount: number, currencyCode: string): string {
  const symbol = getCurrencySymbol(currencyCode);
  return `${symbol}${amount.toFixed(2)}`;
}

/**
 * Format amount with currency symbol and proper locale formatting
 */
export function formatCurrencyWithLocale(amount: number, currencyCode: string, locale: string = 'en-US'): string {
  const symbol = getCurrencySymbol(currencyCode);
  const formattedNumber = amount.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
  return `${symbol}${formattedNumber}`;
}

/**
 * Format as currency for input fields (without symbol)
 */
export function formatAsCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'decimal',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(value);
}

/**
 * Parse currency input back to number
 */
export function parseCurrencyInput(value: string): number {
  // Remove any non-numeric characters except decimal point
  const numericValue = value.replace(/[^0-9.]/g, '');
  return parseFloat(numericValue) || 0;
}

/**
 * Get list of supported currencies for dropdowns
 */
export function getSupportedCurrencies(): CurrencyInfo[] {
  return Object.values(CURRENCIES).sort((a, b) => a.name.localeCompare(b.name));
}

/**
 * Get popular currencies (most commonly used)
 */
export function getPopularCurrencies(): CurrencyInfo[] {
  const popularCodes = ['USD', 'EUR', 'GBP', 'NGN', 'JPY', 'CAD', 'AUD', 'CHF', 'CNY', 'INR'];
  return popularCodes.map(code => CURRENCIES[code]).filter(Boolean);
}

import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { api } from '$lib/api';

export const actions: Actions = {
	default: async ({ request, locals }) => {
		const formData = await request.formData();
		const formId = formData.get('form-id')?.toString();
		
		try {
			// Process different forms based on form-id
			if (formId === 'general-settings-form') {
				// Process general settings form
				console.log('Processing general settings form');
				return {
					success: true,
					message: 'General settings updated successfully'
				};
			} else if (formId === 'seo-settings-form') {
				// Process SEO settings form
				console.log('Processing SEO settings form');
				return {
					success: true,
					message: 'SEO settings updated successfully'
				};
			} else if (formId === 'social-settings-form') {
				// Process social media settings form
				console.log('Processing social media settings form');
				return {
					success: true,
					message: 'Social media settings updated successfully'
				};
			} else if (formId === 'domain-settings-form') {
				// Process domain settings form
				console.log('Processing domain settings form');
				return {
					success: true,
					message: 'Domain settings updated successfully'
				};
			} else if (formId === 'payment-method-form') {
				// Process payment method form
				const methodId = formData.get('method-id')?.toString();
				
				// Extract all settings keys for this payment method
				const entries = Array.from(formData.entries());
				const settings = entries
					.filter(([key]) => key !== 'form-id' && key !== 'method-id')
					.reduce((acc, [key, value]) => {
						acc[key] = value;
						return acc;
					}, {} as Record<string, FormDataEntryValue>);
				
				console.log(`Processing payment method settings for: ${methodId}`, settings);
				
				return {
					success: true,
					message: 'Payment method settings updated successfully'
				};
			} else if (formId === 'billing-settings-form') {
				// Process billing settings form
				console.log('Processing billing settings form');
				// This would connect to payment processor or upgrade subscription
				// For now, just return success
				return {
					success: true,
					message: 'Billing settings updated successfully'
				};
			}
			
			return {
				success: false,
				message: 'Unknown form submitted'
			};
		} catch (error) {
			console.error('Failed to process form:', error);
			return fail(500, {
				success: false,
				message: 'An error occurred while processing your request'
			});
		}
	}
}; 
// schema.ts
import { z } from 'zod';

export const generalSettingsFormSchema = z.object({
    title: z
        .string()
        .min(2, 'title must be at least 2 characters.')
        .max(30, 'title must not be longer than 30 characters'),
    email: z.string({ required_error: 'Please enter an email address' }).email(),
    description: z.string().min(4).max(160).default('Our amazing shop.'),
    phone: z
        .string()
        .min(1, 'Please enter a valid phone number')
        .max(15, 'Please enter a valid phone number')
});

export type GeneralSettingsFormSchema = typeof generalSettingsFormSchema;

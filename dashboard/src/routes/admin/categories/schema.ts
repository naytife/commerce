// schema.ts
import { z } from 'zod';

export const formSchema = z.object({
	name: z.string().min(1, { message: 'Name is required' }),
	description: z.string().min(1, { message: 'Description is required' }),
	// categoryPath: z.string().nullable().optional(),
	// subcategory: z.string().nullable().optional(),
});

export type FormSchema = z.infer<typeof formSchema>;
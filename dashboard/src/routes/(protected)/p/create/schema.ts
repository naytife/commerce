// schema.ts
import { z } from 'zod';

export const formSchema = z.object({
	size: z.array(z.string()).default(['L', 'M', 'S'])
});

export type FormSchema = typeof formSchema;

import type { PageServerLoad, Actions } from "./$types.js";
// import { graphql } from "$houdini"; 
import { fail } from "@sveltejs/kit";
import { superValidate } from 'sveltekit-superforms';
import { generalSettingsFormSchema } from './schema';
import { zod } from 'sveltekit-superforms/adapters';

export const load: PageServerLoad = async () => {
    const form = await superValidate( zod(generalSettingsFormSchema));
    return { form };
};

export const actions: Actions = {
    default: async (event) => {
      const form = await superValidate(event, zod(generalSettingsFormSchema));
  
      if (!form.valid) {
        return fail(400, { form });
      }
      return {form}
    },
  };
import type { PageServerLoad, Actions } from "./$types.js";
import { graphql } from "$houdini"; // Import graphql from Houdini
import { fail } from "@sveltejs/kit";
import { superValidate } from 'sveltekit-superforms';
import { formSchema } from './schema';
import { zod } from 'sveltekit-superforms/adapters';

// Define your mutation using GraphQL
const CreateCategoryMutation = graphql(`
    mutation CreateCategory($parentID: ID!, $title: String!, $description: String) {
        createCategory(category: { parentID: $parentID, title: $title, description: $description }) {
            ...category
        }
    }
`); 

export const load: PageServerLoad = async () => {
    const form = await superValidate( zod(formSchema));
    return { form };
};

export const actions: Actions = {
    default: async (event) => {
      const form = await superValidate(event, zod(formSchema));
  
      if (!form.valid) {
        return fail(400, { form });
      }
      return {form}
    },
  };

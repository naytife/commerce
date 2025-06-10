import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent, fetch }) => {
  // Get session from the parent layout
  const { session } = await parent();

  return {  }; // No initial data needed for create page
};
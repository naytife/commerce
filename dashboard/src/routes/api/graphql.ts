// import { json } from '@sveltejs/kit';
// import { graphql, buildSchema } from 'graphql';
// import type { RequestHandler } from '@sveltejs/kit';
// import { readFileSync } from 'fs';
// import { join } from 'path';

// // Read the GraphQL query from the file
// const queryFilePath = join(process.cwd(), 'src', 'routes','[ID]', '+page.gql');
// const Products = readFileSync(queryFilePath, 'utf8');

// // Define your schema here (you can import it if it's in a separate file)
// const schema = buildSchema(`
//     // Your schema definition here (as you provided)
// `);

// // Define your resolvers for the schema
// const root = {
//     products: () => {
//         return [
//             {
//                 id: '1',
//                 title: 'Sample Product',
//                 description: 'Description of Sample Product',
//                 allowedAttributes: [],
//                 images: [{ url: 'image-url', altText: 'alt text' }],
//                 status: 'PUBLISHED',
//                 createdAt: new Date().toISOString(),
//                 updatedAt: new Date().toISOString(),
//             },
//         ];
//     },
// };

// // Handle the incoming requests
// export const POST: RequestHandler = async ({ request }) => {
//     const { query } = await request.json();

//     // If a query is not provided, use the default query from the file
//     const response = await graphql(schema, query || Products, root);
//     return json(response);
// };

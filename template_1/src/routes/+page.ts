export const prerender = true;
export const ssr = false;
export const csr = true;
import { graphql } from '$houdini'

const Shop = graphql`
    query Shop {
    shop {
        id
        title
        defaultDomain
        contactPhone
        contactEmail
        whatsAppNumber
        whatsAppLink
        facebookLink
        instagramLink
        currencyCode
        about
        shopProductsCategory
        seoDescription
        seoKeywords
        seoTitle
        products {
            totalCount
            edges {
                cursor
                node {
                    id
                    title
                    description
                    updatedAt
                    createdAt
                }
            }
        }
        categories {
            totalCount
            edges {
                cursor
                node {
                    id
                    slug
                    title
                    description
                    updatedAt
                    createdAt
                }
            }
        }
    }
}
`

export const _houdini_load = [Shop]
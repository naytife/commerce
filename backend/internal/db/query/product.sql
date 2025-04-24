-- name: GetProductsByType :many
SELECT 
    p.product_id,
    p.slug,
    p.title,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Aggregate attributes separately
    (
        SELECT COALESCE(
            json_agg(
                json_build_object(
                    'attribute_id', pa.attribute_id,
                    'title', a.title,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::json
        )
        FROM product_attribute_values pa
        LEFT JOIN attributes a ON a.attribute_id = pa.attribute_id
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    ) AS attributes,

    -- Aggregate variants separately
      (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', pv.product_variation_id,
                    'description', pv.description,
                    'price', pv.price,
                    'sku', pv.sku,
                    'available_quantity', pv.available_quantity,
                    'is_default', pv.is_default,
                    'attributes', (
                        SELECT COALESCE(
                            jsonb_agg(
                                jsonb_build_object(
                                    'attribute_id', pva.attribute_id,
                                    'title', a.title,
                                    'attribute_option_id', pva.attribute_option_id,
                                    'value', COALESCE(ao.value, pva.value)
                                )
                            ) FILTER (WHERE pva.attribute_id IS NOT NULL),
                            '[]'::jsonb
                        )
                        FROM product_variation_attribute_values pva
                        LEFT JOIN attributes a ON a.attribute_id = pva.attribute_id
                        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pva.attribute_option_id
                        WHERE pva.product_variation_id = pv.product_variation_id
                    )
                )
            ) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    )::jsonb AS variants,
    
    -- Product images
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', pi.product_image_id,
                    'url', pi.url,
                    'alt', pi.alt
                )
            ) FILTER (WHERE pi.product_image_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_images pi
        WHERE pi.product_id = p.product_id
    )::jsonb AS images


FROM products p
WHERE p.shop_id = $1 
AND p.product_type_id = $2
AND p.product_id > $3
ORDER BY p.product_id
LIMIT $4 
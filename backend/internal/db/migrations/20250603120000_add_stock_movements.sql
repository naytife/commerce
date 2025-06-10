-- Stock movements table for inventory tracking
CREATE TABLE stock_movements (
    movement_id BIGSERIAL PRIMARY KEY,
    product_variation_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    movement_type VARCHAR(50) NOT NULL, -- 'SALE', 'RESTOCK', 'ADJUSTMENT', 'RETURN'
    quantity_change INT NOT NULL, -- positive for increase, negative for decrease
    quantity_before INT NOT NULL,
    quantity_after INT NOT NULL,
    reference_id BIGINT, -- order_id for sales, adjustment_id for manual adjustments
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

-- SET RLS for stock_movements
ALTER TABLE stock_movements ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON stock_movements
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

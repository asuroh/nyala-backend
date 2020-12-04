-- Test Query 
  -- a. Display Customer List including calculating the total order.
    SELECT
      def.customer_id,
      def.customer_name,
      def.email,
      def.phone_number,
      to_char(dob, 'DD-MM-YYYY') as dob,
      CASE
              WHEN sex 
                  THEN 'Male'
              ELSE 'Famale'
        END sex,
        COUNT(o.order_id) as jumlah_order,
        SUM(od.qty*p.basic_price) as total_price_order
    FROM
      customers def
      LEFT JOIN orders o ON o.customer_id = def.customer_id
      LEFT JOIN order_details od on od.order_id = o.order_id
      LEFT JOIN products p on p.product_id = od.product_id
    GROUP BY def.customer_id
    
  -- b. Show Product List including calculating the number of orders sorted by the most in the order.
    SELECT
    pd.product_id,
    pd.product_name,
    pd.basic_price,
    pd.created_at,
    COUNT(o.order_id) as total_order
  FROM
    products pd
    LEFT JOIN order_details ord ON pd.product_id = ord.product_id
    LEFT JOIN orders o ON o.order_id = ord.order_id
  GROUP BY pd.product_id
  ORDER BY total_order DESC 

  -- c. Display the sort payment method data most frequently used by customers.
    SELECT
    pm.payment_method_id,
    pm.method_name,
    pm.code,
    pm.created_at
  FROM
    payment_methods pm
  WHERE
    (SELECT COUNT(order_id) FROM orders WHERE payment_method_id = pm.payment_method_id) > 0;
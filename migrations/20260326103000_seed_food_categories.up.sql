INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
VALUES
  (0, 'Fresh Vegetables', 'fresh_vegetables', 10, 1, 'leaf', 'fresh grocery category'),
  (0, 'Fresh Fruits', 'fresh_fruits', 20, 1, 'apple', 'fresh grocery category'),
  (0, 'Meat And Poultry', 'meat_poultry', 30, 1, 'drumstick', 'fresh grocery category'),
  (0, 'Aquatic Seafood', 'aquatic_seafood', 40, 1, 'fish', 'fresh grocery category'),
  (0, 'Eggs And Dairy', 'eggs_dairy', 50, 1, 'milk', 'fresh grocery category'),
  (0, 'Staple Grains Oils', 'staple_grains_oils', 60, 1, 'basket', 'fresh grocery category'),
  (0, 'Frozen Fast Food', 'frozen_fast_food', 70, 1, 'snowflake', 'fresh grocery category'),
  (0, 'Prepared Dishes', 'prepared_dishes', 80, 1, 'bento', 'fresh grocery category'),
  (0, 'Snacks And Bakery', 'snacks_bakery', 90, 1, 'cookie', 'fresh grocery category'),
  (0, 'Beverages', 'beverages', 100, 1, 'cup', 'fresh grocery category'),
  (0, 'Condiments', 'condiments', 110, 1, 'seasoning', 'fresh grocery category'),
  (0, 'Alcohol', 'alcohol', 120, 1, 'wine', 'fresh grocery category')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `parent_id` = VALUES(`parent_id`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `icon` = VALUES(`icon`),
  `remark` = VALUES(`remark`),
  `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Leafy Vegetables', 'leafy_vegetables', 11, 1, 'leaf', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_vegetables'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Root Vegetables', 'root_vegetables', 12, 1, 'carrot', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_vegetables'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Melon Vegetables', 'melon_vegetables', 13, 1, 'cucumber', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_vegetables'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Tropical Fruits', 'tropical_fruits', 21, 1, 'banana', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_fruits'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Berry Fruits', 'berry_fruits', 22, 1, 'strawberry', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_fruits'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Citrus Fruits', 'citrus_fruits', 23, 1, 'orange', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'fresh_fruits'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Pork', 'pork', 31, 1, 'pork', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'meat_poultry'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Beef And Lamb', 'beef_lamb', 32, 1, 'beef', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'meat_poultry'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Poultry', 'poultry', 33, 1, 'chicken', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'meat_poultry'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Fish', 'fish_products', 41, 1, 'fish', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'aquatic_seafood'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Shrimp Crab Shellfish', 'shrimp_crab_shellfish', 42, 1, 'shrimp', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'aquatic_seafood'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Ready Seafood', 'ready_seafood', 43, 1, 'seafood', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'aquatic_seafood'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Eggs', 'eggs', 51, 1, 'egg', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'eggs_dairy'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Milk And Yogurt', 'milk_yogurt', 52, 1, 'milk', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'eggs_dairy'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Cheese And Butter', 'cheese_butter', 53, 1, 'cheese', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'eggs_dairy'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Rice And Noodles', 'rice_noodles', 61, 1, 'rice', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'staple_grains_oils'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Flour And Baking', 'flour_baking', 62, 1, 'flour', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'staple_grains_oils'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Edible Oil', 'edible_oil', 63, 1, 'oil', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'staple_grains_oils'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Frozen Dumplings Buns', 'frozen_dumplings_buns', 71, 1, 'dumpling', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'frozen_fast_food'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Frozen Meatballs', 'frozen_meatballs', 72, 1, 'meatball', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'frozen_fast_food'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Instant Meals', 'instant_meals', 73, 1, 'noodle', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'frozen_fast_food'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Ready To Cook', 'ready_to_cook', 81, 1, 'pan', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'prepared_dishes'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Ready To Heat', 'ready_to_heat', 82, 1, 'microwave', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'prepared_dishes'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Salads Cold Dishes', 'salads_cold_dishes', 83, 1, 'salad', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'prepared_dishes'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Bread And Pastry', 'bread_pastry', 91, 1, 'bread', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'snacks_bakery'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Chips And Nuts', 'chips_nuts', 92, 1, 'chips', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'snacks_bakery'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Desserts', 'desserts', 93, 1, 'cake', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'snacks_bakery'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Water And Tea', 'water_tea', 101, 1, 'water', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'beverages'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Juice And Milk Drinks', 'juice_milk_drinks', 102, 1, 'juice', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'beverages'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Coffee And Soda', 'coffee_soda', 103, 1, 'coffee', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'beverages'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Sauces', 'sauces', 111, 1, 'sauce', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'condiments'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Spices', 'spices', 112, 1, 'pepper', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'condiments'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Hotpot Soup Base', 'hotpot_soup_base', 113, 1, 'hotpot', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'condiments'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Beer', 'beer', 121, 1, 'beer', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'alcohol'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Wine', 'wine', 122, 1, 'wine', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'alcohol'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

INSERT INTO `categories` (`parent_id`, `name`, `code`, `sort`, `status`, `icon`, `remark`)
SELECT c.`id`, 'Spirits', 'spirits', 123, 1, 'whisky', 'seeded sub category'
FROM `categories` c
WHERE c.`code` = 'alcohol'
ON DUPLICATE KEY UPDATE `parent_id` = VALUES(`parent_id`), `name` = VALUES(`name`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `icon` = VALUES(`icon`), `remark` = VALUES(`remark`), `deleted_at` = NULL;

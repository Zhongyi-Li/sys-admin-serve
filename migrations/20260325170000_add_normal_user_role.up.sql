INSERT INTO `roles` (`name`, `code`, `status`, `sort`, `remark`)
VALUES ('Normal User', 'normal_user', 1, 100, 'default role for registered users')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `status` = VALUES(`status`),
  `sort` = VALUES(`sort`),
  `remark` = VALUES(`remark`),
  `deleted_at` = NULL;

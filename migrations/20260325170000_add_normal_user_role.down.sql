DELETE ur
FROM `user_roles` ur
INNER JOIN `roles` r ON ur.role_id = r.id
WHERE r.code = 'normal_user';

DELETE rm
FROM `role_menus` rm
INNER JOIN `roles` r ON rm.role_id = r.id
WHERE r.code = 'normal_user';

DELETE FROM `roles`
WHERE `code` = 'normal_user';

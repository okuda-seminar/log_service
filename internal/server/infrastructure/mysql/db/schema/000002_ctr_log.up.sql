CREATE TABLE IF NOT EXISTS `ctr_logs` (
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created_At',
  `object_id` VARCHAR(100) NOT NULL COMMENT 'Object_ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
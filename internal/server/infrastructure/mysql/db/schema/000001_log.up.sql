CREATE TABLE IF NOT EXISTS `logs` (
  `log_level` VARCHAR(100) NOT NULL COMMENT 'Log_Level',
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Date',
  `destination_service` VARCHAR(100) NOT NULL COMMENT 'Destination_Service',
  `source_service` VARCHAR(100) NOT NULL COMMENT 'Source_Service',
  `request_type` VARCHAR(100) NOT NULL COMMENT 'Request_Type',
  `content` TEXT NOT NULL COMMENT 'Content'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
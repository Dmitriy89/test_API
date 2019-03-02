CREATE TABLE `testing` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`reqjson` LONGTEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

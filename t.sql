CREATE TABLE `tb_url` (
  `id` int(11) NOT NULL,
  `int_field` int(11) NOT NULL COMMENT '整型',
  `date_field` date NOT NULL,
  `datetime_filed` datetime NOT NULL,
  `timestamp_filed` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `text_field` text COLLATE utf8_bin NOT NULL,
  `double_filed` double NOT NULL,
  `tiny_int` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
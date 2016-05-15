CREATE DATABASE IF NOT EXISTS `user_health_information`;
\connect `user_health_information`

--
-- Table structure for users --
--

CREATE TABLE IF NOT EXISTS `user`(
  `id` INT PRIMARY KEY,
  `name` VARCHAR(256) NOT NULL,
  `email` VARCHAR(256) NOT NULL
  `username` VARCHAR(256) NOT NULL UNIQUE,
  UNIQUE,
);

CREATE INDEX IF NOT EXISTS `name_idx` ON user(`name`);


--
-- Table structure for table health_stats --
--

CREATE TABLE IF NOT EXISTS `health_stats` (
  `id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `step_count` INT NOT NULL,
  `flights_climbed` INT NOT NULL,
  `distance` FLOAT NOT NULL,
  `time` DATETIME NOT NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`user_id`) REFERENCES user(`id`),
);

CREATE INDEX IF NOT EXISTS `step_count_idx` ON health_stats(`step_count`);
CREATE INDEX IF NOT EXISTS `flights_climbed_idx` ON health_stats(`flights_climbed`);
CREATE INDEX IF NOT EXISTS `distance_idx` ON health_stats(`distance`);
CREATE INDEX IF NOT EXISTS `time_idx` ON health_stats(`time`);


--
-- Table structure for users --
--

CREATE TABLE users(
  id INT,
  name VARCHAR(256) NOT NULL,
  email VARCHAR(256) NOT NULL,
  username VARCHAR(256) NOT NULL,
  password VARCHAR(256) NOT NULL,
  PRIMARY KEY(id),
  UNIQUE(email),
  UNIQUE(username)
);
CREATE INDEX IF NOT EXISTS name_idx ON users(name);


--
-- Table structure for table health_stats --
--

CREATE TABLE health_stats(
  id INT NOT NULL,
  user_id INT NOT NULL,
  step_count INT NOT NULL,
  flights_climbed INT NOT NULL,
  distance FLOAT NOT NULL,
  time timestamp NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS step_count_idx ON health_stats(step_count);
CREATE INDEX IF NOT EXISTS flights_climbed_idx ON health_stats(flights_climbed);
CREATE INDEX IF NOT EXISTS distance_idx ON health_stats(distance);
CREATE INDEX IF NOT EXISTS time_idx ON health_stats(time);

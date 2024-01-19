CREATE TABLE users (
  id   INT          NOT NULL AUTO_INCREMENT,
  name VARCHAR(256) NOT NULL,
  PRIMARY KEY (id)
);

INSERT INTO users (name) VALUES ("alice");
INSERT INTO users (name) VALUES ("bob");

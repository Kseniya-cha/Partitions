CREATE TABLE global_cycles (
  id bigserial PRIMARY KEY,
  start_datetime timestamp NOT NULL DEFAULT (now()),
  end_datetime timestamp
);

CREATE TABLE results (
  id bigserial,
  cycles_id int NOT NULL,
  uuid int NOT NULL,
  start_datetime timestamp NOT NULL DEFAULT (now()),
	PRIMARY KEY (id, start_datetime)
) PARTITION BY RANGE (start_datetime);

CREATE INDEX ON global_cycles (start_datetime);

CREATE INDEX ON results (cycles_id);

ALTER TABLE results ADD FOREIGN KEY (cycles_id) REFERENCES global_cycles (id);
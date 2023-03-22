CREATE TABLE "monitoring_cycles" (
  id bigserial PRIMARY KEY,
  start_datetime timestamp NOT NULL DEFAULT (now()),
  end_datetime timestamp
);

CREATE TABLE "device_testing_results" (
  id bigserial,
  cycles_id int NOT NULL,
  start_datetime timestamp NOT NULL DEFAULT (now()),
	PRIMARY KEY (id, start_datetime)
) PARTITION BY RANGE (start_datetime);

CREATE INDEX ON monitoring_cycles (start_datetime);

CREATE INDEX ON device_testing_results (cycles_id);

ALTER TABLE device_testing_results ADD FOREIGN KEY (cycles_id) REFERENCES monitoring_cycles (id);

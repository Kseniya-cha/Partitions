CREATE IF NOT EXISTS TABLE my_table (
    id SERIAL,
    name VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE 
) PARTITION BY RANGE (created_at);

CREATE INDEX my_table_idx ON my_table (created_at);

CREATE TABLE my_table_def PARTITION OF my_table DEFAULT;

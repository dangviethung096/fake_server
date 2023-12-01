DROP TABLE IF EXISTS scheduler_tasks;

CREATE TABLE scheduler_tasks (
    id VARCHAR(20) PRIMARY KEY,
    queue_name TEXT,
    data BYTEA,
    done BOOLEAN,
    loop_index BIGINT,
    loop_count BIGINT,
    next BIGINT,
    interval BIGINT,
    start_time TEXT
);

DROP TABLE IF EXISTS scheduler_todo;

CREATE TABLE scheduler_todo (
    id SERIAL,
    task_id VARCHAR(20),
    bucket BIGINT
);

DROP TABLE IF EXISTS scheduler_done;

CREATE TABLE scheduler_done (
    id VARCHAR(20) PRIMARY KEY,
    bucket BIGINT
);
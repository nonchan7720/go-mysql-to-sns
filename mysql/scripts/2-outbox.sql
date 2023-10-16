CREATE TABLE outbox(
  id bigint auto_increment,
  aggregate_type varchar(255),
  aggregate_id varchar(255),
  event_type varchar(255),
  payload json,
  primary key(id)
)

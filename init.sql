CREATE TABLE IF NOT EXISTS users
(
    user_id   varchar(20) PRIMARY KEY,
    user_name varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS result
(
    result_id SERIAL PRIMARY KEY,
    user_id   varchar(20)    NOT NULL,
    max_speed numeric(10, 5) NOT NULL,
    distance  numeric(10, 5) NOT NULL,
    CONSTRAINT fk_user_result FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE VIEW rating(user_id, user_name, max_speed, max_distance, distance) AS
SELECT u.user_id, u.user_name, max(r.max_speed), max(r.distance), sum(r.distance)
FROM result r JOIN users u on u.user_id = r.user_id
GROUP BY u.user_id, u.user_name;
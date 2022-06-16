CREATE TABLE workouts
(
    id             SERIAL       NOT NULL UNIQUE,
    user_id        integer      NOT NULL,
    title          varchar(255) NOT NULL DEFAULT '',
    description    varchar(500) NOT NULL DEFAULT '',
    is_done        boolean      NOT NULL DEFAULT false,
    appointed_time timestamp    ,
    created_at     timestamp    NOT NULL DEFAULT now(),
    updated_at     timestamp    NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);


CREATE TABLE exercises
(
    id          SERIAL       NOT NULL UNIQUE,
    workout_id  integer      NOT NULL,
    title       varchar(255) NOT NULL default '',
    description varchar(500) NOT NULL default '',
    is_done     boolean      NOT NULL DEFAULT false
);


CREATE TABLE repetitions
(
    id          SERIAL  NOT NULL UNIQUE,
    exercise_id integer NOT NULL,
    weight      float   NOT NULL DEFAULT 0.0,
    count       integer NOT NULL DEFAULT 0,
    is_done     boolean NOT NULL DEFAULT false
);

CREATE TABLE exercises_attachments
(
    id          SERIAL       NOT NULL UNIQUE,
    exercise_id integer      NOT NULL,
    url         varchar(500) NOT NULL default ''
);


alter table exercises
    add constraint fk_exercise
        foreign key (workout_id)
            REFERENCES workouts (id)
            ON DELETE CASCADE;

alter table repetitions
    add constraint fk_repetition
        foreign key (exercise_id)
            REFERENCES exercises (id)
            ON DELETE CASCADE;

alter table exercises_attachments
    add constraint fk_exercise_attachment
        foreign key (exercise_id)
            REFERENCES exercises (id)
            ON DELETE CASCADE;


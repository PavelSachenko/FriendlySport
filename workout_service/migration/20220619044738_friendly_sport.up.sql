CREATE VIEW workouts_objects AS
select workouts.id,
       workouts.user_id,
       workouts.title,
       workouts.description,
       workouts.is_done,
       workouts.appointed_time,
       workouts.created_at,
       workouts.updated_at,
       (select array_to_json(array_agg(json_build_object('id', exercise.id, 'title', exercise.title, 'description',
                                                         exercise.description, 'is_done', exercise.is_done,
                                                         'repetitions',
                                                         (select array_to_json(array_agg(json_build_object('id',
                                                                                                           repetition.id,
                                                                                                           'weight',
                                                                                                           repetition.weight,
                                                                                                           'count',
                                                                                                           repetition.count,
                                                                                                           'is_done',
                                                                                                           repetition.is_done)))
                                                          from (select *
                                                                from repetitions
                                                                WHERE repetitions.exercise_id = exercise.id) as repetition))
           ))
        from (select *
              FROM exercises
              WHERE exercises.workout_id = workouts.id) as exercise) as exercises_json
from workouts;
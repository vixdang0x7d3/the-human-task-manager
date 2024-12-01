-- +goose Up
-- +goose StatementBegin
CREATE  OR REPLACE FUNCTION fn_lookup_coeff(arg_name TEXT)
returns NUMERIC
AS
$$
DECLARE
	ret NUMERIC := 0;
BEGIN
WITH lk_coeff(name, coefficient) AS (
	SELECT * FROM ( VALUES
	 	('tag_next', 15.0),
		('deadline', 12.0),
		('scheduled', 5.0),
		('wait', -3.0),
		('priority_h', 6.0),
		('priority_m', 3.9),
		('priority_l', -1.8),
		('age', 2.0),
		('tags', 1.0),
		('project', 1.0),
	        ('none', 0.0)
	) AS lk_coeff( name, coefficient ) )
SELECT lk_coeff.coefficient INTO ret FROM lk_coeff WHERE lk_coeff.name=arg_name;
RETURN(ret);
END;
$$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION fn_urgency(task_record TASKS)
returns NUMERIC
AS
$$
DECLARE
	urgency NUMERIC;
	time_zero timestamptz := '0001-01-01 00:00:00'::timestamptz;

	lk_project text;
	lk_priority text;
	lk_tag_next text;
	lk_status text;
	age	NUMERIC;
	scheduled NUMERIC;
	deadline NUMERIC;
	waited NUMERIC;
BEGIN
	SELECT CASE
		WHEN task_record.project_id IS NOT NULL THEN 'project'
		ELSE 'none'
	END INTO lk_project;

	SELECT CASE
		WHEN task_record.priority = 'H'::task_priority THEN 'priority_h'
		WHEN task_record.priority = 'M'::task_priority THEN 'priority_m'
		WHEN task_record.priority = 'L'::task_priority THEN 'priority_l'
		ELSE 'none'
	END INTO lk_priority;

	SELECT CASE
		WHEN ARRAY['next'] <@ task_record.tags THEN 'tag_next'
		ELSE 'none'
	END INTO lk_tag_next;

	SELECT current_timestamp::DATE - task_record.CREATE::DATE INTO age;

	WITH tr AS (
		SELECT task_record.deadline::DATE - current_timestamp::DATE AS days
	) SELECT CASE
		WHEN task_record.deadline = time_zero
			OR task_record.state = 'waiting'::task_state
			OR tr.days > 3
		THEN 0
		ELSE 1
	END INTO deadline
	FROM tr;

	WITH tr AS (
		SELECT task_record.schedule::DATE - current_timestamp::DATE AS days
	) SELECT CASE
		WHEN task_record.schedule = time_zero
			OR task_record.state = 'waiting'::task_state
			OR tr.days > 3
		THEN 0
		ELSE 1
	END INTO scheduled
	FROM tr;

	WITH tr AS (
		SELECT task_record.wait::DATE - current_timestamp::DATE AS days
	) SELECT CASE
		WHEN task_record.wait = time_zero
			OR tr.days <= 0
		THEN 0
		ELSE tr.days
	END INTO waited
	FROM tr;

	SELECT
		fn_lookup_coeff(lk_project) +
		fn_lookup_coeff(lk_priority) +
		fn_lookup_coeff(lk_tag_next) +
		fn_lookup_coeff('age') * age +
		fn_lookup_coeff('deadline') * deadline +
		fn_lookup_coeff('scheduled') * scheduled +
		fn_lookup_coeff('wait') * waited
	INTO urgency;


	/* -- raise notice 'project term: %', fn_lookup_coeff(lk_project); */
	/* -- raise notice 'priority term: %', fn_lookup_coeff(lk_priority); */
	/* -- raise notice 'tag_next term: %', fn_lookup_coeff(lk_tag_next); */
	/* -- raise notice 'age term: %', fn_lookup_coeff('age') * age; */
	/* -- raise notice 'deadline term: %', fn_lookup_coeff('deadline') * deadline; */
	/* -- raise notice 'scheduled term: %', fn_lookup_coeff('scheduled') * scheduled; */
	/* -- raise notice 'waited term: %', fn_lookup_coeff('wait') * waited; */
	/* -- raise notice '---------'; */

	RETURN(urgency);
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE VIEW task_items AS
SELECT
	t.id,
	t.user_id,
	coalesce(u.username, '') AS username,
	t.project_id,
	coalesce(p.title, '') AS project_title,
	t.completed_by AS completed_by,
	coalesce(c.username, '') AS completed_by_name,
	t.description,
	t.priority,
	t.state,
	t.deadline,
	t.schedule,
	t.wait,
	t.create,
	t.end,
	t.tags,
	fn_urgency(t.*) AS urgency
FROM tasks t
LEFT JOIN users u ON t.user_id = u.id
LEFT JOIN users c ON t.completed_by = c.id AND t.completed_by IS NOT NULL
LEFT JOIN projects p ON t.project_id = p.id AND t.project_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW task_items; 
DROP FUNCTION IF EXISTS fn_lookup_coeff(TEXT);
DROP FUNCTION IF EXISTS fn_urgency(TASKS);
-- +goose StatementEnd

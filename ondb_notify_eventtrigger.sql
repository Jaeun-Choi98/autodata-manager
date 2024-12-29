
-- < 데이터베이스 레벨(ddl)의 트리거에서는 트리거가 아닌 이벤트 트리거를 사용한다. > 

CREATE OR REPLACE FUNCTION notify_ddl_trigger()
RETURNS event_trigger AS $$
DECLARE
    notifications TEXT = '';
    row record;
BEGIN
    FOR row IN SELECT * FROM pg_event_trigger_ddl_commands() LOOP
        IF row.command_tag = 'CREATE TABLE' THEN
            notifications := notifications || 'Table created: ' || row.object_identity || E'\n';
        END IF;
    END LOOP;

    IF notifications != '' THEN
        PERFORM pg_notify('table_events', notifications);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- DROP by DDL COMMAND의 처리는 따로 설정해줘야함. pg_event_trigger_ddl_commands에 찍히지 않음.
-- 참고 문서(https://www.postgresql.org/docs/current/functions-event-triggers.html#PG-EVENT-TRIGGER-DDL-COMMAND-END-FUNCTIONS)
CREATE EVENT TRIGGER table_ddl_trigger
ON ddl_command_end
EXECUTE FUNCTION notify_ddl_trigger();

CREATE OR REPLACE FUNCTION notify_dropped_trigger()
RETURNS event_trigger AS $$
DECLARE
    notifications TEXT = '';
    obj record;
BEGIN
    FOR obj IN SELECT * FROM pg_event_trigger_dropped_objects() LOOP
        IF obj.object_type = 'table' THEN
            notifications := notifications || 'Table deleted: ' || obj.object_identity || E'\n';
        END IF;
    END LOOP;

    IF notifications != '' THEN
        PERFORM pg_notify('table_events', notifications);
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE EVENT TRIGGER table_dropped_trigger
ON sql_drop
EXECUTE FUNCTION notify_dropped_trigger();
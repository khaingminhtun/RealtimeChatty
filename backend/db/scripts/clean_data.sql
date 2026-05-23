DO $$
DECLARE
  tbl text;
BEGIN
  FOR tbl IN
    SELECT tablename
    FROM pg_tables
    WHERE schemaname = 'public' AND tablename != 'goose_db_version'
  LOOP
    EXECUTE 'TRUNCATE TABLE ' || quote_ident(tbl) || ' RESTART IDENTITY CASCADE';
  END LOOP;
END $$;

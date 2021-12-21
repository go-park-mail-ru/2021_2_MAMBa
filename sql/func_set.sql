Create or replace FUNCTION  CountFilm() RETURNS bigint as $$
    begin
    Return(SELECT reltuples AS estimate FROM pg_class where relname = 'film');
    end;
    $$
    language plpgsql;

vacuum

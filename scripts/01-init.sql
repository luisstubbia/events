-- scripts/01-init.sql

-- Crear usuario y base de datos
CREATE USER eventdbuser WITH PASSWORD 'password';
CREATE DATABASE eventsdb OWNER eventdbuser;

-- Conectarse a la base de datos
\c eventsdb;

-- Otorgar todos los permisos necesarios
GRANT ALL PRIVILEGES ON DATABASE eventsdb TO eventdbuser;

-- Permisos sobre el esquema public
GRANT USAGE, CREATE ON SCHEMA public TO eventdbuser;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO eventdbuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO eventdbuser;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO eventdbuser;

-- Permisos por defecto para objetos futuros
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO eventdbuser;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO eventdbuser;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO eventdbuser;

-- Hacer al usuario propietario del esquema (opcional)
ALTER SCHEMA public OWNER TO eventdbuser;
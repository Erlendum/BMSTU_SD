FROM postgres:latest
RUN apt-get update
ADD sql_init/01_create.sql /docker-entrypoint-initdb.d
ADD sql_init/02_constraints.sql /docker-entrypoint-initdb.d
ADD sql_init/03_copy.sql /docker-entrypoint-initdb.d
RUN chmod a+r /docker-entrypoint-initdb.d/*
EXPOSE 6666

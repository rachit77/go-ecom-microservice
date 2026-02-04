FROM postgres:16.2

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
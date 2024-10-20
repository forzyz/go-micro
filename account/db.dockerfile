FROM postgres:16.4

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
FROM haproxy:2.3

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

EXPOSE 15432
EXPOSE 15433
EXPOSE 15434
EXPOSE 7000
FROM haproxy:2.3

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

EXPOSE 15432
EXPOSE 25432
EXPOSE 35432
EXPOSE 7000
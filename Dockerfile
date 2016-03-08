FROM ubuntu:14.04

COPY bonbon.conf /tmp/
COPY bonbon-server /tmp/

EXPOSE 8080
CMD /tmp/bonbon-server -config /tmp/bonbon.conf

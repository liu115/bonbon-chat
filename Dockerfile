FROM ubuntu:14.04

RUN apt-get update  
RUN apt-get install -y ca-certificates

COPY static/ /tmp/static/

RUN ls /tmp/static

COPY bonbon.conf /tmp/
COPY bonbon-server /tmp/

EXPOSE 8080
CMD /tmp/bonbon-server -config /tmp/bonbon.conf -static /tmp/static

FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN useradd -rm -d /home/coap -s /bin/bash -u 1002 coap

COPY --chown=coap:coap ./server/coap-server /home/coap/coap-server

USER coap
WORKDIR /home/coap

EXPOSE 5688/udp

CMD [ "/home/coap/coap-server" ]
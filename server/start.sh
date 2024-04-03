#!/bin/bash

FILE=/run/secrets/STORAGE_CONNECTION_STRING     
if [ -f $FILE ]; then
    export STORAGE_CONNECTION_STRING=$(cat /run/secrets/STORAGE_CONNECTION_STRING)
fi 

nohup /home/coap/coap-server -network udp4 -password "connect:anything" -dTLS &
nohup /home/coap/coap-server -network udp4 &
nohup /home/coap/coap-server -network udp6 -password "connect:anything" -dTLS &
nohup /home/coap/coap-server -network udp6 &

sleep 1

tail -f nohup.out
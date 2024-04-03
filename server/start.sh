#!/bin/bash

FILE=/run/secrets/STORAGE_CONNECTION_STRING     
if [ -f $FILE ]; then
    export STORAGE_CONNECTION_STRING=$(cat /run/secrets/STORAGE_CONNECTION_STRING)
fi 

nohup /home/coap/coap-server -password "connect:anything" -dTLS &
nohup /home/coap/coap-server &

sleep 1

tail -f nohup.out
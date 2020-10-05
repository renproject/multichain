#!/bin/bash

service rabbitmq-server start

tbears genconf
tbears -v start

tail -f /dev/null
#!/bin/bash

YANG_FILE=lagopus-switch.yang

MODULES=./modules

openconfigd -2 -z --yang-paths=${MODULES} ${YANG_FILE}

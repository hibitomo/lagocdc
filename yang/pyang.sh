#!/bin/bash

YANG_FILE=$@
if [ "$YANG_FILE" = "" ] ; then
    YANG_FILE="lagopus-switch.yang"
fi

MODULES=./modules

FORMAT=tree

pyang --format=${FORMAT} --path="${MODULES}" $YANG_FILE

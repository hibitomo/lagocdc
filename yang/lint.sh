#!/bin/bash

path=`dirname $0`
cd ${path}

YANG_FILE=$@

MODULES=./modules

while getopts a OPT
do
    case $OPT in
	"a" )
	    YANG_FILE=`find . -name "*.yang" -printf "%p "` ;;
	* )
	    exit 1 ;;
    esac
done

if [ "$YANG_FILE" = "" ] ; then
    YANG_FILE="lagopus-switch.yang"
fi

pyang --strict --lint --path="${MODULES}"  $YANG_FILE

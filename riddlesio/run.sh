#!/bin/bash

go build .

BASEDIR=`pwd`

java -jar $BASEDIR/match-wrapper-1.3.2.jar "$(cat wrapper-commands.json)"
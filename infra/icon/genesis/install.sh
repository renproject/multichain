#!/bin/bash

USE_PYPI=true

if [ $USE_PYPI == true ]; then
    pip3 install tbears
else
    WORKDIR=$(dirname $0)
    FILES=$(echo *.whl)
    for f in $FILES; do
        pip3 install $WORKDIR/$f
    done
fi
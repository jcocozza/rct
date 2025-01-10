#!/bin/bash

USER=$1

rct gen-config local > .rct.json
rct gen-config local > remote.rct.json

remotes=( "10.0.0.1" "10.0.0.2" "10.0.0.3" )

for item in "${remotes[@]}" do
    scp remote.rct.json "$USER@$item:.rct.json"
done

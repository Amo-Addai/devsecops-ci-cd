#!/bin/bash
# Activate virtual environment
. /appenv/bin/activate

## DOWNLOAD DEPENDENCIES FOR BUILDING

# Run entrypoint.sh arguments  eg. "npm start"
exec $@
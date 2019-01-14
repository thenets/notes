#!/bin/sh

# Start Redis server
redis-server &
#sleep 1

# Enable virtualenv
. /venv/bin/activate

# Start flask server
cd /app/src
python main.py
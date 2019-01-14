#!/bin/sh

set -e

# Start Redis server
cd /tmp
redis-server &
#sleep 1

# Enable virtualenv
. /venv/bin/activate

# Start flask server
cd /app
python main.py
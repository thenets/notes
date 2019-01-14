FROM alpine

RUN set -e \
    # Install general dev tools
    && apk add build-base python3-dev gcc alpine-sdk \
    # Install python3 and virtualenv
    && apk add py3-virtualenv \
    # Install Redis
    && apk add redis

# Install Python dependencies
ADD ./requirements.txt /app/
RUN set -x \
    # Create virtualenv
    && virtualenv -p python3 /venv \
    # Install libs
    && set +x \
    && . /venv/bin/activate \
    && set -x \
    && pip install -r /app/requirements.txt

ADD ./src /app
ADD ./*.sh /

RUN set -x \
    # Fix permissions
    && chmod +x /*.sh

WORKDIR /app

CMD [ "/entrypoint.sh" ]
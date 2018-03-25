FROM golang:1.8

RUN mkdir /src
WORKDIR /src

RUN mkdir /config
ADD ./config /config/

ADD ./src /src

RUN chmod +x /config/entrypoint.sh
RUN chmod +rx MatchingEngine.go
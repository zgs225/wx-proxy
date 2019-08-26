FROM golang:1.12

COPY . /app
RUN cd /app && go install -v && rm -rf /app

ENTRYPOINT ["wxproxy"]

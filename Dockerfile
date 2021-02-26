FROM golang:1.14
ADD . /app/
WORKDIR /app/src
RUN go get && go build -o falcon && chmod +x /app/src/falcon
CMD [ "./falcon" ]
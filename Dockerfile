FROM golang:latest
RUN mkdir app
ADD . /app
WORKDIR /app
RUN make
EXPOSE 8080
CMD ["/app/apiserver"]

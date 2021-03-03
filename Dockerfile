FROM node:12.18 AS front
WORKDIR /home/mo2front
COPY ./mo2front .
RUN npm install
RUN npm run build


#build stage
FROM golang:1.15 AS builder
RUN apt-get install git
WORKDIR /go/src/app
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...



#final stage
FROM ubuntu:latest
RUN apt-get update
RUN apt-get install ca-certificates -y
# RUN apt-get install graphviz -y
# RUN apt-get install wget -y
# RUN wget -O go.tar.gz https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz && \
#     tar -C /usr/local -xzf go.tar.gz && \
#     rm go.tar.gz
# ENV PATH=$PATH:/usr/local/go/bin
# RUN go version
COPY --from=front /home/dist /app/dist
WORKDIR /app
RUN chmod -R 777 .
COPY --from=builder /go/bin /app
ENV GIN_MODE=release
ENV MO2_MONGO_URL=mongodb://mongodb:27017
ENTRYPOINT /app/mo2
LABEL Name=mo2 Version=0.0.1
EXPOSE 5001

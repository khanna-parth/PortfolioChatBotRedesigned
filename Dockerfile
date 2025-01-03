# Dockerfile for compiling backend (go server, python internals)


FROM python:3.11-slim as python-builder

RUN apt-get update && apt-get install -y python3 python3-pip

WORKDIR /project/backend/internal/chain

COPY backend/internal/chain/requirements.txt /project/backend/internal/chain/requirements.txt
COPY backend/internal/chain/generator.py /project/backend/internal/chain/generator.py

RUN pip3 install --no-cache-dir -r /project/backend/internal/chain/requirements.txt


FROM golang:1.23-alpine as go-builder

RUN apk add --no-cache libc6-compat
WORKDIR /project/backend/server
COPY backend/server/ /project/backend/server
RUN go build -o server .


FROM python:3.11-slim

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y python3-pip && \
    ln -sf /usr/bin/python3 /usr/bin/python

RUN which python3 || echo "python3 not found" && \
    which python || echo "python not found" && \
    python3 --version || echo "python3 execution failed" && \
    python --version || echo "python execution failed"

COPY --from=python-builder /project/backend/internal/chain/requirements.txt /project/backend/internal/chain/requirements.txt
RUN pip3 install --no-cache-dir -r /project/backend/internal/chain/requirements.txt

COPY --from=python-builder /project/backend/ /project/backend/

COPY backend/server/prod.env /project/backend/server/prod.env

COPY --from=go-builder /project/backend/server/server /project/backend/server/server

RUN mkdir -p /project/backend/uploaded

WORKDIR /project/backend/server

EXPOSE 8080

CMD ["./server"]

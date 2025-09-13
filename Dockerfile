FROM python:3.11-slim as python-builder

WORKDIR /app/python

COPY backend/internal/chain/requirements.txt .

RUN pip install --no-cache-dir --extra-index-url https://download.pytorch.org/whl/cpu -r requirements.txt

COPY backend/internal/chain/ .

FROM golang:1.23-alpine as go-builder

RUN apk add --no-cache libc6-compat

WORKDIR /app/go

COPY backend/server/ .
RUN go build -o server .

FROM python:3.11-slim

WORKDIR /project

COPY --from=python-builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin

COPY --from=python-builder /app/python /project/backend/internal/chain

COPY --from=go-builder /app/go/server /project/backend/server/server

COPY backend/server/prod.env /project/backend/server/prod.env

RUN mkdir -p /project/backend/uploaded/STATIC/@Parth

WORKDIR /project/backend/server

EXPOSE 8080

CMD ["./server"]


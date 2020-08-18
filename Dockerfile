FROM golang:latest as builder

WORKDIR /usr/src/
COPY main.go go.mod ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo


FROM alpine:latest

ENV HOME /app
ENV PORT 6666
WORKDIR /app
COPY --from=builder /usr/src/sudoku_solver_api_server /app/
CMD /app/sudoku_solver_api_server

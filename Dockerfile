FROM golang:1-alpine as builder

# Build dependencies
RUN apk add --no-cache gcc make git

# Add all the code
ADD cmd/ /go/src/github.com/MathWebSearch/mwsapi/cmd
ADD elasticutils/ /go/src/github.com/MathWebSearch/mwsapi/elasticutils
ADD tema/ /go/src/github.com/MathWebSearch/mwsapi/tema

# And the makefile
ADD Makefile /go/src/github.com/MathWebSearch/mwsapi/

# Run make
WORKDIR /go/src/github.com/MathWebSearch/mwsapi
RUN make all

# and create a new image form scratch
FROM scratch

# copy all the built images
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/out/mwsapid /mwsapid
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/out/mwsquery /mwsquery
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/out/elasticsync /elasticsync

ENTRYPOINT [ "/mwsapid" ]
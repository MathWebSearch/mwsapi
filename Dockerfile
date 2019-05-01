FROM golang as builder

# copy everything into the $GOPATH
WORKDIR $GOPATH/src/github.com/MathWebSearch/mwsapi/
COPY . .

# And run make
RUN make all

# and create a new image form scratch
FROM scratch

# copy all the built images
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/mwsapid /mwsapid
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/mwsquery /mwsquery
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/elasticquery /elasticquery
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/elasticsync /elasticsync
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/temaquery /temaquery

ENTRYPOINT [ "/mwsapid" ]
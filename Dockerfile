FROM golang as builder

ENV MWSAPID_HOST 0.0.0.0
ENV MWSAPID_PORT 3000

ENV MWSAPID_MWS_HOST ""
ENV MWSAPID_MWS_PORT 8080

ENV MWSAPID_ELASTIC_HOST ""
ENV MWSAPID_ELASTIC_PORT 9200

# copy everything into the $GOPATH
WORKDIR $GOPATH/src/github.com/MathWebSearch/mwsapi/
COPY . .

# And run make
RUN make mwsapid

# and create a new image form scratch
FROM scratch

# copy all the built images
COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/mwsapid /mwsapid
#COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/mwsquery /mwsquery
#COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/elasticquery /elasticquery
#COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/elasticsync /elasticsync
#COPY --from=builder /go/src/github.com/MathWebSearch/mwsapi/temaquery /temaquery

EXPOSE 3000
ENTRYPOINT [ "/mwsapid" ]
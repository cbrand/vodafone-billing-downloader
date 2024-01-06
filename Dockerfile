FROM alpine:3.19

LABEL org.opencontainers.image.source https://github.com/cbrand/vodafone-billing-downloader

ENTRYPOINT ["/vodafone-billing-downloader"]
USER 1000

COPY vodafone-billing-downloader /

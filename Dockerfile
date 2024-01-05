FROM scratch

ENTRYPOINT ["/vodafone-billing-downloader"]
USER 1000

COPY vodafone-billing-downloader /

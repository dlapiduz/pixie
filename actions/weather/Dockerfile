FROM alpine

RUN apk add --update jq curl && rm -rf /var/cache/apk/*

ADD run.sh /scripts/run.sh

ENTRYPOINT ["/scripts/run.sh"]
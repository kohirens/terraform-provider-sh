FROM hashicorp/terraform:1.10.5 AS dev

ARG REPO
ARG USER_NAME
RUN apk update \
 && apk add aws-cli go

ENV GOPATH=/root/
ENV CGO_ENABLED=0

WORKDIR /${USER_NAME}/src/${REPO}

ENTRYPOINT [ "start.sh" ]

CMD []

COPY --chmod=555 .docker/start.sh /usr/local/bin/
COPY --chmod=640 .docker/terraform/.terraformrc /root/
COPY --chmod=440 go.mod ./
COPY --chmod=440 go.sum ./

RUN go mod tidy
FROM alpine:3.12

LABEL maintainer="cloud.jung@acornsoft.io"

ENV ANSIBLE_VERSION 2.9.7
ENV HELM_VERSION 3.6.3
ENV KUBECTL_VERSION 1.19.10

RUN	apk --update --no-cache add \
    ca-certificates \
    bash \
    git \
    openssh-client \
    openssl \
    python3 \
    py3-pip \
    py3-cryptography \
    rsync \
    sshpass \
    sudo \
    && apk --update add --virtual \
    .build-deps \
    python3-dev \
    libffi-dev \
    openssl-dev \
    build-base \
	&& pip3 install --upgrade \
		pip \
		cffi \
	&& pip3 install \
		ansible==${ANSIBLE_VERSION} \
		netaddr \
	&& apk del .build-deps \
	&& rm -rf /var/cache/apk/*

RUN	mkdir -p /etc/ansible && \
	echo 'localhost' > /etc/ansible/hosts && \
	mkdir -p ~/.ssh && touch ~/.ssh/known_hosts && \
	mkdir -p /knit/log

RUN wget -q -O /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && chmod +x /usr/bin/kubectl
RUN wget -q https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxvf /helm-v${HELM_VERSION}-linux-amd64.tar.gz \
    && cp /linux-amd64/helm /usr/bin/

ENV PATH="/knit:${PATH}"

# add openssh and clean
RUN apk add --update openssh \
    && rm  -rf /tmp/* /var/cache/apk/*

#make sure we get fresh keys
RUN rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key

VOLUME /knit

FROM regi.k3.acornsoft.io/k3lab/knit-base:1.0

ENV ANSIBLE_CONFIG=/koreon/scripts/ansible.cfg

COPY scripts /koreon/scripts
COPY inventory /koreon/inventory
COPY entrypoint.sh /koreon/entrypoint.sh

RUN chmod +x /koreon/entrypoint.sh
RUN rm -rf /sbin/apk

WORKDIR /koreon/work
ENTRYPOINT ["/koreon/entrypoint.sh"]
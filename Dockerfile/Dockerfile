FROM regi.k3.acornsoft.io/k3lab/knit-base:1.0

ENV ANSIBLE_CONFIG=/knit/scripts/ansible.cfg

COPY scripts /knit/scripts
COPY inventory /knit/inventory
COPY entrypoint.sh /knit/entrypoint.sh

RUN chmod +x /knit/entrypoint.sh
RUN rm -rf /sbin/apk

WORKDIR /knit/work
ENTRYPOINT ["/knit/entrypoint.sh"]
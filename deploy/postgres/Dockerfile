FROM postgres:18

# https://hub.docker.com/_/postgres/#locale-customization
RUN localedef -i et_EE -c -f UTF-8 -A /usr/share/locale/locale.alias et_EE.UTF-8
ENV LANG=et_EE.UTF-8

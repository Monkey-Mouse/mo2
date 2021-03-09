FROM nginx:1.18-alpine

COPY ./conf/nginx.conf /etc/nginx/nginx.conf

COPY ./script/entrypoint.sh /entrypoint.sh
COPY ./script/le.sh /le.sh

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN rm /etc/nginx/conf.d/default.conf 

RUN chmod +x /entrypoint.sh
RUN chmod +x /le.sh
RUN apk add --no-cache --update certbot tzdata openssl


ENTRYPOINT [ "/entrypoint.sh" ]

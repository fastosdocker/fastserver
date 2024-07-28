# 编译容器
FROM alpine:3.18
LABEL MAINTAINER=dockernb.com
WORKDIR /fast
EXPOSE 8081


RUN date
#健康检测
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
CMD pgrep dockercurl || exit 1
#后端安装
RUN wget http://192.168.0.170:9000/fastosdocker/fastserver/$(uname -s)-$(uname -m)/dockercurl -P /fast/
RUN chmod +x /fast/dockercurl
#前端安装
RUN mkdir /fast/static/ && wget http://192.168.0.170:9000/fastosdocker/web/pc.zip && unzip pc.zip && \
mv -f pc /fast/static/pc && rm /fast/pc.zip

RUN wget -O /fast/static/index.html http://192.168.0.170:9000/fastosdocker/web/index.html 


RUN sed -i 's@https://dl-cdn.alpinelinux.org/alpine@http://192.168.0.227:3142/alpine@g' /etc/apk/repositories && \
    apk add --update --no-cache curl git docker-compose
#时间时区
RUN apk -U add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && date \
    && apk del tzdata
#RUN curl -SL https://github.com/docker/compose/releases/download/v2.24.6/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
RUN chmod +x  /usr/bin/docker-compose && mkdir data
CMD ["/fast/dockercurl"] 

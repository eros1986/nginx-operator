FROM phusion/baseimage:0.9.22

RUN apt-get update
RUN apt-get install --reinstall tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN apt-get install -y supervisor

ENV APP_DIR /app

RUN mkdir -p ${APP_DIR}

ADD ./operator ${APP_DIR}
ADD ./run_operator.sh ${APP_DIR}
ADD ./operator.conf /etc/supervisor/conf.d/

RUN mkdir /etc/service/app
COPY ./run.sh /etc/service/app/run
RUN chmod +x /etc/service/app/run

CMD ["/sbin/my_init"]
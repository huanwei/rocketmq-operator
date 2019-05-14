FROM rocketmqinc/rocketmq:4.3.0

USER root
COPY brokerGenConfig.sh brokerStart.sh ${ROCKETMQ_HOME}/bin/
RUN chmod +x ${ROCKETMQ_HOME}/bin/brokerGenConfig.sh \
 && chmod +x ${ROCKETMQ_HOME}/bin/brokerStart.sh
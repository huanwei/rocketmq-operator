FROM apache/rocketmq-base:4.3.0

EXPOSE 10909 10911

RUN mv ${ROCKETMQ_HOME}/bin/runbroker-customize.sh ${ROCKETMQ_HOME}/bin/runbroker.sh \
 && chmod +x ${ROCKETMQ_HOME}/bin/runbroker.sh \
 && chmod +x ${ROCKETMQ_HOME}/bin/mqbroker

COPY brokerGenConfig.sh brokerStart.sh ${ROCKETMQ_HOME}/bin/
RUN chmod +x ${ROCKETMQ_HOME}/bin/brokerGenConfig.sh \
 && chmod +x ${ROCKETMQ_HOME}/bin/brokerStart.sh

WORKDIR ${ROCKETMQ_HOME}/bin/

#CMD cd ${ROCKETMQ_HOME}/bin \
# && export JAVA_OPT=" -Duser.home=/opt" \
# && sh mqbroker -n localhost:9876
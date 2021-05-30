#topic produce: test -> topic response: ("%v.response.%v", topic, NODE_ID) => test.response.123212321
#MSG
export MSG_TEMPLATE="TECHX"
#NODE
export NODE_ID="1"     #`date '+%s'`
export APPLICATION_NAME="motor"

#MYSQL
export MYSQL_USERNAME="mo"
export MYSQL_PASSWORD="123qwe"
export MYSQL_DB="tradex-order"
export MYSQL_HOST="localhost"
export MYSQL_PORT="3306"

#REDIS
export REDIS_HOST="localhost"
export REDIS_PORT="6379"

#KAFKA
export KAFKA_SERVERS="localhost:9092"
export KAFKA_MAIN_TOPIC="order"
export KAFKA_GROUP_ID="conditional-order"
export KAFKA_REQUEST_TIMEOUT="30"  #second
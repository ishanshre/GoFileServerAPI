createMongoDBcontainer:
	docker run -d -p 27017:27017 --name=GoFileServerAPI -v mongo_data:/data/db mongo

createRedisContainer:
	docker run -d --name GoFileServertRedis -p 6379:6379 redis:latest

createRedisInsightContainer:
	docker run -d --name redis-insight -p 8001:8001 redislabs/redisinsight:latest

startContainer:
	docker start GoFileServerAPI GoFileServertRedis



stopContainer:
	docker stop GoFileServerAPI GoFileServertRedis

run:
	go run cmd/api/main.go
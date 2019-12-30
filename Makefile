.PHONY: docker

GIT_VER = $(shell git describe --tags )

help:
	# try something else
	#
	# available targets
	# 
	#   docker        build docker iamge (version $(GIT_VER))
	#
docker:
	docker build -t paste:latest .
	docker tag paste:latest paste:$(GIT_VER)

docker-push:
	# Push latest
	docker tag paste:latest sigmonsays/paste:latest
	docker tag paste:latest sigmonsays/paste:$(GIT_VER)

	# Push version
	docker push sigmonsays/paste
	docker push sigmonsays/paste:$(GIT_VER)

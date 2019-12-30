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

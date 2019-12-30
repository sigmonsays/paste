.PHONY: docker

help:
	# try something else
	#
	# available targets
	# 
	#   docker        build docker iamge
	#
docker:
	docker build -t paste .

all:
	docker build -t timberslide/twitterstreamer .

push:
	docker push timberslide/twitterstreamer

clean:
	docker rmi -f timberslide/twitterstreamer

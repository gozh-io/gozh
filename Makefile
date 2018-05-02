IMAGE=gozh
TAG=1.0
REGISTRY=cloud.docker.com
REGISTRY_USER=blade2iron
PUBLISH_REGISTRY=cloud.docker.com
VENDOR=blade2iron
FQIN=$(VENDOR)/$(IMAGE)

all: build publish publish-test

build:
	docker build -t $(FQIN) .

rmi:
	docker rmi -f $(PUBLISH_REGISTRY)/$(FQIN):latest
	docker rmi -f $(PUBLISH_REGISTRY)/$(FQIN):$(TAG)
	docker rmi -f $(FQIN)
	docker rmi -f $(REGISTRY)/$(FQIN):latest
	docker rmi -f $(REGISTRY)/$(FQIN):$(TAG)

install:
	echo "TBD."
	
test:
	# sudo docker search $(REGISTRY)/$(FQIN)
	echo "TBD."

clean:
	# sudo atomic uninstall --force $(TEST_IMAGE)
	echo "TBD."

publish:
	docker login -u ${REGISTRY_USER} $(PUBLISH_REGISTRY)
	docker tag -f $(FQIN) $(PUBLISH_REGISTRY)/$(FQIN):latest
	docker push $(PUBLISH_REGISTRY)/$(FQIN):latest
	docker tag -f $(FQIN) $(PUBLISH_REGISTRY)/$(FQIN):$(TAG)
	docker push $(PUBLISH_REGISTRY)/$(FQIN):$(TAG)
	echo "RUN exmaple:"
	echo "docker run -p 80:80  $(REGISTRY)/$(FQIN):latest"
	echo "docker run --rm -it -p 80:80 $(REGISTRY)/$(FQIN):latest bash"

publish-test:
	docker search $(REGISTRY)/$(FQIN)
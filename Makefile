IMAGE=$(shell basename $(shell pwd))
#git rev-parse --show-toplevel
APP_NAME=$(IMAGE)
GIT_URL=$(shell git remote -v | grep '^origin\s.*(fetch)$$' | awk '{print $$2}' | sed 's/\.git//' | sed 's/http[s:/]*//')
TAG=1.0
REGISTRY=registry.hub.docker.com
#账号名称
REGISTRY_USER=blade2iron

FQIN=$(REGISTRY_USER)/$(IMAGE)

build:
	# echo "APP_NAME=$(APP_NAME) GIT_URL=$(GIT_URL)"
	docker build -t $(FQIN) --build-arg PROJECT_NAME=$(APP_NAME) --build-arg PROJECT_URL=$(GIT_URL) .

rmi:
	#docker rmi -f $(FQIN)
	#docker rmi -f $(REGISTRY)/$(FQIN):latest
	#docker rmi -f $(REGISTRY)/$(FQIN):$(TAG)
	docker rmi -f $(FQIN)
	docker rmi -f $(FQIN):latest
	docker rmi -f $(FQIN):$(TAG)

install:
	echo "TBD."

test:
	# sudo docker search $(REGISTRY)/$(FQIN)
	echo "TBD."

clean:
	# sudo atomic uninstall --force $(TEST_IMAGE)
	echo "TBD."

publish:
	docker login -u ${REGISTRY_USER} $(REGISTRY)
	#docker tag  $(FQIN) $(REGISTRY)/$(FQIN):latest
	#docker push $(REGISTRY)/$(FQIN):latest
	#docker tag  $(FQIN) $(REGISTRY)/$(FQIN):$(TAG)
	#docker push $(REGISTRY)/$(FQIN):$(TAG)
	#echo "RUN exmaple:"
	#echo "docker run --rm -p 80:80  $(REGISTRY)/$(FQIN):latest"
	#echo "docker run --rm -it -p 80:80 $(REGISTRY)/$(FQIN):latest bash"

	docker tag  $(FQIN) $(FQIN):latest
	docker push $(FQIN):latest
	docker tag  $(FQIN) $(FQIN):$(TAG)
	docker push $(FQIN):$(TAG)
	echo "RUN exmaple:"
	echo "docker run --rm -p 80:80  $(FQIN):latest"
	echo "docker run --rm -it -p 80:80 $(FQIN):latest bash"

publish-test:
	#docker search $(REGISTRY)/$(FQIN)
	docker search $(FQIN)

all: build publish publish-test

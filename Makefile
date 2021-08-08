HUGO_IMAGE := klakegg/hugo:ext
HUGO_PORT := 1313

CONTAINER_COMMAND := ${shell command -v podman 2>/dev/null || command -v docker 2>/dev/null}

BASE_COMMAND := \
	${CONTAINER_COMMAND} run \
	--tty \
	--interactive \
	--rm=true \
	-v "${PWD}":/src \
	--security-opt label=disable

YARN_COMMAND := \
	${BASE_COMMAND} \
	--entrypoint="yarn" \
	${HUGO_IMAGE}

HUGO_COMMAND := \
	${BASE_COMMAND} \
	-p ${HUGO_PORT}:${HUGO_PORT} \
	${HUGO_IMAGE}

.PHONY: all
all: build

.PHONY: dependencies
dependencies:
ifndef CONTAINER_COMMAND
    $(error "Neither Docker nor Podman could be found.")
endif
	${YARN_COMMAND} install

.PHONY: build
build: dependencies
	${HUGO_COMMAND} --minify
	# If we run using Docker, we should reset file ownership afterwards.
ifneq (,$(findstring docker,${CONTAINER_COMMAND}))
	sudo chown -R ${shell id -u ${USER}}:${shell id -g ${USER}} public/
endif

.PHONY: server
server: dependencies
	${HUGO_COMMAND} server --minify --buildDrafts

.PHONY: clean
clean:
	rm -rf ./node_modules/
	rm -rf ./public/
	rm -rf ./resources/_gen/

#!/bin/sh -e
# watchexec make 
PROXY_CONTAINER=qrcode-proxy

function _do_up(){
	# run api server
	#bin/qrcode-api &

	# run grpc proxy
	docker run -d \
		-p 8080:80 \
		-v ${PWD}/nginx.conf:/etc/nginx/conf.d/default.conf:ro \
		--name ${PROXY_CONTAINER} nginx
}

function _do_dev_server() {
	watchexec -v --print-events -r --delay-run 1 \
		-e go \
		-- 'make && bin/qrcodeapi grpc-server'
}

function _do_down() {
	docker rm -f ${PROXY_CONTAINER}
}

cmd=$1
shift

_do_${cmd} @$


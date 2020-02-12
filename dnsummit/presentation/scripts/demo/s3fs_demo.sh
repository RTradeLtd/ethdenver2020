#! /bin/bash

case "$1" in 

    mount)
        if [[ ! -d ./dnet-s3fs-demo ]]; then
            mkdir ./dnet-s3fs-demo
        fi
        sudo s3fs -o url=http://127.0.0.1:9000 -o use_path_request_style -o default_acl=public-read -o complement_stat -o sigv2 -o mp_umask=0022 -o dbglevel=debug -f -o curldbg testdata ./dnet-s3fs-demo
        ;;
    *)
        echo "invalid invocation"
        echo "./s3fs_demo.sh <mount | unmount>"
        exit 1
        ;;
esac

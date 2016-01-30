Barebones Docket Volume Plugin
==============================

This is a simple example of a docker volume plugin. This requires a working
installation of go to build.

To build (on a linux system running docker)

    go install

To run (on a linux system running docker)

    go run main.go

    # or after building

    $GOPATH/bin/barebones

If you are working on a Mac, you can use the gox cross compiler

    go get github.com/mitchellh/gox
    gox -osarch="linux/amd64"

If cross compiling, you may then wish to copy to a local docker-machine
instance and run for testing.

    docker-machine scp ./barebones_linux_amd64 <machine name>:~/
    docker-machine ssh <machine name>
    sudo ./barebones_linux_amd64

Once you have a copy of docker daemon and barebones running, test it with:

    docker run -ti -v volumename:/data --volume-driver=barebones busybox sh
    touch /data/helloworld

Verify the volume was created and the file exists (on the machine with docker
daemon)

    ls /tmp/volumename
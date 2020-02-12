title: Building (Decentralized) Applications On IPFS Using S3
author:
   name: Alex Trottier
   url: https://github.com/bonedaddy
   email: postables@rtradetechnologies.com
theme: themes/cleaver-dark
controls: false
--

## Building (Decentralized) Applications On IPFS Using S3

<html>
<div style="text-align: center;">
Demo: turn web 2.0 tooling into web 3.0 in 5 minutes 
</div>
<html>

--

## <u> Current Problems With Using IPFS </u>

* Unstable API implementations
* Lack of built-in API authentication and authorization
* Slowness that gets worse as your data demands scale
* Requires specific integration
  * Design your application from the ground up for using IPFS
  * Redesign your existing application for using IPFS

--

## <u> Why S3 </u>

* The S3 API is probably one of the most widely used APIs
  * Majority of fortune 500 companies use S3 or MinIO in some capacity
* Robust support by Amazon, and a lively third-party ecosystem
* Provide an S3 API into IPFS, and you have instant IPFS access for tens of thousands of companies

-- 

## <u> How? </u>

* TemporalX - High performance, production ready IPFS
  * Fast data Ingestion and processing
  * Built-in data replication
* MinIO - High performance object storage
  * 95% of the work already done for us (primary reason tbh)
  * A large number of companies are already familiar with minio
  * Docker support

--

## <u> S3X - TemporalX MinIO Gateway - 1 </u>

* Using TemporalX gives us an interface into IPFS suitable for production workloads
  * Fast, lightweight, and microservice friendly
* Uses a custom gateway that satisfies the MinIO `ObjectLayer` interface

--

## <u> S3X - TemporalX MinIO Gateway - 2 </u>

* Because we satisfy the `ObjectLayer` interface, we inherit a few benefits:
  * Authentication and authorization
  * Stable API implementation, usable by any S3 client
* Custom gateway has an internal "ledger store"
  * Basically a dope key-value store with a few nice features

--

## <u> S3X - TemporalX MinIO Gateway - 3 </u>

* Ledger store provides ordered, and synchronized access to the same bucket and objects
  * No race conditions regarding multiple readers/writers accessing the same bucket
  * Downside to this is that certain portions of the code paths are blocking when handling multiple read/writes for the same bucket+object, although this has been minimized as much as possible
  * Upside is that multiple requests for different buckets+objects will not be blocked
* Because of this, you can use our S3 API as a backend for `go-ds-s3` which has problems with multiple readers/writers to the same buckets ;)

--

## <u> Demos </u>

* S3 API can be used for a wide variety of tasks
* Not enough time to build a full application
* Instead I'll show how you can take existing tooling and make it IPFS compatible in 5 minutes:
  * Database and file system backups
  * FUSE filesystem mount

--

## <u> Database Backups  </u>

* For this we're going to use [wal-g](https://github.com/wal-g/wal-g)
* `wal-g` allows backing up different databases to S3

<html>
<style>
.center {
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: 50%;
}
</style>
<video width="426" height="320" controls class="center">
  <source src="https://gateway.temporal.cloud/ipfs/QmNso48sCzCiJVWLnic45wGofJ5Ucqx1tcEZXe4h3hNEzT" type="video/webm">
</video>
</html>

--

## <u> Filesystem Backup </u>

* For this we're going to use [rclone](https://github.com/rclone/rclone)
* `rclone` can be used to backup "stuff" to a variety of backends including S3

<html>
<style>
.center {
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: 50%;
}
</style>
<video width="426" height="320" controls class="center">
  <source src="https://gateway.temporal.cloud/ipfs/QmboeuYfu7k333bhiSZ6GNjHGqtEkDnthYHr1oyAn4TrY4" type="video/webm">
</video>
</html>

--

## <u> FUSE Filesystem </u>

* This demo is really finnicky so unfortunately a pre-recorded video will have to suffice.
* TODO(bonedaddy): incldue video

--

## <u> The Data On IPFS - 1 </u>

* Ultimately the data is stored on IPFS
* Which means we have two ways to get to our data:
  * Through the app that generated the data
    * Boring, so I don't demo this
    * See me after if you want to see this
  * Through IPFS:
    * Cool, so I'll demo this
    * Pinning services
    * IPFS nodes

--

## <u> The Data On IPFS - 2 </u>

* We'll get the swarm address of the TemporalX node so we can connect to `go-ipfs`.
* Using the docker logs we'll grab a couple of the CID's generated for some of the buckets and objects
* Then we'll pin the CID's on the `go-ipfs` node, and run a couple of commands:
  * `ipfs dag get`
  * `ipfs block get`
  * `ipfs block stat`

--

## <u> The Data On IPFS - 3 </u>

* Note that the way the data is displayed will look weird because they're not a default IPLD object type
  * The type definitions are available on [github](https://github.com/RTradeLtd/s3x/blob/master/cmd/gateway/s3x/s3.proto)
* Using the type definitions you can then start consuming this data within your IPFS based apps

--

## <u> The Data On IPFS - Video Backup </u>

* Here in case the live demo messes up
* Our local TemporalX docker node can be reached via `/ip4/0.0.0.0/tcp/4005/p2p/12D3KooWGSTKKGK99jMqLWyqxwtnT7tRkLyve6hhXRk5dopHKGrZ`
* So we'll connect to the TemporalX node through a `go-ipfs` node, and inspect the data

<html>
<style>
.center {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

</style>
<video width="426" height="320" controls class="center">
  <source src="https://gateway.temporal.cloud/ipfs/QmQJqb1ZQ6i8ZVK3xAQ4GiEPyYfdLNmLaeXPZgMqXyCf2A" type="video/webm">
</video>
</html>

--

## <u> Other Things You Can Do - 1 </u>

<b>IPFS Inception</b>
* Run the TemporalX node the MinIO gateway needs on a private network
* Run a go-ipfs node on the public network using `go-ds-s3` pointed to our MinIO gateway
* Your go-ipfs node is now storing data on two different networks at once
* Due to the way our gateway is designed, this content can be accessed directly on either of the IPFS networks

--

## <u> Other Things You Can Do - 2 </u>

* [IPFS powered FTP server](https://cloudacademy.com/blog/s3-ftp-server/)
* [IPFS powered video streaming service](https://antmedia.io/)
* [IPFS powered TravisCI code release storage](https://docs.travis-ci.com/user/deployment/s3/)
* [IPFS powered S3 Jenkins Plugin](https://plugins.jenkins.io/s3)
* [IPFS powered literally anything on the minio awesome list](https://github.com/minio/awesome-minio)

--

## <u> Other Things You Can Do - 3 </u>

* Switch your existing applications that use S3 and transfer them over to IPFS simply by changing a URL
* Backup your website using the MinIO command line client `mc`
* Actually have a reasonable migration path off AWS 
* Many more...

--

## <u> Can This Be Done With Go-IPFS </u>

* Undoubtedly some of are you are thinking "well we can just do this with go-ipfs, no need for TemporalX"
  * This is probably true if you don't want to support production workloads
  * But, if you need to support production workloads it is woefully impossible
* Pinning system negatively impacts performance at scale
* Memory consumption of an equally loaded go-ipfs node is 7-10x higher than TemporalX
* TemporalX's gRPC API lead to drastically faster API response times

--

## <u> Big Shout Out To MinIO </u>

* We aren't affiliated with MinIO
* Without them however, this likely wouldn't be possible, or at the very least have required an immense amount of work, and resources that aren't available to us

--

## <u> Tools Used </u>

* [`wal-g`](https://github.com/wal-g/wal-g) (open-source)
* [`rclone`](https://github.com/rclone/rclone) (open-source)
* [`minio`](https://github.com/minio/minio) (open-source)
* [`s3x` - minio fork](https://github.com/RTradeLtd/s3x) (open-source)
* [`mc`](https://github.com/minio/mc) (open-source)
* [`temporalx`](https://temporal.cloud/temporalx/) (partially open-source)

--

## <u> Questions? </u>

* Note #1: if you want to discuss the architecture of the gateway lets do that after the talk as it can take 
* Note #2: More resources available at our website https://temporal.cloud and on telegram https://t.me/RTradeTEMPORAL
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

  <style>
    .holder {   
      width: auto;
      margin-top: 300px;
      position: absolute;
      display: inline-block;    
    }
    .holder img {
      width: 30%; /* Will shrink image to 30% of its original width */
      position: relative;
      display: block;
      margin-left: 500px;
      margin-bottom: 225px;
    }​
    .ip {
      margin-left: 500px;
      margin-bottom: 225px;
    }
  </style>
  <div style="text-align: center;">
    Demo: turn web 2.0 tooling into web 3.0 in 5 minutes 
  </div>
  <div class="holder">
    <p class="ip" style="margin-left: 475px">Powered By:</p>
    <img src="https://gateway.temporal.cloud/ipfs/QmWf9krER4BVZidfYjDSB2BrPdGPQfB2kkFFdXZSNUj9Eu"/>
  </div>

</html>

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

* S3 API one of the most widely deployed APIs
* Robust support by Amazon + third party ecosystem
* S3 IPFS API = instant ipfs access for thousands of companies

-- 

## <u> How? </u>

* TemporalX - High performance, production ready IPFS
  * Fast data Ingestion and processing
  * Built-in data replication
* MinIO - High performance object storage
  * Battle-tested S3 framework
  * Used by thousands of companies already
  * Docker support

--

## <u> What? </u>

* TemporalX + MinIO Gateway = S3X
* S3X = S3 API powered by IPFS

--

## <u> S3X - 1 </u>

* TemporalX provides the IPFS interface
* Custom gateway bridges TemporalX and MinIO through an `ObjectLayer`
* "Object data" stored as-is on IPFS separate from S3 data

--

## <u> S3X- 2 </u>

* Satisfying `ObjectLayer` means we inherit MinIO benefits:
  * Authentication and authorization
  * Stable API implementation (usable by many)
  * etc...
* Uses an internal "ledger store" to track metadata

--

## <u> S3X - 3 </u>

* Ledger store facilitates synchronized and ordered access to data
* Multiple requests + different buckets/objects = non-blocking
* Multiple requests + same buckets/objects = blocking
* Can use S3X as a backend for `go-ds-s3` and multiple `go-ipfs` nodes without issues :O 

--

## <u> Demos </u>

* S3 API can be used for a wide variety of tasks
* Transform non-IPFS tooling into IPFS tooling +/- 5 minutes
  * Database + filesystem backups
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

* Using `s3fs` we can mount buckets locally via FUSE
* Basically an IPFS fuse filesystem without the `go-ipfs` issues
* TODO(bonedaddy): try and get working before presentation

--

## <u> The Data On IPFS - 1 </u>

* Ultimately the data is stored on IPFS
* Which means we have two ways to get to our data:
  * Through the app that generated the data (boring)
  * Through IPFS (cool)

--

## <u> The Data On IPFS - 2 </u>

* Actual object data stored separately
* Object data can be loaded through a gateway
* Objects+buckets stored as [IPLD objects](https://github.com/RTradeLtd/s3x/blob/master/cmd/gateway/s3x/s3.proto)

--

## <u> The Data On IPFS - 3 </u>

* Connect TemporalX and `go-ipfs` together
* Inspect data on IPFS through `go-ipfs`

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
* Run a `go-ipfs` node on the public network using `go-ds-s3` pointed to S3X
* Your `go-ipfs` node is now storing data on two different networks at once
* Because object data stored as is, content can be accessed directly from either network

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

* If you don't want to support production workloads, <i>maybe</I>
* Pinning negatively impacts performance at scale
  * 1TB of data pinned at 256KB chunksize = 50% performance reduction when adding/removing pins
* Memory consumption of an equally loaded go-ipfs node is 7-10x higher than TemporalX
* TemporalX gRPC API superior than HTTP API 
  * Low-latency
  * Multiplexed
  * HTTP/2

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
package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	host "github.com/libp2p/go-libp2p-host"
	dopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	datastore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	record "github.com/libp2p/go-libp2p-record"
	minio "github.com/minio/minio-go/v6"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

// from https://gist.github.com/ZenGround0/49e4a1aa126736f966a1dfdcb84abdae

const partBoundary = "123456789000000000000987654321"

const boundary = "\r\n--" + partBoundary + "\r\n"

var (
	hostAddress   = flag.String("host", "192.168.0.2:80", "the host to connect to")
	addr          = flag.String("multi.addr", "/ip4/0.0.0.0/tcp/4006", "the multiaddr for libp2p host")
	accessKey     = flag.String("access.key", "minio", "minio access key")
	secretKey     = flag.String("secret.key", "minio123", "minio secret key")
	minioEndpoint = flag.String("minio.endpoint", "0.0.0.0:9000", "minio endpoint")
	setup         = flag.Bool("setup", true, "setup the testenv then exit")
	mux           sync.Mutex
)

func init() {
	flag.Parse()
}
func main() {
	minioClient, err := minio.New(*minioEndpoint, *accessKey, *secretKey, false)
	if err != nil {
		log.Fatal("failed to access minio endpoint ", err)
	}
	exists, err := minioClient.BucketExists("testbucket")
	if err != nil {
		log.Fatal("failed to check if bucket exists ", err)
	}
	if !exists {
		if err := minioClient.MakeBucket("testbucket", "us-east-1"); err != nil {
			log.Fatal("failed to make bucket ", err)
		}
	}
	if *setup {
		log.Println("finished setup")
		os.Exit(0)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ds := dssync.MutexWrap(datastore.NewMapDatastore())
	ps := pstoremem.NewPeerstore()
	pk, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 256)
	if err != nil {
		log.Fatal("failed to generate key ", err)
	}
	logger := zap.NewExample()
	maddr, err := multiaddr.NewMultiaddr(*addr)
	if err != nil {
		log.Fatal("failed to get multiaddr ", err)
	}
	h, dt, err := newLibp2pHostAndDHT(ctx, logger, ds, ps, pk, []multiaddr.Multiaddr{maddr})
	if err != nil {
		log.Fatal("libp2p host creation failed ", err)
	}
	h.Close()
	dt.Close()
	tr := http.DefaultTransport
	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   *hostAddress,
			Path:   "/",
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to send request ", err)
	}
	// remove previous data
	minioClient.RemoveObject("testbucket", "videofeed")
	// create a new empty object
	_, err = minioClient.PutObject("testbucket", "videofeed", bytes.NewReader(nil), 0, minio.PutObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// this streams the video feed over libp2p
	libp2pStreamData := func(stream network.Stream) {
		defer stream.Reset()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			obj, err := minioClient.GetObject("testbucket", "videofeed", minio.GetObjectOptions{})
			if err != nil {
				log.Fatal(err)
			}
			data, err := ioutil.ReadAll(obj)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := stream.Write(data); err != nil {
				log.Fatal(err)
			}
		}
	}
	// register the libp2p protocol handler
	h.SetStreamHandler(protocol.ID("texiot/videostream/0.0.1"), libp2pStreamData)
	// launch a goroutine that will serve out mp4 video as a "stream"
	go func() {
		streamSrvMux := http.NewServeMux()
		streamSrvMux.HandleFunc("/", streamHandler)
		streamServer := http.Server{
			Addr:    "0.0.0.0:6969",
			Handler: streamSrvMux,
		}
		go func() {
			if err := streamServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ctx.Done():
				streamServer.Shutdown(ctx)
				return
			case <-ticker.C:
				goto START
			}
		START:
			func() {
				mux.Lock()
				defer mux.Unlock()
				obj, err := minioClient.GetObject("testbucket", "videofeed", minio.GetObjectOptions{})
				if err != nil {
					continue
				}
				objData, err := ioutil.ReadAll(obj)
				if err != nil {
					log.Fatal(err)
				}
				if err := ioutil.WriteFile("videofeed.mjpeg", objData, os.FileMode(0640)); err != nil {
					log.Fatal(err)
				}
				cmd := exec.Command(
					"ffmpeg",
					"-i",
					"videofeed.mjpeg",
					"-c:v",
					"libx264",
					"-preset",
					"veryslow",
					"-crf",
					"18",
					"output.mp4",
				)
				if err := cmd.Run(); err != nil {
					log.Fatal(err)
				}
			}()
		}
	}()
	var buf = make([]byte, 1024*1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			log.Fatal("failed to read body ", err)
		}
		obj, err := minioClient.GetObject("testbucket", "videofeed", minio.GetObjectOptions{})
		if err != nil {
			continue
		}
		data, err := ioutil.ReadAll(obj)
		if err != nil {
			log.Fatal("failed to read object ", err)
		}
		data = append(data, buf[:n]...)
		_, err = minioClient.PutObject(
			"testbucket",
			"videofeed",
			bytes.NewReader(data),
			int64(len(data)),
			minio.PutObjectOptions{},
		)
		if err != nil {
			log.Fatal("failed to put object ", err)
		}
	}

}

func handleExit(ctx context.Context, cancelFunc context.CancelFunc, wg *sync.WaitGroup, doneChan chan bool) {
	defer wg.Done()
	// make a channel to catch os signals on
	quitCh := make(chan os.Signal, 1)
	// register the types of os signals to trap
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// wait until we receive an exit signal
	<-quitCh
	// cancel the context which will trigger shutdown of service components
	cancelFunc()
	// notify that we are finished handling all exit procedures
	doneChan <- true
}

func newLibp2pHostAndDHT(
	ctx context.Context,
	logger *zap.Logger,
	ds datastore.Batching,
	ps peerstore.Peerstore,
	pk crypto.PrivKey,
	addrs []multiaddr.Multiaddr) (host.Host, *dht.IpfsDHT, error) {
	var opts []libp2p.Option
	opts = append(opts,
		libp2p.Identity(pk),
		libp2p.ListenAddrs(addrs...),
		libp2p.Peerstore(ps),
		libp2p.DefaultMuxers,
		libp2p.DefaultTransports,
		libp2p.DefaultSecurity)
	h, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, nil, err
	}

	idht, err := dht.New(ctx, h,
		dopts.Validator(record.NamespacedValidator{
			"pk":   record.PublicKeyValidator{},
			"ipns": ipns.Validator{KeyBook: ps},
		}),
	)
	if err != nil {
		return nil, nil, err
	}
	rHost := routedhost.Wrap(h, idht)
	return rHost, idht, nil
}

const BUFSIZE = 1024 * 8

func streamHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	file, err := os.Open("output.mp4")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	defer file.Close()

	fi, err := file.Stat()

	if err != nil {
		w.WriteHeader(500)
		return
	}

	fileSize := int(fi.Size())

	if len(r.Header.Get("Range")) == 0 {

		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes 0-"+contentEnd+"/"+contentLength)
		w.WriteHeader(200)

		buffer := make([]byte, BUFSIZE)

		for {
			n, err := file.Read(buffer)

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}

	} else {

		rangeParam := strings.Split(r.Header.Get("Range"), "=")[1]
		splitParams := strings.Split(rangeParam, "-")

		// response values

		contentStartValue := 0
		contentStart := strconv.Itoa(contentStartValue)
		contentEndValue := fileSize - 1
		contentEnd := strconv.Itoa(contentEndValue)
		contentSize := strconv.Itoa(fileSize)

		if len(splitParams) > 0 {
			contentStartValue, err = strconv.Atoi(splitParams[0])

			if err != nil {
				contentStartValue = 0
			}

			contentStart = strconv.Itoa(contentStartValue)
		}

		if len(splitParams) > 1 {
			contentEndValue, err = strconv.Atoi(splitParams[1])

			if err != nil {
				contentEndValue = fileSize - 1
			}

			contentEnd = strconv.Itoa(contentEndValue)
		}

		contentLength := strconv.Itoa(contentEndValue - contentStartValue + 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes "+contentStart+"-"+contentEnd+"/"+contentSize)
		w.WriteHeader(206)

		buffer := make([]byte, BUFSIZE)

		file.Seek(int64(contentStartValue), 0)

		writeBytes := 0

		for {
			n, err := file.Read(buffer)

			writeBytes += n

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			if writeBytes >= contentEndValue {
				data := buffer[:BUFSIZE-writeBytes+contentEndValue+1]
				w.Write(data)
				w.(http.Flusher).Flush()
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}
	}
}

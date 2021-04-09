package main

import (
	"archive/tar"
	"compress/gzip"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:generate tar czf headless.apkovl.tar.gz --exclude .DS_Store etc

//go:embed headless.apkovl.tar.gz
var headless []byte

var params = struct {
	arch            string
	version         string
	ssid            string
	passphrase      string
	authorized_keys string
	dist            string
}{}

func addHeadlessApkVol(dist string) error {
	fp, err := os.Create(filepath.Join(dist, "headless.apkovl.tar.gz"))
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fp.Write(headless); err != nil {
		return err
	}
	if err := fp.Sync(); err != nil {
		return err
	}
	return nil
}

func addWiFiText(dist, ssid, passphrase string) error {
	fp, err := os.Create(filepath.Join(dist, "wifi.txt"))
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fmt.Fprintf(fp, "%s %s\n", ssid, passphrase); err != nil {
		return err
	}
	if err := fp.Sync(); err != nil {
		return err
	}
	return nil
}

func addAuthorizedKeysText(dist, keys string) error {
	fp, err := os.Create(filepath.Join(dist, "authorized_keys"))
	if err != nil {
		return err
	}
	defer fp.Close()
	src, err := os.Open(keys)
	if err != nil {
		return err
	}
	defer src.Close()
	if _, err := io.Copy(fp, src); err != nil {
		return err
	}
	return nil
}

func makeURL(arch, version string) string {
	major := strings.Join(strings.Split(version, ".")[0:2], ".")
	version = strings.TrimPrefix(version, "v")
	return fmt.Sprintf("https://dl-cdn.alpinelinux.org/alpine/%s/releases/aarch64/alpine-rpi-%s-%s.tar.gz", major, version, arch)
}

func writeItem(fpath string, r io.Reader) error {
	fp, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, r); err != nil {
		return err
	}
	if err := fp.Sync(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.StringVar(&params.arch, "arch", "aarch64", "target architecture armv7/armhf/aarch64")
	flag.StringVar(&params.version, "version", "v3.13.4", "target alpinelinux version")
	flag.StringVar(&params.ssid, "ssid", "", "initial Wi-Fi ssid")
	flag.StringVar(&params.passphrase, "passphrase", "", "initial Wi-Fi passphrase")
	flag.StringVar(&params.authorized_keys, "authorized_keys", "", "initial root ssh authorized_keys")
	flag.StringVar(&params.dist, "dist", "dist", "output directory")
	flag.Parse()
	u := makeURL(params.arch, params.version)
	log.Println("download:", u)
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
	}
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fpath := filepath.Join(params.dist, header.Name)
		log.Print(fpath)
		switch header.Typeflag {
		default:
			log.Print("unknown type:", header.Typeflag)
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, 0o755); err != nil {
				log.Fatalln(err)
			}
		case tar.TypeReg:
			if err := writeItem(fpath, reader); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := addHeadlessApkVol(params.dist); err != nil {
		log.Fatal(err)
	}
	if params.ssid != "" {
		if err := addWiFiText(params.dist, params.ssid, params.passphrase); err != nil {
			log.Fatal(err)
		}
	}
	if params.authorized_keys != "" {
		if err := addAuthorizedKeysText(params.dist, params.authorized_keys); err != nil {
			log.Fatal(err)
		}
	}
}

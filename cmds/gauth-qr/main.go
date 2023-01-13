package main

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/jbert/gauthQR"
	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Must provide exactly one filename argument")
	}
	err := run(os.Args[1], os.Stdout)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func run(imageFilename string, w io.Writer) error {
	imageData, err := runZBar(imageFilename)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(imageData, "QR-Code:") {
		return fmt.Errorf("Couldn't find 'QR-Code:' prefix in [%s]", imageData)
	}
	imageData = imageData[8:]

	u, err := url.Parse(imageData)
	if err != nil {
		return fmt.Errorf("Can't parse image data [%s] as url: %s", imageData, err)
	}
	values := u.Query()

	authDataB64, ok := values["data"]
	if !ok {
		return fmt.Errorf("Can't find 'data' query param in url [%s]", u)
	}
	if len(authDataB64) != 1 {
		return fmt.Errorf("Didn't find exactly one data query param (%d)", len(authDataB64))
	}
	authData, err := base64.StdEncoding.DecodeString(authDataB64[0])
	if err != nil {
		return fmt.Errorf("Can't decode [%s] as base64: %s", authDataB64[0], err)
	}

	// OK, this is a protobuf
	mp, err := parseProtobuf(authData)
	if err != nil {
		return err
	}

	for _, op := range mp.OtpParameters {
		b32secret := base32.StdEncoding.EncodeToString(op.Secret)
		fmt.Printf("%s: %s\n", op.Name, b32secret)
	}

	return nil
}

func parseProtobuf(buf []byte) (*gauthQR.MigrationPayload, error) {
	var mp gauthQR.MigrationPayload
	err := proto.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(buf, &mp)
	if err != nil {
		return nil, fmt.Errorf("Can't parse [%02X] as protobuf: %s", buf, err)
	}
	return &mp, nil
}

func runZBar(fname string) (string, error) {
	cmd := exec.Command("zbarimg", "-q", fname)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error running \"zbar -q %s\": %s", fname, err)
	}
	if out[len(out)-1] != '\n' {
		return "", fmt.Errorf("No trailing newline on [%s]", out)
	}
	return string(out[:len(out)-1]), nil
}

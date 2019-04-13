package tls

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "github.com/quan-to/slog"
    "math/big"
    "net"
    "os"
    "time"
)

func GenerateHTTPSKeyPair() (ca *tls.Certificate, server *tls.Certificate) {
    var err error

    now := time.Now()
    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
    caKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    if err != nil {
        slog.Scope("RSTLS").Fatal(err)
    }

    caTemplate := x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"RadioServer"},
            CommonName:   "RadioServer SelfGenerated CA",
        },
        NotBefore:             now,
        NotAfter:              now.AddDate(1, 0, 0),
        KeyUsage:              x509.KeyUsageCertSign,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
        IsCA:                  true,
    }

    caCertBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
    if err != nil {
        slog.Scope("RSTLS").Fatal(err)
    }
    ca = &tls.Certificate{}
    ca.Certificate = append(ca.Certificate, caCertBytes)

    serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        slog.Scope("RSTLS").Fatal(err)
    }

    serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        slog.Scope("RSTLS").Fatal("failed to generate serial number: %s", err)
    }
    serverTemplate := x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"RadioServer"},
            CommonName:   "RadioServer",
        },
        NotBefore:             now,
        NotAfter:              now.AddDate(1, 0, 0),
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
        IsCA:                  false,
    }

    ips := GetIPAddressList()

    for _, v := range ips {
        serverTemplate.IPAddresses = append(serverTemplate.IPAddresses, v)
    }

    serverTemplate.IPAddresses = append(serverTemplate.IPAddresses, net.IPv4(127, 0, 0, 1))
    serverTemplate.DNSNames = append(serverTemplate.DNSNames, "localhost")

    hostname := GetHostname()
    if hostname != "" {
        serverTemplate.DNSNames = append(serverTemplate.DNSNames, hostname)
    }

    serverCertBytes, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, &serverKey.PublicKey, caKey)
    if err != nil {
        panic(err)
    }


    server = &tls.Certificate{}
    server.Certificate = append(server.Certificate, serverCertBytes)
    server.PrivateKey = serverKey

    return
}

func GetHostname() string {
    name, _ := os.Hostname()
    return name
}

func GetIPAddressList() [][]byte {
    ifaces, _ := net.Interfaces()

    var ips [][]byte

    // handle err
    for _, i := range ifaces {
        addrs, _ := i.Addrs()
        // handle err
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            ips = append(ips, ip)
            // process IP address
        }
    }

    return ips
}
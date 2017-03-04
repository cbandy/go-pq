package pq

// This file contains SSL tests

import (
	_ "crypto/sha256"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func maybeSkipSSLTests(t *testing.T) {
	// Require some special variables for testing certificates
	if os.Getenv("PQSSLCERTTEST_PATH") == "" {
		t.Skip("PQSSLCERTTEST_PATH not set, skipping SSL tests")
	}

	value := os.Getenv("PQGOSSLTESTS")
	if value == "" || value == "0" {
		t.Skip("PQGOSSLTESTS not enabled, skipping SSL tests")
	} else if value != "1" {
		t.Fatalf("unexpected value %q for PQGOSSLTESTS", value)
	}

	// Environment sanity check: should fail without SSL
	if conninfo := "sslmode=disable user=pqgossltest"; checkSSLConn(t, conninfo) == nil {
		t.Fatalf("expected error with conninfo=%q", conninfo)
	}
}

func checkSSLConn(t *testing.T, conninfo string) error {
	db, err := openTestConnConninfo(conninfo)
	if err != nil {
		// should never fail
		t.Fatal(err)
	}
	defer db.Close()

	// Do something with the connection to see whether it's working or not.
	tx, err := db.Begin()
	if err == nil {
		err = tx.Rollback()
	}
	return err
}

func TestSSLRequire(t *testing.T) {
	maybeSkipSSLTests(t)

	err := checkSSLConn(t, "sslmode=require user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}
}

// Test sslmode=verify-full
func TestSSLVerifyFull(t *testing.T) {
	maybeSkipSSLTests(t)

	// Not OK according to the system CA
	err := checkSSLConn(t, "host=postgres sslmode=verify-full user=pqgossltest")
	_, ok := err.(x509.UnknownAuthorityError)
	if !ok {
		t.Fatalf("expected x509.UnknownAuthorityError, got %#+v", err)
	}

	rootCertPath := filepath.Join(os.Getenv("PQSSLCERTTEST_PATH"), "root.crt")
	rootCert := "sslrootcert=" + rootCertPath + " "

	// No match on Common Name
	err = checkSSLConn(t, rootCert+"host=127.0.0.1 sslmode=verify-full user=pqgossltest")
	_, ok = err.(x509.HostnameError)
	if !ok {
		t.Fatalf("expected x509.HostnameError, got %#+v", err)
	}

	// OK
	err = checkSSLConn(t, rootCert+"host=postgres sslmode=verify-full user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}
}

// Test sslmode=require sslrootcert=rootCertPath
func TestSSLRequireWithRootCert(t *testing.T) {
	maybeSkipSSLTests(t)

	bogusRootCertPath := filepath.Join(os.Getenv("PQSSLCERTTEST_PATH"), "bogus_root.crt")
	bogusRootCert := "sslrootcert=" + bogusRootCertPath + " "

	// Not OK according to the bogus CA
	err := checkSSLConn(t, bogusRootCert+"host=postgres sslmode=require user=pqgossltest")
	_, ok := err.(x509.UnknownAuthorityError)
	if !ok {
		t.Fatalf("expected x509.UnknownAuthorityError, got %s, %#+v", err, err)
	}

	nonExistentCertPath := filepath.Join(os.Getenv("PQSSLCERTTEST_PATH"), "non_existent.crt")
	nonExistentCert := "sslrootcert=" + nonExistentCertPath + " "

	// No match on Common Name, but that's OK because we're not validating anything.
	err = checkSSLConn(t, nonExistentCert+"host=127.0.0.1 sslmode=require user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}

	rootCertPath := filepath.Join(os.Getenv("PQSSLCERTTEST_PATH"), "root.crt")
	rootCert := "sslrootcert=" + rootCertPath + " "

	// No match on Common Name, but that's OK because we're not validating the CN.
	err = checkSSLConn(t, rootCert+"host=127.0.0.1 sslmode=require user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}

	// Everything OK
	err = checkSSLConn(t, rootCert+"host=postgres sslmode=require user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}
}

// Test sslmode=verify-ca
func TestSSLVerifyCA(t *testing.T) {
	maybeSkipSSLTests(t)

	// Not OK according to the system CA
	err := checkSSLConn(t, "host=postgres sslmode=verify-ca user=pqgossltest")
	_, ok := err.(x509.UnknownAuthorityError)
	if !ok {
		t.Fatalf("expected x509.UnknownAuthorityError, got %#+v", err)
	}

	rootCertPath := filepath.Join(os.Getenv("PQSSLCERTTEST_PATH"), "root.crt")
	rootCert := "sslrootcert=" + rootCertPath + " "

	// No match on Common Name, but that's OK
	err = checkSSLConn(t, rootCert+"host=127.0.0.1 sslmode=verify-ca user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}

	// Everything OK
	err = checkSSLConn(t, rootCert+"host=postgres sslmode=verify-ca user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}
}

func getCertConninfo(t *testing.T, source string) string {
	var sslkey string
	var sslcert string

	certpath := os.Getenv("PQSSLCERTTEST_PATH")

	switch source {
	case "missingkey":
		sslkey = "/tmp/filedoesnotexist"
		sslcert = filepath.Join(certpath, "postgresql.crt")
	case "missingcert":
		sslkey = filepath.Join(certpath, "postgresql.key")
		sslcert = "/tmp/filedoesnotexist"
	case "certtwice":
		sslkey = filepath.Join(certpath, "postgresql.crt")
		sslcert = filepath.Join(certpath, "postgresql.crt")
	case "valid":
		sslkey = filepath.Join(certpath, "postgresql.key")
		sslcert = filepath.Join(certpath, "postgresql.crt")
	default:
		t.Fatalf("invalid source %q", source)
	}
	return fmt.Sprintf("sslmode=require user=pqgosslcert sslkey=%s sslcert=%s", sslkey, sslcert)
}

// Authenticate over SSL using client certificates
func TestSSLClientCertificates(t *testing.T) {
	maybeSkipSSLTests(t)

	// Should also fail without a valid certificate
	err := checkSSLConn(t, "sslmode=require user=pqgosslcert")
	pge, ok := err.(*Error)
	if !ok {
		t.Fatal("expected pq.Error")
	}
	if pge.Code.Name() != "invalid_authorization_specification" {
		t.Fatalf("unexpected error code %q", pge.Code.Name())
	}

	// Should work
	err = checkSSLConn(t, getCertConninfo(t, "valid"))
	if err != nil {
		t.Fatal(err)
	}
}

// Test errors with ssl certificates
func TestSSLClientCertificatesMissingFiles(t *testing.T) {
	maybeSkipSSLTests(t)

	// Key missing, should fail
	err := checkSSLConn(t, getCertConninfo(t, "missingkey"))
	_, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("expected PathError, got %#+v", err)
	}

	// Cert missing, should fail
	err = checkSSLConn(t, getCertConninfo(t, "missingcert"))
	_, ok = err.(*os.PathError)
	if !ok {
		t.Fatalf("expected PathError, got %#+v", err)
	}

	// Key has wrong permissions, should fail
	err = checkSSLConn(t, getCertConninfo(t, "certtwice"))
	if err != ErrSSLKeyHasWorldPermissions {
		t.Fatalf("expected ErrSSLKeyHasWorldPermissions, got %#+v", err)
	}
}

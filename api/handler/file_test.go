package handler

import (
	"net/http/httptest"
	"testing"

	"file-service/service/file"
)

func TestWithPublicStorageURLReplacesLocalhostWithRequestHost(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest("GET", "http://147.45.98.120/api/file-service/files/7", nil)
	stored := file.File{
		ID:       7,
		FileName: "document.docx",
		FileSize: 19280,
		Bucket:   file.BucketDocuments,
		URL:      "http://localhost/storage/documents/document.docx",
	}

	got := withPublicStorageURL(request, stored)

	if got.URL != "http://147.45.98.120/storage/documents/document.docx" {
		t.Fatalf("URL = %q, want %q", got.URL, "http://147.45.98.120/storage/documents/document.docx")
	}
}

func TestWithPublicStorageURLUsesForwardedHeaders(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest("GET", "http://file-service/files/7", nil)
	request.Header.Set("X-Forwarded-Host", "tns.quassbot.ru")
	request.Header.Set("X-Forwarded-Proto", "https")
	stored := file.File{
		ID:       7,
		FileName: "document.docx",
		FileSize: 19280,
		Bucket:   file.BucketDocuments,
		URL:      "http://localhost/storage/documents/document.docx",
	}

	got := withPublicStorageURL(request, stored)

	if got.URL != "https://tns.quassbot.ru/storage/documents/document.docx" {
		t.Fatalf("URL = %q, want %q", got.URL, "https://tns.quassbot.ru/storage/documents/document.docx")
	}
}

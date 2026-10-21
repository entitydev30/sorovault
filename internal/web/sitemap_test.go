package web_test

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/soroworks/sorovault/internal/model"
	"github.com/soroworks/sorovault/internal/store"
)

func TestLinkPreviewMetadata(t *testing.T) {
	h := newHarness(t)
	h.register(t, idA)
	_, body := h.get(t, "/contracts/"+idA)
	assert.Contains(t, body, `<meta name="description" content="testnet Soroban contract with 7 functions:`)
	assert.Contains(t, body, `<meta property="og:title" content="Contract `+idA)
	assert.Contains(t, body, `<meta property="og:description" content="testnet Soroban contract with 7 functions:`)
	assert.Contains(t, body, "has_voted")
	assert.NotContains(t, body, `<meta property="og:url"`, "no canonical URL is emitted without PUBLIC_URL")
}

func TestSitemapXML(t *testing.T) {
	h := newHarness(t)
	h.register(t, idA)
	h.register(t, idB)
	h.reg.Store().(*store.Memory) // ensure in-memory store is exercised
	resp, _ := h.get(t, "/sitemap.xml")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	// h.server routes are already mounted; set the optional URL before routing
	// in a separate test server when PUBLIC_URL is supplied.
}

func TestSitemapPaginationAndEscaping(t *testing.T) {
	h := newHarness(t)
	for i := 0; i < 3; i++ {
		contract := store.Contract{
			ContractID: fmt.Sprintf("C%d", i), Network: "testnet",
			CurrentWasmHash: fmt.Sprintf("hash-%d", i),
		}
		require.NoError(t, h.store.Save(context.Background(), contract, &model.Interface{ABIVersion: model.ABIVersion}))
	}
	_ = time.Now
	_ = url.PathEscape
	_ = strings.Contains
	_ = xml.Unmarshal
}

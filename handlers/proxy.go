package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"proxy/config"
	"proxy/utils"
)

var excludedHeaders = map[string]bool{
	"content-encoding":  true,
	"content-length":    true,
	"transfer-encoding": true,
	"connection":        true,
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := config.TargetURL + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") && len(body) > 0 {
		newBody, converted, original, convertedValue, err := utils.ConvertPreciseAmount(body)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidPreciseAmount) {
				fmt.Printf("⚠ No se pudo convertir precise_amount: %s\n", original)
			} else {
				fmt.Printf("Error procesando JSON: %v\n", err)
			}
		} else if converted {
			body = newBody
			fmt.Printf("✓ Convertido precise_amount: '%s' (str) -> %d (int)\n", original, convertedValue)
		}
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Error creando request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		if strings.ToLower(key) != "host" {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}
	}

	if len(body) > 0 {
		proxyReq.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	for _, cookie := range r.Cookies() {
		proxyReq.AddCookie(cookie)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al conectar con el servidor: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		if !excludedHeaders[strings.ToLower(key)] {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

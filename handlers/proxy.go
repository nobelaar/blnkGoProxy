package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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
		msg := fmt.Sprintf("error leyendo el body de la solicitud entrante: %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logHTTPMessage("Solicitud entrante", r.Method, r.URL.String(), r.Header, body)

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
			fmt.Printf("✓ Convertido precise_amount: '%s' (str) -> %s (numero)\n", original, convertedValue.String())
		}
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		msg := fmt.Sprintf("error creando la request hacia el destino: %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
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

	logHTTPMessage("Solicitud reenviada al destino", proxyReq.Method, proxyReq.URL.String(), proxyReq.Header, body)

	resp, err := client.Do(proxyReq)
	if err != nil {
		msg := fmt.Sprintf("error al conectar con el servidor de destino: %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("error leyendo la respuesta del servidor de destino: %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadGateway)
		return
	}

	logHTTPResponse("Respuesta del servidor de destino", resp, respBody)

	for key, values := range resp.Header {
		if !excludedHeaders[strings.ToLower(key)] {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(resp.StatusCode)

	if len(respBody) > 0 {
		if _, err := w.Write(respBody); err != nil {
			log.Printf("error enviando la respuesta al cliente: %v\n", err)
		}
	}
}

// logHTTPMessage imprime de forma legible el request que llega o se envía.
func logHTTPMessage(title, method, url string, headers http.Header, body []byte) {
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Printf("%s %s\n", method, url)
	fmt.Println("Headers:")
	for key, values := range headers {
		fmt.Printf("- %s: %s\n", key, strings.Join(values, ", "))
	}
	if len(body) > 0 {
		fmt.Printf("Body: %s\n", string(body))
	} else {
		fmt.Println("Body: <vacío>")
	}
	fmt.Println("===============")
}

// logHTTPResponse imprime el status, headers y body de la respuesta del destino.
func logHTTPResponse(title string, resp *http.Response, body []byte) {
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		fmt.Printf("- %s: %s\n", key, strings.Join(values, ", "))
	}
	if len(body) > 0 {
		fmt.Printf("Body: %s\n", string(body))
	} else {
		fmt.Println("Body: <vacío>")
	}
	fmt.Println("===============")
}

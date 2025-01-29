package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dostonlv/url-shortner-http3-based/internal/service"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

type Handler struct {
	services *service.URLService
}

func NewHandler(services *service.URLService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.handleShorten)
	mux.HandleFunc("/{short}", h.handleRedirect)
	mux.HandleFunc("/stats/", h.handleStats)
	mux.HandleFunc("/docs", h.docs)

	return mux
}

// @Summary Create short URL
// @Tags urls
// @Accept json
// @Produce json
// @Param input body models.URL true "URL info"
// @Success 200 {object} models.URL
// @Router /shorten [post]
func (h *Handler) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		URL string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	input.URL = validateAndFormatURL(input.URL)

	url, err := h.services.CreateShortURL(input.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": fmt.Sprintf("https://localhost:4433/%s", url.ShortCode),
	})
}

func validateAndFormatURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}

// @Summary Redirect to original URL
// @Tags urls
// @Accept json
// @Produce json
// @Param short_code path string true "Short code"
// @Success 302 {object} models.URL
// @Router /{short_code} [get]
// Optimized handleRedirect function
func (h *Handler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	url, err := h.services.GetURL(shortCode)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Redirecting to: %s\n", url.OriginalURL)

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

// @Summary Get URL stats
// @Tags urls
// @Accept json
// @Produce json
// @Param short_code path string true "Short code"
// @Success 200 {object} models.URL
// @Router /stats/{short_code} [get]
func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[len("/stats/"):]
	url, err := h.services.GetURL(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

// router.Get("/reference", func(w http.ResponseWriter, r *http.Request) {
// 	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
// 		SpecURL: "./docs/swagger.json",
// 		CustomOptions: scalar.CustomOptions{
// 			PageTitle: "Simple API",
// 		},
// 		DarkMode: true,
// 	})

// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	fmt.Fprintln(w, htmlContent)
// })

func (h *Handler) docs(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/swagger.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "URL Shortener API",
		},
		DarkMode: true,
	})
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Fprintln(w, htmlContent)

}

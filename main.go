// Package main is the entry point of the go-web-app web server.
//
// This service is a simple Go HTTP server that serves four static HTML pages.
// The application is intentionally minimal. The real complexity lives in the
// DevOps infrastructure around it — Docker, Kubernetes, Helm, and CI/CD.
//
// Routes:
//
//	GET /home     → serves static/home.html
//	GET /courses  → serves static/courses.html
//	GET /about    → serves static/about.html
//	GET /contact  → serves static/contact.html
//	GET /healthz  → liveness probe  (used by Kubernetes to check if the app is running)
//	GET /readyz   → readiness probe (used by Kubernetes to check if the app is ready for traffic)
//
// The server listens on 0.0.0.0:8080 so it is reachable inside a Docker
// container from any network interface.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ---------------------------------------------------------------------------
// Page Handlers
// ---------------------------------------------------------------------------

// homePage serves the home page.
func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/home.html")
}

// coursePage serves the courses page.
func coursePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/courses.html")
}

// aboutPage serves the about page.
func aboutPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/about.html")
}

// contactPage serves the contact page.
func contactPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/contact.html")
}

// ---------------------------------------------------------------------------
// Health Check Handlers
// ---------------------------------------------------------------------------

// healthzHandler is the liveness probe endpoint.
//
// Kubernetes calls this endpoint on a schedule. If it returns a non-200 status,
// Kubernetes restarts the container. This handler must be very fast and must
// never do any I/O — its only job is to confirm the process is still running.
//
// Route: GET /healthz
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

// readyzHandler is the readiness probe endpoint.
//
// Kubernetes calls this before sending any traffic to a new pod. If it returns
// non-200, the pod is temporarily removed from the load balancer but is NOT
// restarted. This lets the app signal that it is still starting up.
//
// For this static-file server there is no startup work, so readiness mirrors
// liveness. In a more complex app, you would check database connections here.
//
// Route: GET /readyz
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ready")
}

// ---------------------------------------------------------------------------
// Server Setup
// ---------------------------------------------------------------------------

// main starts the HTTP server and waits for a shutdown signal.
//
// When Kubernetes wants to stop a pod (during a rolling update or scale-down),
// it sends a SIGTERM signal. Without a graceful shutdown handler, any requests
// currently being processed would be dropped immediately. This function waits
// up to 30 seconds for in-flight requests to finish before the server exits.
func main() {
	// Register all URL routes on a new request multiplexer.
	// Using a local mux (instead of the default global one) avoids accidental
	// route registration from imported packages.
	mux := http.NewServeMux()

	// Application pages
	mux.HandleFunc("/home", homePage)
	mux.HandleFunc("/courses", coursePage)
	mux.HandleFunc("/about", aboutPage)
	mux.HandleFunc("/contact", contactPage)

	// Kubernetes health check endpoints
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readyz", readyzHandler)

	// Configure the server with explicit timeouts.
	// Without timeouts, a slow or malicious client can hold a connection open
	// indefinitely, eventually exhausting server resources.
	//
	//   ReadTimeout  — maximum time to read the full request (header + body).
	//   WriteTimeout — maximum time to write the full response.
	//   IdleTimeout  — how long a keep-alive connection stays open when idle.
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a background goroutine so the main goroutine is
	// free to listen for OS signals.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server starting on %s", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Block until we receive SIGINT (Ctrl+C) or SIGTERM (from Kubernetes).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		// The server stopped unexpectedly — e.g. the port was already in use.
		log.Fatalf("server error: %v", err)

	case sig := <-quit:
		// A shutdown signal was received. Start graceful shutdown.
		log.Printf("shutdown signal received: %s", sig)

		// Give in-flight requests up to 30 seconds to complete.
		// This matches terminationGracePeriodSeconds in the Kubernetes Deployment.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}

		log.Println("server stopped cleanly")
	}
}

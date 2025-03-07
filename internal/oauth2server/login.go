package oauth2server

import (
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/valkey-io/valkey-go"
)

// LoginPageData holds the data to be passed to the login template
type LoginPageData struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	State        string
	Scope        string
	Error        string
	ErrorMessage string
	PageKey      string
}

// HandleLoginPage renders a styled login form
func (self *Oauth2Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")
	scope := r.URL.Query().Get("scope")
	errorParam := r.URL.Query().Get("error")
	// Unique key specific to the rendered form
	pageKey := uuid.NewString()

	// Create page data
	data := LoginPageData{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		ResponseType: responseType,
		State:        state,
		Scope:        scope,
		Error:        errorParam,
		PageKey:      pageKey,
	}

	// Set error message based on error code
	if errorParam == "invalid_credentials" {
		data.ErrorMessage = "Invalid email or password. Please try again."
	}

	// Create a template with the HTML and CSS
	tmpl := template.Must(template.New("login").Parse(loginTemplate))

	// Store the state in cache to prevent double submission
	self.StringCache.SetWithExpiration(self.Ctx, pageKey, "1", 30*time.Minute)

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// HandleLoginSubmit processes the login form submission
func (self *Oauth2Server) HandleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	clientID := r.Form.Get("client_id")
	redirectURI := r.Form.Get("redirect_uri")
	responseType := r.Form.Get("response_type")
	state := r.Form.Get("state")
	scope := r.Form.Get("scope")
	pageKey := r.Form.Get("page_key")

	// Validate state to prevent double submission
	if _, err := self.StringCache.Getdel(self.Ctx, pageKey); err != nil {
		if err == valkey.Nil {
			// Return ok empty status
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Error("Error validating request: ", err)
		http.Error(w, "Error validating request", http.StatusInternalServerError)
		return
	}

	// Validate credentials against your repository
	user, err := self.Repository.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		// Authentication failed - redirect back to login with error
		// Preserve all original query parameters
		q := url.Values{}
		q.Set("client_id", clientID)
		q.Set("redirect_uri", redirectURI)
		q.Set("response_type", responseType)
		q.Set("state", state)
		q.Set("scope", scope)
		q.Set("error", "invalid_credentials")

		http.Redirect(w, r, "/login?"+q.Encode(), http.StatusFound)
		return
	}

	// Authentication succeeded - redirect back to authorize endpoint with user info
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", responseType)
	q.Set("state", state)
	q.Set("scope", scope)
	q.Set("user_id", user.Email) // Pass the authenticated user ID

	http.Redirect(w, r, "/authorize?"+q.Encode(), http.StatusFound)
}

// Login page HTML template with embedded CSS
const loginTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Sign In - Unbind</title>
    <style>
      :root {
          --background: hsl(120 35% 95%);
          --foreground: hsl(120 90% 5%);
          --primary: hsl(120 90% 5%);
          --destructive: hsl(0 84% 30%);
          --border: hsl(120 20% 85%);
          --input: hsl(120 35% 95%);
          --muted-foreground: hsl(120 10% 35%);
          --top-loader: hsl(120 100% 30%);
          --border-radius: 0.6rem;
      }

      @media (prefers-color-scheme: dark) {
        :root {
          --background: hsl(120 9% 6%);
          --foreground: hsl(120 100% 97%);
          --primary: hsl(120 100% 97%);
          --destructive: hsl(0 80% 69%);
          --border: hsl(120 5% 12%);
          --input: hsl(120 7% 8%);
          --muted-foreground: hsl(120 6% 65%);
          --muted-more-foreground: hsl(120 6% 40%);
          --top-loader: hsl(120 70% 70%);
        }
      }

      * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
      }

      body {
          font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
          background-color: var(--background);
          color: var(--foreground);
          font-size: 1rem;
          line-height: 1.6;
          display: flex;
          align-items: center;
          justify-content: center;
          min-height: 100vh;
          padding: 2rem 1rem calc(2rem + 8vh) 1rem;
          word-wrap: break-word;
      }

      .login-container {
          width: 100%;
          max-width: 20rem;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          gap: 1rem
      }

      .form {
        width: 100%;
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      label {
          display: block;
          padding: 0rem 0.25rem 0.25rem 0.25rem;
          font-size: 0.85rem;
          font-weight: 400;
          color: var(--muted-foreground);
      }

      .logo {
        width: 100%;
        max-width: 8rem;
        height: auto;
      }

      input {
        width: 100%;
        color: var(--foreground);
        background: var(--input);
        padding: 0.75rem 1rem;
        border: 1px solid var(--border);
        border-radius: var(--border-radius);
        font-size: 1rem;
      }

      input::placeholder{
          color: color-mix(in oklab, var(--foreground) 50%, transparent);
      }


      input:focus {
          outline: none;
          border-color: color-mix(in oklab, var(--primary) 50%, transparent);
      }

			button:disabled {
				opacity: 0.7;
				cursor: not-allowed;
			}

      button {
          display: block;
          width: 100%;
          padding: 0.75rem 1rem;
          background-color: var(--primary);
          color: var(--background);
          border: none;
          border-radius: var(--border-radius);
          font-weight: 700;
          font-size: 1rem;
          cursor: pointer;
      }

      button:hover {
          background-color: color-mix(in oklab, var(--primary) 85%, transparent);
      }

      .error-message {
          color: var(--destructive);
          background-color: color-mix(in oklab, var(--destructive) 15%, transparent);
          padding: 0.5rem 0.75rem;
          border-radius: var(--border-radius);
          font-size: 0.85rem;
          width: 100%;
          display: {{if .ErrorMessage}}block{{else}}none{{end}};
      }

      input:-webkit-autofill,
      input:-webkit-autofill:hover,
      input:-webkit-autofill:focus,
      input:-webkit-autofill:active {
        -webkit-box-shadow: 0 0 0 100px var(--top-loader) inset !important;
        -webkit-text-fill-color: color-mix(in oklab, var(--top-loader) 15%, transparent) !important;
        transition: background-color 5000s ease-in-out 0s;
      }
    </style>
		<script>
			document.getElementById('loginForm').addEventListener('submit', function(e) {
				// Disable the submit button
				const submitButton = document.getElementById('submitButton');
				if (submitButton.disabled) {
					e.preventDefault();
					return false;
				}
				
				submitButton.disabled = true;
				submitButton.textContent = 'Signing in...';
				
				return true;
			});
		</script>		
  </head>
  <body>
    <div class="login-container">
      <svg
        class="logo"
        xmlns="http://www.w3.org/2000/svg"
        width="99"
        height="24"
        fill="none"
        viewBox="0 0 99 24"
      >
        <path
          fill="currentColor"
          fillRule="evenodd"
          clipRule="evenodd"
          d="M12 21C5.373 21 0 15.627 0 9V3h16v6a4 4 0 0 1-8 0V7h2v2a2 2 0 1 0 4 0V5H2v4c0 5.523 4.477 10 10 10s10-4.477 10-10V5h-2v4A8 8 0 1 1 4 9V7h2v2a6 6 0 0 0 12 0V3h6v6c0 6.627-5.373 12-12 12m24-3.38q1.26.62 2.8.62t2.82-.62a4.8 4.8 0 0 0 2.06-1.86q.76-1.24.76-3.16V4h-2.7v8.62q0 1.04-.34 1.76-.34.7-1 1.04-.64.34-1.56.34-.9 0-1.56-.34a2.25 2.25 0 0 1-.98-1.04q-.34-.72-.34-1.76V4h-2.7v8.6q0 1.92.74 3.16a4.7 4.7 0 0 0 2 1.86m10.611-9.7V18h2.7v-5.38q0-.84.3-1.42.3-.6.82-.92.54-.32 1.22-.32 1.06 0 1.56.64.52.64.52 1.84V18h2.68v-5.82q0-1.46-.46-2.46-.44-1-1.3-1.52t-2.12-.52q-1.179 0-2.04.52-.84.52-1.3 1.36l-.2-1.64zm16.504 10.12q.62.2 1.4.2 1.4 0 2.48-.68t1.7-1.88q.64-1.2.64-2.7 0-1.54-.64-2.72a4.7 4.7 0 0 0-1.72-1.88q-1.08-.7-2.48-.7-1.18 0-1.98.48-.78.46-1.26 1.16V3.6h-2.7V18h2.4l.3-1.3q.32.44.78.8.48.34 1.08.54m2.16-2.52q-.6.36-1.4.36-.78 0-1.4-.36-.6-.36-.96-1.02a3.4 3.4 0 0 1-.34-1.54q0-.86.34-1.52.36-.66.96-1.02.62-.38 1.4-.38.8 0 1.4.38.62.36.96 1.04.34.66.34 1.52t-.34 1.52a2.6 2.6 0 0 1-.96 1.02m8.746-7.6h-2.7V18h2.7zm-2.54-1.8q.48.42 1.2.42.74 0 1.2-.42.48-.44.48-1.1t-.48-1.08q-.46-.44-1.2-.44-.72 0-1.2.44-.46.42-.46 1.08t.46 1.1m7.248 1.8h-2.38V18h2.7v-5.38q0-.84.3-1.42.3-.6.82-.92.54-.32 1.22-.32 1.06 0 1.56.64.52.64.52 1.84V18h2.68v-5.82q0-1.46-.46-2.46-.44-1-1.3-1.52t-2.12-.52q-1.18 0-2.04.52-.84.52-1.3 1.36zm11.564 9.64q1.1.68 2.48.68.84 0 1.46-.22a3.3 3.3 0 0 0 1.06-.6q.46-.38.76-.8l.3 1.38h2.4V3.6h-2.7v5.62a3.3 3.3 0 0 0-1.32-1.14q-.82-.4-1.9-.4a4.6 4.6 0 0 0-2.5.7 4.83 4.83 0 0 0-1.74 1.88q-.62 1.18-.62 2.72 0 1.5.62 2.7t1.7 1.88m4.52-2.04q-.6.36-1.4.36-.78 0-1.4-.36-.6-.38-.96-1.04-.34-.66-.34-1.54 0-.84.34-1.5.36-.66.98-1.02.62-.38 1.38-.38.8 0 1.4.38.62.36.96 1.02t.34 1.52-.34 1.52-.96 1.04"
        />
      </svg>
      <form id="loginForm" class="form" method="post" action="/login">
        <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}" />
        <input type="hidden" name="client_id" value="{{.ClientID}}" />
        <input type="hidden" name="response_type" value="{{.ResponseType}}" />
        <input type="hidden" name="state" value="{{.State}}" />
        <input type="hidden" name="scope" value="{{.Scope}}" />
				<input type="hidden" name="page_key" value="{{.PageKey}}">
        <div class="input-group">
          <label for="username">Email</label>
          <input
            type="email"
            id="username"
            name="username"
            placeholder="you@email.com"
            required
            autofocus
          />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input type="password" id="password" name="password" placeholder="••••••••" required />
        </div>
        <button id="submitButton" type="submit">Sign in</button>
      </form>
      <div class="error-message">{{.ErrorMessage}}</div>
    </div>
  </body>
</html>`

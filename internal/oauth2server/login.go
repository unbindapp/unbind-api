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
const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign In - Unbind</title>
    <style>
        :root {
            --primary-color: #4a6cf7;
            --primary-hover: #3858d6;
            --dark-text: #212b36;
            --gray-text: #637381;
            --light-gray: #f8fafc;
            --error-color: #ff4757;
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
            background-color: #f9fafb;
            color: var(--dark-text);
            line-height: 1.6;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
        }
        
        .login-container {
            background-color: white;
            border-radius: 10px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
            width: 100%;
            max-width: 400px;
            padding: 40px;
        }
        
        .logo-container {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .logo {
            font-size: 24px;
            font-weight: 700;
            color: var(--primary-color);
        }
        
        h1 {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 24px;
            text-align: center;
        }
        
        .form-group {
            margin-bottom: 24px;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: var(--gray-text);
        }
        
        input[type="email"],
        input[type="password"] {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #e2e8f0;
            border-radius: 6px;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        
        input[type="email"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 3px rgba(74, 108, 247, 0.1);
        }
        
        button {
            display: block;
            width: 100%;
            padding: 12px 16px;
            background-color: var(--primary-color);
            color: white;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        
        button:hover {
            background-color: var(--primary-hover);
        }
        
        .error-message {
            color: var(--error-color);
            background-color: rgba(255, 71, 87, 0.1);
            padding: 12px 16px;
            border-radius: 6px;
            margin-bottom: 24px;
            font-size: 14px;
            display: {{if .ErrorMessage}}block{{else}}none{{end}};
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="logo-container">
            <div class="logo">Unbind</div>
        </div>
        
        <h1>Sign in to your account</h1>
        
        <div class="error-message">
            {{.ErrorMessage}}
        </div>
        
        <form method="post" action="/login">
            <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}">
            <input type="hidden" name="client_id" value="{{.ClientID}}">
            <input type="hidden" name="response_type" value="{{.ResponseType}}">
            <input type="hidden" name="state" value="{{.State}}">
            <input type="hidden" name="scope" value="{{.Scope}}">
						<input type="hidden" name="page_key" value="{{.PageKey}}">
            
            <div class="form-group">
                <label for="username">Email</label>
                <input type="email" id="username" name="username" placeholder="your@email.com" required autofocus>
            </div>
            
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" placeholder="••••••••" required>
            </div>
            
            <button type="submit">Sign in</button>
        </form>
    </div>
</body>
</html>`

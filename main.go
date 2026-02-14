package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const postsAPI = "https://jsonplaceholder.typicode.com/posts"

const (
	btnCSS            = `btn relative inline-flex items-center justify-center w-full py-3 px-5 text-sm font-semibold bg-emerald-400 text-gray-950 rounded-[10px] cursor-pointer transition-all duration-200 hover:brightness-110 hover:shadow-[0_4px_20px_rgba(52,211,153,0.2)] active:scale-[0.98]`
	btnSecondaryCSS   = `btn relative inline-flex items-center justify-center w-full py-3 px-5 text-sm font-semibold text-gray-200 bg-white/[0.03] border border-white/[0.07] backdrop-blur-sm rounded-[10px] cursor-pointer transition-all duration-200 hover:bg-white/[0.07] hover:border-white/[0.14] hover:text-white active:scale-[0.98]`
	btnGhostCSS       = `btn relative inline-flex items-center justify-center w-full py-3 px-5 text-sm font-semibold text-gray-500 bg-transparent border border-white/[0.07] rounded-[10px] cursor-pointer transition-all duration-200 hover:bg-white/[0.05] hover:border-white/[0.14] hover:text-gray-200 active:scale-[0.98]`
	btnDangerCSS      = `btn relative inline-flex items-center justify-center w-full py-3 px-5 text-sm font-semibold bg-red-400 text-white rounded-[10px] cursor-pointer transition-all duration-200 hover:brightness-110 hover:shadow-[0_4px_20px_rgba(248,113,113,0.2)] active:scale-[0.98]`
	spinnerCSS        = `btn-spinner hidden w-[1em] h-[1em] border-2 border-gray-950 border-t-transparent rounded-full animate-spin`
	spinnerLightCSS   = `btn-spinner hidden w-[1em] h-[1em] border-2 border-white border-t-transparent rounded-full animate-spin`
	spinnerSecCSS     = `btn-spinner hidden w-[1em] h-[1em] border-2 border-gray-200 border-t-transparent rounded-full animate-spin`
	labelCSS          = `block text-[0.65rem] font-semibold uppercase tracking-wider text-gray-500 mb-2`
	inputCSS          = `w-full px-3.5 py-2.5 mb-4 font-mono text-sm text-gray-200 bg-[#12151b] border border-white/[0.07] rounded-[10px] transition-all duration-200 placeholder:text-gray-500 focus:outline-none focus:border-emerald-400 focus:ring-2 focus:ring-emerald-400/10`
	radioOptionCSS    = `radio-option flex items-center gap-3 px-3.5 py-2.5 bg-[#12151b] border border-white/[0.07] rounded-[10px] cursor-pointer transition-all duration-200 hover:border-white/[0.14] hover:bg-white/[0.04] has-[:checked]:border-emerald-400 has-[:checked]:bg-emerald-400/10`
	radioInputCSS     = `radio-dot appearance-none w-4 h-4 border-2 border-gray-500 rounded-full cursor-pointer transition-all duration-150 shrink-0 checked:border-emerald-400 checked:bg-emerald-400 checked:shadow-[inset_0_0_0_3px_#12151b]`
)

// ---------------------------------------------------------------------------
// Mock token data
// ---------------------------------------------------------------------------

type mockToken struct {
	ID     string
	Type   string
	Name   string
	Masked string
}

var mockTokens = []mockToken{
	{ID: "tok_1", Type: "GitLab PAT", Name: "gitlab-deploy-key", Masked: "glpat-****Xk9f"},
	{ID: "tok_2", Type: "JIRA", Name: "jira-automation", Masked: "jira-****m3Qz"},
	{ID: "tok_3", Type: "GitLab PAT", Name: "gitlab-ci-runner", Masked: "glpat-****Lw2d"},
}

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

var httpClient = &http.Client{Timeout: 10 * time.Second}

func callPlaceholderPost(body map[string]interface{}) bool {
	data, err := json.Marshal(body)
	if err != nil {
		return false
	}
	resp, err := httpClient.Post(postsAPI, "application/json", bytes.NewReader(data))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func callPlaceholderDelete() bool {
	req, err := http.NewRequest(http.MethodDelete, postsAPI+"/1", nil)
	if err != nil {
		return false
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// ---------------------------------------------------------------------------
// HTML builders
// ---------------------------------------------------------------------------

func actionsAreaHTML(message string, isError bool) string {
	msgHTML := ""
	if message != "" {
		if isError {
			msgHTML = fmt.Sprintf(`<p class="text-sm text-center px-3 py-2.5 rounded-[10px] text-red-400 bg-red-400/[0.12] mb-1">%s</p>`, message)
		} else {
			msgHTML = fmt.Sprintf(`<p class="text-sm text-center px-3 py-2.5 rounded-[10px] text-emerald-400 bg-emerald-400/10 mb-1">%s</p>`, message)
		}
	}
	return fmt.Sprintf(`<div id="connect-area" class="flex flex-col gap-2 mt-1">
        %s
        <button class="%s" type="button"
            hx-get="./add-key-form"
            hx-target="#connect-area"
            hx-swap="outerHTML"
            hx-indicator="#add-spinner"
            hx-disabled-elt="this">
            <span class="btn-text">Add Wallet Key</span>
            <span id="add-spinner" class="%s" aria-hidden="true"></span>
        </button>
        <div class="flex gap-2">
            <button class="%s" type="button"
                hx-get="./update-key-form"
                hx-target="#connect-area"
                hx-swap="outerHTML"
                hx-indicator="#update-spinner"
                hx-disabled-elt="this">
                <span class="btn-text">Update key</span>
                <span id="update-spinner" class="%s" aria-hidden="true"></span>
            </button>
            <button class="%s" type="button"
                hx-get="./delete-key-form"
                hx-target="#connect-area"
                hx-swap="outerHTML"
                hx-indicator="#delete-spinner"
                hx-disabled-elt="this">
                <span class="btn-text">Delete key</span>
                <span id="delete-spinner" class="%s" aria-hidden="true"></span>
            </button>
        </div>
    </div>`, msgHTML, btnCSS, spinnerCSS, btnSecondaryCSS, spinnerSecCSS, btnSecondaryCSS, spinnerSecCSS)
}

func addKeyFormHTML(isError bool) string {
	errHTML := ""
	if isError {
		errHTML = `<p class="text-[0.78rem] text-red-400 bg-red-400/[0.12] px-3 py-2 rounded-[10px] mb-3">Failed to add key. Try again.</p>`
	}
	return fmt.Sprintf(`<div id="connect-area" class="mt-1">
        %s
        <form class="text-left" id="add-key-form" hx-post="./add-key" hx-target="#connect-area" hx-swap="outerHTML" hx-indicator="#submit-spinner" hx-disabled-elt="find button[type='submit']">
            <label class="%s">Token type</label>
            <div class="flex flex-col gap-1.5 mb-4">
                <label class="%s">
                    <input class="%s" type="radio" name="token_type" value="gitlab_pat" required />
                    <span class="text-sm font-medium text-gray-200">GitLab PAT token</span>
                </label>
                <label class="%s">
                    <input class="%s" type="radio" name="token_type" value="jira" />
                    <span class="text-sm font-medium text-gray-200">JIRA token</span>
                </label>
            </div>
            <div id="add-key-input-wrap" class="hidden">
                <label class="%s" for="wallet-key">Key</label>
                <input class="%s" id="wallet-key" name="key" type="text" placeholder="Enter your token" autocomplete="off" required />
            </div>
            <div class="flex gap-2 mt-2">
                <button class="%s flex-1" type="button" hx-get="./actions" hx-target="#connect-area" hx-swap="outerHTML">
                    Cancel
                </button>
                <button class="%s flex-1" type="submit">
                    <span class="btn-text">Save key</span>
                    <span id="submit-spinner" class="%s" aria-hidden="true"></span>
                </button>
            </div>
        </form>
    </div>`, errHTML, labelCSS, radioOptionCSS, radioInputCSS, radioOptionCSS, radioInputCSS, labelCSS, inputCSS, btnGhostCSS, btnCSS, spinnerCSS)
}

func tokenRadiosHTML() string {
	var sb strings.Builder
	for _, t := range mockTokens {
		fmt.Fprintf(&sb, `<label class="%s">
                    <input class="%s" type="radio" name="token_id" value="%s" required />
                    <span>
                        <span class="block text-sm font-medium text-gray-200">%s</span>
                        <span class="block font-mono text-[0.7rem] text-gray-500 mt-0.5">%s &middot; %s</span>
                    </span>
                </label>
`, radioOptionCSS, radioInputCSS, t.ID, t.Name, t.Type, t.Masked)
	}
	return sb.String()
}

func updateKeyFormHTML(isError bool) string {
	errHTML := ""
	if isError {
		errHTML = `<p class="text-[0.78rem] text-red-400 bg-red-400/[0.12] px-3 py-2 rounded-[10px] mb-3">Update failed. Try again.</p>`
	}
	return fmt.Sprintf(`<div id="connect-area" class="mt-1">
        %s
        <form class="text-left" id="update-key-form" hx-post="./update-key" hx-target="#connect-area" hx-swap="outerHTML" hx-indicator="#update-submit-spinner" hx-disabled-elt="find button[type='submit']">
            <label class="%s">Select token to update</label>
            <div class="flex flex-col gap-1.5 mb-4">
                %s
            </div>
            <div id="update-new-key-wrap" class="hidden">
                <label class="%s" for="new-key">New key</label>
                <input class="%s" id="new-key" name="key" type="text" placeholder="Enter new token value" autocomplete="off" required />
            </div>
            <div class="flex gap-2 mt-2">
                <button class="%s flex-1" type="button" hx-get="./actions" hx-target="#connect-area" hx-swap="outerHTML">
                    Cancel
                </button>
                <button class="%s flex-1" type="submit">
                    <span class="btn-text">Save</span>
                    <span id="update-submit-spinner" class="%s" aria-hidden="true"></span>
                </button>
            </div>
        </form>
    </div>`, errHTML, labelCSS, tokenRadiosHTML(), labelCSS, inputCSS, btnGhostCSS, btnCSS, spinnerCSS)
}

func deleteKeyFormHTML(message string, isError bool) string {
	msgHTML := ""
	if message != "" {
		if isError {
			msgHTML = fmt.Sprintf(`<p class="text-[0.78rem] text-red-400 bg-red-400/[0.12] px-3 py-2 rounded-[10px] mb-3">%s</p>`, message)
		} else {
			msgHTML = fmt.Sprintf(`<p class="text-sm text-center px-3 py-2.5 rounded-[10px] text-emerald-400 bg-emerald-400/10 mb-1">%s</p>`, message)
		}
	}
	return fmt.Sprintf(`<div id="connect-area" class="mt-1">
        %s
        <form class="text-left" id="delete-key-form">
            <label class="%s">Select token to delete</label>
            <div class="flex flex-col gap-1.5 mb-4">
                %s
            </div>
            <div class="flex gap-2 mt-2">
                <button class="%s flex-1" type="button" hx-get="./actions" hx-target="#connect-area" hx-swap="outerHTML">
                    Cancel
                </button>
                <button class="%s flex-1" type="button"
                    hx-post="./delete-key"
                    hx-target="#connect-area"
                    hx-swap="outerHTML"
                    hx-include="#delete-key-form"
                    hx-confirm="Are you sure you want to delete this token? This action cannot be undone."
                    hx-disabled-elt="this"
                    hx-indicator="#delete-confirm-spinner">
                    <span class="btn-text">Delete</span>
                    <span id="delete-confirm-spinner" class="%s" aria-hidden="true"></span>
                </button>
            </div>
        </form>
    </div>`, msgHTML, labelCSS, tokenRadiosHTML(), btnGhostCSS, btnDangerCSS, spinnerLightCSS)
}

func indexHTML() string {
	gitEmail := os.Getenv("GIT_COMMITTER_EMAIL")
	emailCls := ""
	emailVal := gitEmail
	if gitEmail == "" {
		emailCls = "text-gray-500 italic"
		emailVal = "Not set"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RDX Wallet</title>
    <link rel="stylesheet" href="./vendor/css/fonts.css">
    <script src="./vendor/js/tailwindcss.js"></script>
    <script>
    tailwind.config = {
        theme: {
            extend: {
                fontFamily: {
                    sans: ['Outfit', 'system-ui', '-apple-system', 'sans-serif'],
                    mono: ['JetBrains Mono', 'ui-monospace', 'monospace'],
                },
            },
        },
    }
    </script>
    <script src="./vendor/js/htmx-2.0.4.min.js"></script>
    <script>
    document.addEventListener("htmx:beforeSwap", function(e) {
        if (e.detail.xhr.status === 422) { e.detail.shouldSwap = true; }
    });
    document.addEventListener("change", function(e) {
        if (e.target.matches('#add-key-form input[name="token_type"]')) {
            document.getElementById('add-key-input-wrap').classList.remove('hidden');
            document.getElementById('wallet-key').required = true;
        }
        if (e.target.matches('#update-key-form input[name="token_id"]')) {
            document.getElementById('update-new-key-wrap').classList.remove('hidden');
            document.getElementById('new-key').required = true;
        }
    });
    </script>
    <style>
        body {
            background-image: radial-gradient(ellipse 60%% 40%% at 50%% -10%%, rgba(52,211,153,0.08), transparent);
        }
        /* htmx spinner states */
        .btn.htmx-request .btn-text { visibility: hidden; }
        .btn.htmx-request .btn-spinner { display: inline-block; position: absolute; }
        form.htmx-request .btn .btn-text { visibility: hidden; }
        form.htmx-request .btn .btn-spinner { display: inline-block; position: absolute; }
    </style>
</head>
<body class="min-h-screen flex items-center justify-center p-6 bg-[#090b0f] text-gray-200 font-sans">
    <div class="w-full max-w-[400px]">
        <h1 class="text-[1.75rem] font-bold tracking-tight text-center text-white mb-0.5">RDX Wallet</h1>
        <p class="text-sm text-gray-500 text-center mb-8">Secure. Simple. Yours.</p>

        <div class="bg-white/[0.03] backdrop-blur-xl border border-white/[0.07] rounded-[14px] px-6 py-5 text-left mb-3">
            <div class="text-[0.65rem] font-semibold uppercase tracking-wider text-gray-500 mb-2">Email</div>
            <div class="font-mono text-sm break-all %s">%s</div>
        </div>

        %s

        <p class="text-xs text-gray-500 text-center mt-5 leading-relaxed">Auto-deletion policy: All tokens are automatically deleted after 30 days. This is non-adjustable.</p>
    </div>
</body>
</html>`, emailCls, emailVal, actionsAreaHTML("", false))
}

// ---------------------------------------------------------------------------
// Route handlers
// ---------------------------------------------------------------------------

func handleIndex(c echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML())
}

func handleActions(c echo.Context) error {
	return c.HTML(http.StatusOK, actionsAreaHTML("", false))
}

func handleAddKeyForm(c echo.Context) error {
	return c.HTML(http.StatusOK, addKeyFormHTML(false))
}

func handleUpdateKeyForm(c echo.Context) error {
	return c.HTML(http.StatusOK, updateKeyFormHTML(false))
}

func handleDeleteKeyForm(c echo.Context) error {
	return c.HTML(http.StatusOK, deleteKeyFormHTML("", false))
}

func handleAddKey(c echo.Context) error {
	tokenType := c.FormValue("token_type")
	key := c.FormValue("key")
	ok := callPlaceholderPost(map[string]interface{}{
		"title":  tokenType,
		"body":   key,
		"userId": 1,
	})
	if ok {
		return c.HTML(http.StatusOK, actionsAreaHTML("Key added.", false))
	}
	return c.HTML(http.StatusUnprocessableEntity, addKeyFormHTML(true))
}

func handleUpdateKey(c echo.Context) error {
	tokenID := c.FormValue("token_id")
	key := c.FormValue("key")
	ok := callPlaceholderPost(map[string]interface{}{
		"title":   "update_key",
		"body":    key,
		"tokenId": tokenID,
	})
	if ok {
		return c.HTML(http.StatusOK, actionsAreaHTML("Key updated.", false))
	}
	return c.HTML(http.StatusUnprocessableEntity, updateKeyFormHTML(true))
}

func handleDeleteKey(c echo.Context) error {
	ok := callPlaceholderDelete()
	if ok {
		return c.HTML(http.StatusOK, actionsAreaHTML("Token deleted.", false))
	}
	return c.HTML(http.StatusUnprocessableEntity, deleteKeyFormHTML("Delete failed. Try again.", true))
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/vendor", "../vendor")

	e.GET("/", handleIndex)
	e.GET("/actions", handleActions)
	e.GET("/add-key-form", handleAddKeyForm)
	e.GET("/update-key-form", handleUpdateKeyForm)
	e.GET("/delete-key-form", handleDeleteKeyForm)
	e.POST("/add-key", handleAddKey)
	e.POST("/update-key", handleUpdateKey)
	e.POST("/delete-key", handleDeleteKey)

	e.Logger.Fatal(e.Start(":8000"))
}

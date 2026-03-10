package confluence

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/pages/123" {
			t.Errorf("Expected to request '/wiki/api/v2/pages/123', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "123",
			"title": "Test Page",
			"body": map[string]interface{}{
				"storage": map[string]interface{}{
					"value": "<p>Hello World</p>",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	page, err := client.GetPage("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if page.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", page.Title)
	}
}

func TestSearchPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"id":    "123",
					"title": "Found Page",
					"space": map[string]interface{}{
						"key":  "test",
						"name": "Test Space",
					},
					"excerpt": "Excerpt...",
					"_links": map[string]interface{}{
						"webui": "/spaces/test/pages/123",
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	pages, err := client.SearchPages("text~\"test\"", 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 || pages[0].Title != "Found Page" {
		t.Errorf("Unexpected search results: %+v", pages)
	}
}

func TestCreatePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "456",
			"_links": map[string]interface{}{
				"webui": "/spaces/test/pages/456",
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	res, err := client.CreatePage("111", "", "New Page", "<h1>Content</h1>")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.ID != "456" {
		t.Errorf("Expected ID 456, got '%s'", res.ID)
	}
}

func TestGetPageChildren(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"id":    "789",
					"title": "Child Page",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	pages, err := client.GetPageChildren("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 || pages[0].Title != "Child Page" {
		t.Errorf("Unexpected children results: %+v", pages)
	}
}

func TestGetSpaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"key":  "DEV",
					"name": "Development",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	spaces, err := client.GetSpaces(10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(spaces) != 1 || spaces[0].Key != "DEV" {
		t.Errorf("Unexpected space results: %+v", spaces)
	}
}

func TestUpdatePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "123",
				"title": "Old Page",
				"version": map[string]interface{}{
					"number": 1,
				},
			})
			return
		}

		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "123",
				"title": "Updated Page",
				"version": map[string]interface{}{
					"number": 2,
				},
				"_links": map[string]interface{}{
					"webui": "/spaces/test/pages/123",
				},
			})
			return
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	res, err := client.UpdatePage("123", "Updated Page", "<h1>New Content</h1>")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.Title != "Updated Page" {
		t.Errorf("Expected Title 'Updated Page', got '%s'", res.Title)
	}
}

func TestAddComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "999",
			"_links": map[string]interface{}{
				"webui": "/spaces/test/pages/123/comments/999",
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	res, err := client.AddComment("123", "Test Comment")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.ID != "999" {
		t.Errorf("Expected ID 999, got '%s'", res.ID)
	}
}

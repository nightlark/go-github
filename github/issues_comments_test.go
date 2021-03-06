// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestIssuesService_ListComments_allIssues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		testFormValues(t, r, values{
			"since": "2002-02-10T15:30:00Z",
			"page":  "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	since := time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC)
	opt := &IssueListCommentsOptions{
		Since:       &since,
		ListOptions: ListOptions{Page: 2},
	}
	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, "o", "r", 0, opt)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []*IssueComment{{ID: Int64(1)}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}
}

func TestIssuesService_ListComments_specificIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `[{"id":1}]`)
	})

	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []*IssueComment{{ID: Int64(1)}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}
}

func TestIssuesService_ListComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListComments(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_GetComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeReactionsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.GetComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetComment returned error: %v", err)
	}

	want := &IssueComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.GetComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_GetComment_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.GetComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssuesService_CreateComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &IssueComment{Body: String("b")}

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.CreateComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.CreateComment returned error: %v", err)
	}

	want := &IssueComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.CreateComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_CreateComment_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.CreateComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_EditComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &IssueComment{Body: String("b")}

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	comment, _, err := client.Issues.EditComment(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.EditComment returned error: %v", err)
	}

	want := &IssueComment{ID: Int64(1)}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.EditComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_EditComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.EditComment(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_DeleteComment(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Issues.DeleteComment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.DeleteComments returned error: %v", err)
	}
}

func TestIssuesService_DeleteComment_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Issues.DeleteComment(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

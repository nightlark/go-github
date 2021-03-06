// Copyright 2016 The go-github AUTHORS. All rights reserved.
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
)

func TestReviewers_marshall(t *testing.T) {
	testJSONMarshal(t, &Reviewers{}, "{}")

	u := &Reviewers{
		Users: []*User{{
			Login:       String("l"),
			ID:          Int64(1),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		}},
		Teams: []*Team{{
			ID:              Int64(1),
			NodeID:          String("node"),
			Name:            String("n"),
			Description:     String("d"),
			URL:             String("u"),
			Slug:            String("s"),
			Permission:      String("p"),
			Privacy:         String("priv"),
			MembersCount:    Int(1),
			ReposCount:      Int(1),
			Organization:    nil,
			MembersURL:      String("m"),
			RepositoriesURL: String("r"),
			Parent:          nil,
			LDAPDN:          String("l"),
		}},
	}

	want := `{
		"users" : [
			{
				"login": "l",
				"id": 1,
				"avatar_url": "a",
				"gravatar_id": "g",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e",
				"hireable": true,
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"url": "u"
			}
		],
		"teams" : [
			{
				"id": 1,
				"node_id": "node",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "priv",
				"members_count": 1,
				"repos_count": 1,
				"members_url": "m",
				"repositories_url": "r",
				"ldap_dn": "l"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestPullRequestsService_ListReviews(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	reviews, _, err := client.PullRequests.ListReviews(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("PullRequests.ListReviews returned error: %v", err)
	}

	want := []*PullRequestReview{
		{ID: Int64(1)},
		{ID: Int64(2)},
	}
	if !reflect.DeepEqual(reviews, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", reviews, want)
	}
}

func TestPullRequestsService_ListReviews_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviews(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_GetReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.GetReview(ctx, "o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.GetReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.GetReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_GetReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.GetReview(ctx, "%", "r", 1, 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_DeletePendingReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.DeletePendingReview(ctx, "o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.DeletePendingReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.DeletePendingReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_DeletePendingReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.DeletePendingReview(ctx, "%", "r", 1, 1)
	testURLParseError(t, err)
}

func TestPullRequestsService_ListReviewComments(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	ctx := context.Background()
	comments, _, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, nil)
	if err != nil {
		t.Errorf("PullRequests.ListReviewComments returned error: %v", err)
	}

	want := []*PullRequestComment{
		{ID: Int64(1)},
		{ID: Int64(2)},
	}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("PullRequests.ListReviewComments returned %+v, want %+v", comments, want)
	}
}

func TestPullRequestsService_ListReviewComments_withOptions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `[]`)
	})

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviewComments(ctx, "o", "r", 1, 1, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("PullRequests.ListReviewComments returned error: %v", err)
	}
}

func TestPullRequestReviewRequest_isComfortFadePreview(t *testing.T) {
	path := "path/to/file.go"
	body := "this is a comment body"
	left, right := "LEFT", "RIGHT"
	pos1, pos2, pos3 := 1, 2, 3
	line1, line2, line3 := 11, 22, 33

	tests := []struct {
		name     string
		review   *PullRequestReviewRequest
		wantErr  error
		wantBool bool
	}{{
		name:     "empty review",
		review:   &PullRequestReviewRequest{},
		wantBool: false,
	}, {
		name: "old-style review",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos2,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos3,
			}},
		},
		wantBool: false,
	}, {
		name: "new-style review",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}, {
				Path: &path,
				Body: &body,
				Side: &left,
				Line: &line2,
			}, {
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line3,
			}},
		},
		wantBool: true,
	}, {
		name: "blended comment",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1, // can't have both styles.
				Side:     &right,
				Line:     &line1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}, {
		name: "position then line",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}, {
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}, {
		name: "line then position",
		review: &PullRequestReviewRequest{
			Comments: []*DraftReviewComment{{
				Path: &path,
				Body: &body,
				Side: &right,
				Line: &line1,
			}, {
				Path:     &path,
				Body:     &body,
				Position: &pos1,
			}},
		},
		wantErr: ErrMixedCommentStyles,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotBool, gotErr := tc.review.isComfortFadePreview()
			if tc.wantErr != nil {
				if gotErr != tc.wantErr {
					t.Errorf("isComfortFadePreview() = %v, wanted %v", gotErr, tc.wantErr)
				}
			} else {
				if gotBool != tc.wantBool {
					t.Errorf("isComfortFadePreview() = %v, wanted %v", gotBool, tc.wantBool)
				}
			}
		})
	}
}

func TestPullRequestsService_ListReviewComments_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviewComments(ctx, "%", "r", 1, 1, nil)
	testURLParseError(t, err)
}

func TestPullRequestsService_CreateReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewRequest{
		CommitID: String("commit_id"),
		Body:     String("b"),
		Event:    String("APPROVE"),
	}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.CreateReview(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("PullRequests.CreateReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.CreateReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_CreateReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.CreateReview(ctx, "%", "r", 1, &PullRequestReviewRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_UpdateReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.PullRequests.UpdateReview(ctx, "o", "r", 1, 1, "updated_body")
	if err != nil {
		t.Errorf("PullRequests.UpdateReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PullRequests.UpdateReview = %+v, want %+v", got, want)
	}
}

func TestPullRequestsService_SubmitReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewRequest{
		Body:  String("b"),
		Event: String("APPROVE"),
	}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/events", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.SubmitReview(ctx, "o", "r", 1, 1, input)
	if err != nil {
		t.Errorf("PullRequests.SubmitReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.SubmitReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_SubmitReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.SubmitReview(ctx, "%", "r", 1, 1, &PullRequestReviewRequest{})
	testURLParseError(t, err)
}

func TestPullRequestsService_DismissReview(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &PullRequestReviewDismissalRequest{Message: String("m")}

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/dismissals", func(w http.ResponseWriter, r *http.Request) {
		v := new(PullRequestReviewDismissalRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	review, _, err := client.PullRequests.DismissReview(ctx, "o", "r", 1, 1, input)
	if err != nil {
		t.Errorf("PullRequests.DismissReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int64(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.DismissReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_DismissReview_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.PullRequests.DismissReview(ctx, "%", "r", 1, 1, &PullRequestReviewDismissalRequest{})
	testURLParseError(t, err)
}

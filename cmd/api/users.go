package main

import (
	"context"
	"errors"
	"github.com/Kser1x/social/internal/store"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	result, err := app.storage.Users.DeleteByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateUserPayload struct {
	Username *string `json:"username" validate:"omitempty,max=100"`
	Email    *string `json:"email" validate:"omitempty,max=1000"`
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var payload UpdateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Username != nil {
		user.Username = *payload.Username
	}

	if payload.Email != nil {
		user.Email = *payload.Email
	}

	if err := app.storage.Users.Update(r.Context(), user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.storage.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)

			}
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
func getUserFromCtx(r *http.Request) *store.UserModel {
	user, _ := r.Context().Value(userCtx).(*store.UserModel)
	return user

}

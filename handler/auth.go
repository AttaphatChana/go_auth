package handlers

//
//import (
//	"encoding/json"
//	"net/http"
//
//	"github.com/go-redis/redis/v8"
//	"github.com/jackc/pgx/v5/pgxpool"
//
//	_ "test_auth/models"
//	_ "test_auth/utils"
//)
//
// request shapes
//type creds struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//// SignupHandler registers new users.
//func SignupHandler(db *pgxpool.Pool, rdb *redis.Client) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var c creds
//		json.NewDecoder(r.Body).Decode(&c)
//		// TODO: hash, store in Postgres
//		w.WriteHeader(http.StatusNotImplemented)
//	}
//}
//
// LoginHandler authenticates and issues JWT.
//func LoginHandler(db *pgxpool.Pool, rdb *redis.Client) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var c creds
//		json.NewDecoder(r.Body).Decode(&c)
//		// TODO: verify, create token via utils.GenerateJWT
//		w.WriteHeader(http.StatusNotImplemented)
//	}
//}
//
// ProfileHandler returns user info if authenticated.
//func ProfileHandler(db *pgxpool.Pool) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// TODO: get user ID from context, fetch from Postgres
//		w.WriteHeader(http.StatusNotImplemented)
//	}
//}

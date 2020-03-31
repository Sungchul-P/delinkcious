package service

import (
	"errors"
	"log"      // 로깅용
	"net/http" // HTTP 서비스

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	sgm "github.com/sungchul-p/delinkcious/pkg/social_graph_manager"
)

var (
	BadRoutingError = errors.New("inconsistent mapping between route and handler")
)

func Run() {
	// 소셜 그래프 매니저를 위한 데이터 저장소 생성
	store, err := sgm.NewDbSocialGraphStore("localhost", 5432, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}
	svc, err := sgm.NewSocialGraphManager(store)
	if err != nil {
		log.Fatal(err)
	}

	// 각 엔드포인트에 대한 핸들러 구성
	followHandler := httptransport.NewServer(
		makeFollowEndpoint(svc), // Endpoint 팩토리 함수
		decodeFollowRequest,     // request 디코더
		encodeResponse,          // response 인코더
	)
	unfollowHandler := httptransport.NewServer(
		makeUnfollowEndpoint(svc),
		decodeUnfollowRequest,
		encodeResponse,
	)

	getFollowingHandler := httptransport.NewServer(
		makeGetFollowingEndpoint(svc),
		decodeGetFollowingRequest,
		encodeResponse,
	)

	getFollowersHandler := httptransport.NewServer(
		makeGetFollowersEndpoint(svc),
		decodeGetFollowersRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/follow").Handler(followHandler)
	r.Methods("POST").Path("/unfollow").Handler(unfollowHandler)
	r.Methods("GET").Path("/following/{username}").Handler(getFollowingHandler)
	r.Methods("GET").Path("/followers/{username}").Handler(getFollowersHandler)

	log.Println("Listening on port 9090...")
	log.Fatal(http.ListenAndServe(":9090", r))
}

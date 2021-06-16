package server

func (s *server) loadRoutes() {

	s.router.Handle(
		"/posts",
		adapt(
			s.PostsCreate(),
			s.logging(),
			s.addHeader("Content-Type", "application/json"),
		),
	).Methods("POST")

	s.router.Handle(
		"/posts/{postID}",
		adapt(
			s.PostsGet(),
			s.logging(),
			s.addHeader("Content-Type", "application/json"),
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{offset}/{limit}",
		adapt(
			s.PostsGetAll(),
			s.logging(),
			s.addHeader("Content-Type", "application/json"),
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{postID}",
		adapt(
			s.PostsDelete(),
			s.logging(),
			s.addHeader("Content-Type", "application/json"),
		),
	).Methods("DELETE")
}

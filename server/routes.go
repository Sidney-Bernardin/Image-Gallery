package server

func (s *server) loadRoutes() {

	s.router.Handle(
		"/posts",
		adapt(
			s.PostsCreate(),
			s.logging(),
		),
	).Methods("POST")

	s.router.Handle(
		"/posts/{id}",
		adapt(
			s.PostsGet(),
			s.logging(),
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{offset}/{limit}",
		adapt(
			s.PostsGetAll(),
			s.logging(),
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{id}",
		adapt(
			s.PostsDelete(),
			s.logging(),
		),
	).Methods("DELETE")
}

package server

func (s *server) loadRoutes() {

	contentTypeJSON := s.addHeader("Content-Type", "application/json")

	s.router.Handle(
		"/posts",
		adapt(
			s.PostsCreate(),
			s.logging(),
			contentTypeJSON,
		),
	).Methods("POST")

	s.router.Handle(
		"/posts/{postID}",
		adapt(
			s.PostsGet(),
			s.logging(),
			contentTypeJSON,
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{offset}/{limit}",
		adapt(
			s.PostsGetAll(),
			s.logging(),
			contentTypeJSON,
		),
	).Methods("GET")

	s.router.Handle(
		"/posts/{postID}",
		adapt(
			s.PostsDelete(),
			s.logging(),
			contentTypeJSON,
		),
	).Methods("DELETE")

	s.router.Handle(
		"/poststhumbnail/{postID}",
		adapt(
			s.PostsThumbnailGet(),
			s.logging(),
		),
	).Methods("GET")
}

package handler

type Handlers struct {
	WebpageHandler *WebpageHandler
}

func Setup() *Handlers {
	return &Handlers{
		WebpageHandler: NewWebpageHandler(),
	}
}

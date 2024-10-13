package can

type Result struct {
	Allowed bool
	Reason  string
	Error   error
}

func (cr Result) Allow(f func()) Result {
	if cr.Allowed {
		f()
	}

	return cr
}

func (cr Result) Err(f func(error)) Result {
	if cr.Error != nil {
		f(cr.Error)
	}

	return cr
}

func (cr Result) Deny(f func(string)) Result {
	if !cr.Allowed && cr.Error == nil {
		f(cr.Reason)
	}

	return cr
}

func Allowed() Result {
	return Result{Allowed: true}
}

func Denied(reason string) Result {
	return Result{Allowed: false, Reason: reason}
}

func Error(err error) Result {
	return Result{Allowed: false, Error: err}
}

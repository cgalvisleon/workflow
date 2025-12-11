package workflow

type SetFn func(*Instance) error
type GetFn func(string) (*Instance, error)
type DeleteFn func(string) error

var (
	getFn    GetFn
	setFn    SetFn
	deleteFn DeleteFn
)

/**
* OnGet
* @param f GetFn
* @return void
**/
func OnGet(f GetFn) {
	if f == nil {
		return
	}

	getFn = f
}

/**
* OnSet
* @param f SetFn
* @return void
**/
func OnSet(f SetFn) {
	if f == nil {
		return
	}

	setFn = f
}

/**
* OnDelete
* @param f DeleteFn
* @return void
**/
func OnDelete(f DeleteFn) {
	deleteFn = f
}

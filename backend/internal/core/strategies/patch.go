package strategies

type Patch interface {
	ApplyTo(task Task) error
}

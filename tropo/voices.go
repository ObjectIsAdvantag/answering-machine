package tropo


// TODO : dynamically generate the list of voices for each env : Dev / Prod by calling Tropo REST API

type Voice struct {
	Env	 	environment
	Lang	language
	name 	string
}

type environment int // unexported type users can't construct their own.
const (
	DEV environment = iota
	PROD
)

type language string

const (
	FR language = "FR"
)


var AUDREY = Voice{DEV, FR, "Audrey"}



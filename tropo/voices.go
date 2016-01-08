package tropo


type Voice struct {
	name 	string
	Lang	language
	gender	gender
	Env	 	environment
}

func (voice *Voice) Name() string {
	return voice.name
}

type environment int // unexported type users can't construct their own.
const (
	DEV environment = iota
	PROD
)

type language string
const (
	fr_FR language = "fr_FR"
)

type gender int // unexported type users can't construct their own.
const (
	FEMALE gender = iota
	MALE
)


// TODO : extend the list of voices for each env : Dev / Prod
// see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
var AUDREY = Voice{"Audrey", fr_FR, FEMALE, DEV}



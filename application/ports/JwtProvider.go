package ports

type JwtProvider interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error) // returns userName or an error
}

package rpc

type LoginRequest struct {
	ID int
}

type LoginResponse struct {
	Session string
}

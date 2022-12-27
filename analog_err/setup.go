package analog_err

func (er *ErrorRegistry) Setup() {
	analogErrors := []*AnalogError{{
		Code:    001,
		Message: "Invalid request payload",
		Detail:  "The payload received does not match the requirements to be successfully decoded",
		IsError: true,
	}, {
		Code:    002,
		Message: "Missing Pseudo or Password for given user",
		Detail:  "Ensure the payload received contains 'pseudo' and 'password' keys",
		IsError: true,
	}, {
		Code:    003,
		Message: "Password encryption failed",
		Detail:  "For security reason, users password are encrypted in database, it seems your password failed to be encrypted. Please try with another password",
		IsError: true,
	}, {
		Code:    004,
		Message: "User account creation failed",
		Detail:  "The account creation failed, please try again ",
		IsError: true,
	}, {
		Code:    005,
		Message: "User account not found",
		Detail:  "The account requested does not exist in database",
		IsError: true,
	}, {
		Code:    006,
		Message: "User pseudo or password is invalid",
		Detail:  "Authentication failed, retype your password to be sure, and check again for your pseudo",
		IsError: true,
	}, {
		Code:    007,
		Message: "Authentication token generation failed",
		Detail:  "An error occurred while creating your authentication token. We will investigate on this error",
		IsError: true,
	}, {
		Code:    010,
		Message: "Token not found in request header",
		Detail:  "You must log in with your account and then navigates. If you do not have created account please create one.",
		IsError: true,
	}, {
		Code:    011,
		Message: "Token has expired",
		Detail:  "Token has an expiration limit, refresh your token before continuing",
		IsError: true,
	}, {
		Code:    012,
		Message: "Unauthorized action",
		Detail:  "Token is not authorized, please try to logout login.",
		IsError: true,
	}, {
		Code:    013,
		Message: "Camera not found",
		Detail:  "The requested Camera does not exist in database",
		IsError: true,
	}, {
		Code:    014,
		Message: "User camera creation failed",
		Detail:  "The user camera creation failed, please try again",
		IsError: true,
	}, {
		Code:    015,
		Message: "Action not authorized, this is not your resource",
		Detail:  "The current request tries to update a resource of another user, event reported",
		IsError: true,
	}, {
		Code:    016,
		Message: "User camera deletion failed",
		Detail:  "The user camera deletion failed, please try again",
		IsError: true,
	}, {
		Code:    017,
		Message: "User Camera not found",
		Detail:  "The requested User Camera does not exist in database",
		IsError: true,
	}}

	for _, aErr := range analogErrors {
		er.Register(aErr)
	}
}

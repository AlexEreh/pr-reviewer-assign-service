package api

import (
	gohttp "net/http"

	"pr-reviewer-assign-service/pkg/errors"
	"pr-reviewer-assign-service/pkg/errors/http"
)

// ContractError нужна для корректной генерации Swagger.
type ContractError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var ErrUserIDNotProvided = errors.Template{
	Code:    "NO_USER_ID",
	Message: "no user_id provided",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusBadRequest),
	},
}

var ErrTeamNameNotProvided = errors.Template{
	Code:    "NO_TEAM_NAME",
	Message: "no team_name provided",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusBadRequest),
	},
}

var ErrTeamExists = errors.Template{
	Code:    "TEAM_EXISTS",
	Message: "team_name already exists",
	Params: errors.Params{
		http.WithStatus(
			gohttp.StatusBadRequest,
		), // NOTE: почему 400 в исходном API? Больше на 409 похоже.
	},
}

var ErrPRExists = errors.Template{
	Code:    "PR_EXISTS",
	Message: "PR id already exists",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusConflict),
	},
}

var ErrPRMerged = errors.Template{
	Code:    "PR_MERGED",
	Message: "cannot reassign on merged PR",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusConflict),
	},
}

var ErrNotAssigned = errors.Template{
	Code:    "NOT_ASSIGNED",
	Message: "reviewer is not assigned to this PR",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusConflict), // NOTE: почему 409 в исходном API?
	},
}

var ErrNoCandidate = errors.Template{
	Code:    "NO_CANDIDATE",
	Message: "no active replacement candidate in team",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusConflict), // NOTE: почему 409 в исходном API?
	},
}

var ErrNotFound = errors.Template{
	Code:    "NOT_FOUND",
	Message: "resource not found",
	Params: errors.Params{
		http.WithStatus(gohttp.StatusNotFound),
	},
}

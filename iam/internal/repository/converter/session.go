package converter

import (
	"github.com/ChopX4/raketka/iam/internal/model"
	repoModel "github.com/ChopX4/raketka/iam/internal/repository/model"
)

func SessionToRepo(session model.Session) repoModel.Session {
	return repoModel.Session{
		SessionUUID: session.SessionUUID,
		UserUUID:    session.UserUUID,
	}
}

func SessionToModel(session repoModel.Session) model.Session {
	return model.Session{
		SessionUUID: session.SessionUUID,
		UserUUID:    session.UserUUID,
	}
}

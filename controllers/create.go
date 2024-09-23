package controllers

func (u *controllerGeneric[T, X]) Create(model T) (T, error) {
	m, err := u.svcT.Create(model)
	if err != nil {
		u.log.Error(err.Error())
		return m, err
	}

	return m, nil
}

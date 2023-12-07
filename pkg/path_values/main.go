package path_values

/*func ReadUUID(r *http.Request, name string) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName(name))
	if err != nil {
		return "", err
	}

	value := id.String()
	return value, nil
}

func ReadInt(r *http.Request, name string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName(name), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}*/

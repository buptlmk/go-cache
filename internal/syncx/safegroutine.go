package syncx

func SafedGroutine(fn func(), callBack func()) {

	go func() {
		if err := recover(); err != nil {
			if callBack != nil {
				callBack()
			}
		}

		fn()
	}()

}

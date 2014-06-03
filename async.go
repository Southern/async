package async

type Done func(error, ...interface{})

type Routine func(Done, ...interface{})

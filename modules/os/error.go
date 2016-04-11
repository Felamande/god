package os

type readError struct{
    msg  string
}

func (e readError)Error()string{
    return e.msg
}